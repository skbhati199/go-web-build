package cache

import (
	"sync"
	"time"
)

type CacheEntry struct {
	Data      interface{}
	CreatedAt time.Time
	ExpiresAt time.Time
}

type TemplateCache struct {
	entries map[string]*CacheEntry
	mutex   sync.RWMutex
	ttl     time.Duration
}

func NewTemplateCache(ttl time.Duration) *TemplateCache {
	return &TemplateCache{
		entries: make(map[string]*CacheEntry),
		ttl:     ttl,
	}
}

func (c *TemplateCache) Set(key string, value interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.entries[key] = &CacheEntry{
		Data:      value,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(c.ttl),
	}
}

func (c *TemplateCache) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	entry, exists := c.entries[key]
	if !exists {
		return nil, false
	}

	if time.Now().After(entry.ExpiresAt) {
		delete(c.entries, key)
		return nil, false
	}

	return entry.Data, true
}

func (c *TemplateCache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.entries, key)
}

func (c *TemplateCache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.entries = make(map[string]*CacheEntry)
}

func (c *TemplateCache) Cleanup() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	now := time.Now()
	for key, entry := range c.entries {
		if now.After(entry.ExpiresAt) {
			delete(c.entries, key)
		}
	}
}
