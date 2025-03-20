package cache

import (
    "context"
    "time"
)

type CacheStrategy struct {
    provider CacheProvider
    config   Config
}

type Config struct {
    Strategy     string
    TTL          time.Duration
    MaxSize      int64
    Compression  bool
    Distribution []string
}

func NewCacheStrategy(config Config) *CacheStrategy {
    return &CacheStrategy{
        provider: selectProvider(config.Strategy),
        config:   config,
    }
}

func (s *CacheStrategy) Set(ctx context.Context, key string, value interface{}) error {
    metadata := s.generateMetadata(value)
    return s.provider.Store(ctx, key, value, metadata)
}

func (s *CacheStrategy) Get(ctx context.Context, key string) (interface{}, error) {
    return s.provider.Retrieve(ctx, key)
}