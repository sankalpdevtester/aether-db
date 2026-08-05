// Package utils provides utility functions for the AetherDB project.
package utils

import (
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

// CacheConfig represents the configuration for the in-memory cache.
type CacheConfig struct {
	// DefaultTTL is the default time-to-live for cache entries.
	DefaultTTL time.Duration
	// MaxSize is the maximum number of entries in the cache.
	MaxSize int
}

// Cache is an in-memory cache with TTL for database query results.
type Cache struct {
	cache *cache.Cache
	mu    sync.RWMutex
}

// NewCache returns a new instance of the in-memory cache.
func NewCache(config CacheConfig) *Cache {
	c := cache.New(5*time.Minute, 10*time.Minute)
	return &Cache{
		cache: c,
	}
}

// Get returns the value associated with the given key from the cache.
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.cache.Get(key)
}

// Set sets the value associated with the given key in the cache.
func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache.Set(key, value, cache.DefaultExpiration)
}

// Delete removes the entry associated with the given key from the cache.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache.Delete(key)
}

// Invalidate invalidates all entries in the cache.
func (c *Cache) Invalidate() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache.Flush()
}

// Example usage:
func main() {
	cache := NewCache(CacheConfig{
		DefaultTTL: 1 * time.Hour,
		MaxSize:     1000,
	})

	// Set a value in the cache
	cache.Set("key", "value")

	// Get a value from the cache
	value, found := cache.Get("key")
	if found {
		println(value.(string)) // prints: value
	}

	// Delete a value from the cache
	cache.Delete("key")

	// Invalidate all entries in the cache
	cache.Invalidate()
}
``` 
// Integration with existing files:
// The cache can be used in the src/feature/data_compression.go file to cache compressed data.
// The cache can be used in the src/feature/time_series_index.go file to cache time series data.
// The cache can be used in the src/feature/secondary_index.go file to cache secondary index data.
// The cache can be used in the src/sharding/shard_manager.go file to cache shard metadata.
// The cache can be used in the src/sharding/shard_replicator.go file to cache replicated data.
// The cache can be used in the src/sharding/shard_router.go file to cache routing information.