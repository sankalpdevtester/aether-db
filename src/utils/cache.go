// Package utils provides utility functions for the AetherDB project.
package utils

import (
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

// CacheConfig represents the configuration for the in-memory cache.
type CacheConfig struct {
	TTL time.Duration
}

// Cache is an in-memory cache with a TTL (time to live) for stored values.
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

// Get returns the value associated with the given key if it exists in the cache.
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

// cacheResult caches the result of a database query.
func cacheResult(cache *Cache, key string, query func() (interface{}, error)) (interface{}, error) {
	// Check if the result is already cached.
	if cachedValue, found := cache.Get(key); found {
		return cachedValue, nil
	}

	// If not, execute the query and cache the result.
	result, err := query()
	if err != nil {
		return nil, err
	}

	// Cache the result with a TTL.
	cache.Set(key, result)

	return result, nil
}

// Example usage:
func main() {
	cache := NewCache(CacheConfig{TTL: 1 * time.Minute})

	// Define a query function that retrieves data from the database.
	queryData := func() (interface{}, error) {
		// Simulate a database query.
		time.Sleep(100 * time.Millisecond)
		return "Query result", nil
	}

	// Cache the result of the query.
	result, err := cacheResult(cache, "query_result", queryData)
	if err != nil {
		panic(err)
	}

	// Print the cached result.
	println(result.(string))
}