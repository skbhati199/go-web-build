package cache

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

type DistributedCache struct {
	nodes    []CacheNode
	balancer LoadBalancer
	mu       sync.RWMutex
	options  CacheOptions
}

type CacheNode struct {
	ID       string
	Endpoint string
	Weight   int
	Health   bool
	client   CacheClient
}

type CacheOptions struct {
	ReplicationFactor int
	Timeout           time.Duration
	RetryAttempts     int
	ConsistencyLevel  string
}

type CacheClient interface {
	Get(ctx context.Context, key string) (interface{}, error)
	Set(ctx context.Context, key string, value interface{}) error
	Delete(ctx context.Context, key string) error
	Health(ctx context.Context) bool
}

func NewDistributedCache(nodes []CacheNode, opts CacheOptions) *DistributedCache {
	if opts.ReplicationFactor <= 0 {
		opts.ReplicationFactor = 2
	}
	if opts.Timeout <= 0 {
		opts.Timeout = 5 * time.Second
	}

	return &DistributedCache{
		nodes:    nodes,
		balancer: newLoadBalancer(),
		options:  opts,
	}
}

func (c *DistributedCache) Get(ctx context.Context, key string) (interface{}, error) {
	node := c.balancer.SelectNode(key)
	if !node.Health {
		return nil, fmt.Errorf("selected node is unhealthy: %s", node.ID)
	}

	ctx, cancel := context.WithTimeout(ctx, c.options.Timeout)
	defer cancel()

	value, err := node.client.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("cache get failed: %w", err)
	}

	return value, nil
}

func (c *DistributedCache) Set(ctx context.Context, key string, value interface{}) error {
	nodes := c.balancer.GetReplicationNodes(key)
	return c.replicateToNodes(ctx, nodes, key, value)
}

func (c *DistributedCache) replicateToNodes(ctx context.Context, nodes []CacheNode, key string, value interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, c.options.Timeout)
	defer cancel()

	var wg sync.WaitGroup
	errChan := make(chan error, len(nodes))

	for _, node := range nodes {
		if !node.Health {
			continue
		}

		wg.Add(1)
		go func(n CacheNode) {
			defer wg.Done()
			if err := n.client.Set(ctx, key, value); err != nil {
				errChan <- fmt.Errorf("replication to node %s failed: %w", n.ID, err)
			}
		}(node)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return fmt.Errorf("replication failed: %v", errors)
	}

	return nil
}

func (c *DistributedCache) Delete(ctx context.Context, key string) error {
	nodes := c.balancer.GetReplicationNodes(key)
	var errors []error

	for _, node := range nodes {
		if err := node.client.Delete(ctx, key); err != nil {
			errors = append(errors, fmt.Errorf("delete from node %s failed: %w", node.ID, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("delete failed: %v", errors)
	}

	return nil
}

func (c *DistributedCache) Clear(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var errors []error
	for _, node := range c.nodes {
		if !node.Health {
			continue
		}

		if err := node.client.Delete(ctx, ""); err != nil {
			errors = append(errors, fmt.Errorf("clear node %s failed: %w", node.ID, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("clear cache failed: %v", errors)
	}

	return nil
}

func (c *DistributedCache) MonitorHealth(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.checkNodesHealth(ctx)
		}
	}
}

func (c *DistributedCache) checkNodesHealth(ctx context.Context) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i := range c.nodes {
		c.nodes[i].Health = c.nodes[i].client.Health(ctx)
	}
}

// func newDistributedCache() CacheProvider {
// 	nodes := []CacheNode{
// 		{
// 			ID:       "node-1",
// 			Endpoint: "localhost:6379",
// 			Weight:   100,
// 			Health:   true,
// 			client:   newRedisClient("localhost:6379"),
// 		},
// 		{
// 			ID:       "node-2",
// 			Endpoint: "localhost:6380",
// 			Weight:   100,
// 			Health:   true,
// 			client:   newRedisClient("localhost:6380"),
// 		},
// 	}

// 	options := CacheOptions{
// 		ReplicationFactor: 2,
// 		Timeout:           5 * time.Second,
// 		RetryAttempts:     3,
// 		ConsistencyLevel:  "quorum",
// 	}

// 	cache := NewDistributedCache(nodes, options)

// 	// Start health monitoring in background
// 	go func() {
// 		ctx := context.Background()
// 		cache.MonitorHealth(ctx)
// 	}()

// 	return cache
// }

func newRedisClient(endpoint string) CacheClient {
	return &redisClient{
		endpoint: endpoint,
		client:   redis.NewClient(&redis.Options{Addr: endpoint}),
	}
}
