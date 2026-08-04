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
	// CleanupInterval is the interval at which the cache is cleaned up.
	CleanupInterval time.Duration
}

// Cache is an in-memory cache with TTL for database query results.
type Cache struct {
	cache *cache.Cache
	mu    sync.RWMutex
}

// NewCache returns a new instance of the Cache.
func NewCache(config CacheConfig) *Cache {
	c := cache.New(config.DefaultTTL, config.CleanupInterval)
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

// Delete removes the value associated with the given key from the cache.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache.Delete(key)
}

// Flush removes all values from the cache.
func (c *Cache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache.Flush()
}

// Example usage:
func main() {
	cacheConfig := CacheConfig{
		DefaultTTL:       5 * time.Minute,
		CleanupInterval: 10 * time.Minute,
	}
	cache := NewCache(cacheConfig)

	// Set a value in the cache
	cache.Set("key", "value")

	// Get a value from the cache
	value, found := cache.Get("key")
	if found {
		println(value.(string)) // prints "value"
	}

	// Delete a value from the cache
	cache.Delete("key")

	// Flush the cache
	cache.Flush()
}