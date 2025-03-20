package cache

import (
	"context"
	"sync"
	"time"
)

type memoryCache struct {
	data    map[string]*cacheEntry
	mu      sync.RWMutex
	maxSize int64
	janitor *time.Ticker
}

type cacheEntry struct {
	value    interface{}
	metadata Metadata
}

func newMemoryCache() CacheProvider {
	mc := &memoryCache{
		data:    make(map[string]*cacheEntry),
		janitor: time.NewTicker(5 * time.Minute),
	}
	go mc.startCleanup()
	return mc
}

func (c *memoryCache) Store(ctx context.Context, key string, value interface{}, metadata Metadata) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = &cacheEntry{
		value:    value,
		metadata: metadata,
	}
	return nil
}

func (c *memoryCache) Retrieve(ctx context.Context, key string) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.data[key]
	if !exists {
		return nil, nil
	}

	if time.Now().After(entry.metadata.ExpiresAt) {
		go c.Delete(ctx, key)
		return nil, nil
	}

	return entry.value, nil
}

func (c *memoryCache) Delete(ctx context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, key)
	return nil
}

func (c *memoryCache) Clear(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = make(map[string]*cacheEntry)
	return nil
}

func (c *memoryCache) startCleanup() {
	for range c.janitor.C {
		c.cleanup()
	}
}

func (c *memoryCache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, entry := range c.data {
		if now.After(entry.metadata.ExpiresAt) {
			delete(c.data, key)
		}
	}
}

func (c *memoryCache) Stop() {
	c.janitor.Stop()
}
