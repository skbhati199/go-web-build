package cache

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync/atomic"
	"time"
)

var (
	ErrInvalidStrategy = errors.New("invalid cache strategy")
	ErrValueTooLarge   = errors.New("value exceeds max size")
)

type CacheProvider interface {
	Store(ctx context.Context, key string, value interface{}, metadata Metadata) error
	Retrieve(ctx context.Context, key string) (interface{}, error)
	Delete(ctx context.Context, key string) error
	Clear(ctx context.Context) error
}

type Metadata struct {
	CreatedAt   time.Time
	ExpiresAt   time.Time
	Size        int64
	Compressed  bool
	ContentType string
}

func (s *CacheStrategy) generateMetadata(value interface{}) Metadata {
	size := calculateSize(value)
	return Metadata{
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(s.config.TTL),
		Size:        size,
		Compressed:  s.config.Compression && size > 1024,
		ContentType: detectContentType(value),
	}
}

func selectProvider(strategy string) CacheProvider {
	switch strategy {
	case "memory":
		return newMemoryCache()
	case "redis":
		return newRedisCache()
	// case "distributed":
	// 	return newDistributedCache()
	default:
		return newMemoryCache() // Default to memory cache
	}
}

func (s *CacheStrategy) Set(ctx context.Context, key string, value interface{}) error {
	metadata := s.generateMetadata(value)
	return s.provider.Store(ctx, key, value, metadata)
}

type CacheStrategy struct {
	provider CacheProvider
	config   Config
	metrics  *CacheMetrics
}

type CacheMetrics struct {
	Hits        int64
	Misses      int64
	Evictions   int64
	Size        int64
	LastUpdated time.Time
}

func (s *CacheStrategy) Get(ctx context.Context, key string) (interface{}, error) {
	if err := s.validateKey(key); err != nil {
		return nil, err
	}

	value, err := s.provider.Retrieve(ctx, key)
	if err != nil {
		atomic.AddInt64(&s.metrics.Misses, 1)
		return nil, err
	}

	if value == nil {
		atomic.AddInt64(&s.metrics.Misses, 1)
		return nil, nil
	}

	atomic.AddInt64(&s.metrics.Hits, 1)

	if s.config.Compression {
		return s.decompress(value)
	}

	return value, nil
}

func (s *CacheStrategy) decompress(value interface{}) (interface{}, error) {
	compressed, ok := value.([]byte)
	if !ok {
		return value, nil
	}

	reader := bytes.NewReader(compressed)
	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		return nil, fmt.Errorf("decompression failed: %w", err)
	}
	defer gzipReader.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, gzipReader); err != nil {
		return nil, fmt.Errorf("decompression read failed: %w", err)
	}

	return buf.Bytes(), nil
}

func (s *CacheStrategy) GetMetrics() CacheMetrics {
	return CacheMetrics{
		Hits:        atomic.LoadInt64(&s.metrics.Hits),
		Misses:      atomic.LoadInt64(&s.metrics.Misses),
		Evictions:   atomic.LoadInt64(&s.metrics.Evictions),
		Size:        atomic.LoadInt64(&s.metrics.Size),
		LastUpdated: s.metrics.LastUpdated,
	}
}

func NewCacheStrategy(config Config) *CacheStrategy {
	if config.TTL == 0 {
		config.TTL = 24 * time.Hour
	}

	return &CacheStrategy{
		provider: selectProvider(config.Strategy),
		config:   config,
		metrics: &CacheMetrics{
			LastUpdated: time.Now(),
		},
	}
}

func (s *CacheStrategy) Delete(ctx context.Context, key string) error {
	return s.provider.Delete(ctx, key)
}

func (s *CacheStrategy) Clear(ctx context.Context) error {
	return s.provider.Clear(ctx)
}

func calculateSize(value interface{}) int64 {
	switch v := value.(type) {
	case string:
		return int64(len(v))
	case []byte:
		return int64(len(v))
	case nil:
		return 0
	default:
		if data, err := json.Marshal(v); err == nil {
			return int64(len(data))
		}
		return 0
	}
}

func (s *CacheStrategy) validateKey(key string) error {
	if key == "" {
		return errors.New("cache key cannot be empty")
	}
	if len(key) > 256 {
		return errors.New("cache key too long")
	}
	return nil
}

func (s *CacheStrategy) compress(value interface{}) ([]byte, error) {
	var data []byte
	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		jsonData, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		data = jsonData
	}

	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	if _, err := w.Write(data); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func detectContentType(value interface{}) string {
	switch value.(type) {
	case string:
		return "text/plain"
	case []byte:
		return "application/octet-stream"
	default:
		return "application/json"
	}
}

func (s *CacheStrategy) validateValue(value interface{}) error {
	size := calculateSize(value)
	if s.config.MaxSize > 0 && size > s.config.MaxSize {
		return fmt.Errorf("%w: size %d exceeds limit %d", ErrValueTooLarge, size, s.config.MaxSize)
	}
	return nil
}

type Config struct {
	Strategy     string
	TTL          time.Duration
	MaxSize      int64
	Compression  bool
	Distribution []string
}
