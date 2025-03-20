package cache

import (
    "context"
    "sync"
)

type DistributedCache struct {
    nodes    []CacheNode
    balancer LoadBalancer
    mu       sync.RWMutex
}

type CacheNode struct {
    ID       string
    Endpoint string
    Weight   int
    Health   bool
}

func NewDistributedCache(nodes []CacheNode) *DistributedCache {
    return &DistributedCache{
        nodes:    nodes,
        balancer: newLoadBalancer(),
    }
}

func (c *DistributedCache) Get(ctx context.Context, key string) (interface{}, error) {
    node := c.balancer.SelectNode(key)
    return node.Get(ctx, key)
}

func (c *DistributedCache) Set(ctx context.Context, key string, value interface{}) error {
    nodes := c.balancer.GetReplicationNodes(key)
    return c.replicateToNodes(ctx, nodes, key, value)
}