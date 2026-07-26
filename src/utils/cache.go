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
	MaxSize    int
}

// Cache is an in-memory cache with TTL for database query results.
type Cache struct {
	cache *cache.Cache
	mu    sync.RWMutex
}

// NewCache returns a new instance of the Cache.
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

// Flush removes all values from the cache.
func (c *Cache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache.Flush()
}

// GetOrSet returns the value associated with the given key from the cache.
// If the key is not present in the cache, it sets the value and returns it.
func (c *Cache) GetOrSet(key string, value interface{}, ttl time.Duration) interface{} {
	c.mu.RLock()
	val, found := c.cache.Get(key)
	c.mu.RUnlock()
	if found {
		return val
	}
	c.mu.Lock()
	c.cache.Set(key, value, ttl)
	c.mu.Unlock()
	return value
}

// Example usage:
func main() {
	cache := NewCache(CacheConfig{
		DefaultTTL: 1 * time.Hour,
		MaxSize:     1000,
	})

	// Set a value in the cache
	cache.Set("key", "value", 1*time.Hour)

	// Get a value from the cache
	val, found := cache.Get("key")
	if found {
		println(val.(string)) // Output: value
	}

	// Delete a value from the cache
	cache.Delete("key")

	// Flush the cache
	cache.Flush()
}