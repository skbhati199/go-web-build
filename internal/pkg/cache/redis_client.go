package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

type redisClient struct {
	endpoint string
	client   *redis.Client
}

func (c *redisClient) Get(ctx context.Context, key string) (interface{}, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var result interface{}
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *redisClient) Set(ctx context.Context, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, 24*time.Hour).Err()
}

func (c *redisClient) Delete(ctx context.Context, key string) error {
	if key == "" {
		return c.client.FlushDB(ctx).Err()
	}
	return c.client.Del(ctx, key).Err()
}

func (c *redisClient) Health(ctx context.Context) bool {
	return c.client.Ping(ctx).Err() == nil
}
