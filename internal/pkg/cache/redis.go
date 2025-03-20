package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type redisCache struct {
	client  *redis.Client
	options *redis.Options
}

func newRedisCache() CacheProvider {
	options := &redis.Options{
		Addr:     "localhost:6379",
		Password: "", // Set from environment
		DB:       0,
	}

	return &redisCache{
		client:  redis.NewClient(options),
		options: options,
	}
}

func (c *redisCache) Store(ctx context.Context, key string, value interface{}, metadata Metadata) error {
	data := struct {
		Value    interface{} `json:"value"`
		Metadata Metadata    `json:"metadata"`
	}{
		Value:    value,
		Metadata: metadata,
	}

	encoded, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to encode cache data: %w", err)
	}

	ttl := metadata.ExpiresAt.Sub(time.Now())
	return c.client.Set(ctx, key, encoded, ttl).Err()
}

func (c *redisCache) Retrieve(ctx context.Context, key string) (interface{}, error) {
	result, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to retrieve from cache: %w", err)
	}

	var data struct {
		Value    interface{} `json:"value"`
		Metadata Metadata    `json:"metadata"`
	}

	if err := json.Unmarshal(result, &data); err != nil {
		return nil, fmt.Errorf("failed to decode cache data: %w", err)
	}

	if time.Now().After(data.Metadata.ExpiresAt) {
		go c.Delete(ctx, key)
		return nil, nil
	}

	return data.Value, nil
}

func (c *redisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func (c *redisCache) Clear(ctx context.Context) error {
	return c.client.FlushDB(ctx).Err()
}

func (c *redisCache) Close() error {
	return c.client.Close()
}
