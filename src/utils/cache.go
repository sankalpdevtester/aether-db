// Package utils provides utility functions for the AetherDB project.
package utils

import (
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

// CacheConfig represents the configuration for the in-memory cache.
type CacheConfig struct {
	DefaultTTL time.Duration
}

// Cache is an in-memory cache with TTL for database query results.
type Cache struct {
	cache *cache.Cache
	mu    sync.RWMutex
}

// NewCache returns a new instance of the Cache.
func NewCache(config CacheConfig) *Cache {
	return &Cache{
		cache: cache.New(5*time.Minute, 10*time.Minute),
	}
}

// Get returns the value associated with the given key from the cache.
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.cache.Get(key)
}

// Set sets the value associated with the given key in the cache.
func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache.Set(key, value, ttl)
}

// Delete removes the value associated with the given key from the cache.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache.Delete(key)
}

// Invalidate invalidates all cache entries.
func (c *Cache) Invalidate() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache.Flush()
}

// GetCache returns the underlying cache instance.
func (c *Cache) GetCache() *cache.Cache {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.cache
}

// cacheKey generates a cache key based on the given parameters.
func cacheKey(shardID string, key string) string {
	return shardID + ":" + key
}

// Example usage:
func main() {
	cache := NewCache(CacheConfig{
		DefaultTTL: 1 * time.Minute,
	})

	// Set a value in the cache
	cache.Set("shard1:mykey", "myvalue", 30*time.Second)

	// Get a value from the cache
	value, found := cache.Get("shard1:mykey")
	if found {
		println(value.(string)) // prints "myvalue"
	}

	// Delete a value from the cache
	cache.Delete("shard1:mykey")

	// Invalidate all cache entries
	cache.Invalidate()
}