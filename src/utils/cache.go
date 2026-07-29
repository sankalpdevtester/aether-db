// Package utils provides utility functions for AetherDB.
package utils

import (
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

// CacheConfig represents the configuration for the cache.
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

// NewCache returns a new instance of the cache.
func NewCache(config CacheConfig) *Cache {
	c := cache.New(5*time.Minute, 10*time.Minute)
	return &Cache{
		cache: c,
	}
}

// Get returns the value associated with the given key.
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.cache.Get(key)
}

// Set sets the value associated with the given key.
func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache.Set(key, value, ttl)
}

// Delete removes the value associated with the given key.
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

// GetOrSet returns the value associated with the given key, or sets it if it doesn't exist.
func (c *Cache) GetOrSet(key string, value interface{}, ttl time.Duration) interface{} {
	c.mu.RLock()
	val, ok := c.cache.Get(key)
	c.mu.RUnlock()
	if ok {
		return val
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache.Set(key, value, ttl)
	return value
}

// Example usage:
func main() {
	cache := NewCache(CacheConfig{
		DefaultTTL: 1 * time.Hour,
		MaxSize:     1000,
	})
	cache.Set("key", "value", 1*time.Hour)
	val, ok := cache.Get("key")
	if ok {
		println(val.(string)) // prints "value"
	}
}