// Package utils provides utility functions for AetherDB.
package utils

import (
	"sync"
	"time"

	"github.com/AetherDB/aetherdb/src/feature/data_compression"
	"github.com/AetherDB/aetherdb/src/sharding"
)

// CacheConfig represents the configuration for the in-memory cache.
type CacheConfig struct {
	TTL time.Duration // Time to live for cache entries
}

// Cache is an in-memory cache with TTL for database query results.
type Cache struct {
	config CacheConfig
	cache  map[string][]byte
	mu     sync.RWMutex
}

// NewCache returns a new instance of the cache.
func NewCache(config CacheConfig) *Cache {
	return &Cache{
		config: config,
		cache:  make(map[string][]byte),
	}
}

// Get returns the cached result for the given key.
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, ok := c.cache[key]
	if !ok {
		return nil, false
	}

	// Check if the cache entry has expired
	if time.Since(time.Now()) > c.config.TTL {
		delete(c.cache, key)
		return nil, false
	}

	return value, true
}

// Set sets the cached result for the given key.
func (c *Cache) Set(key string, value []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Compress the value before caching
	compressedValue, err := data_compression.Compress(value)
	if err != nil {
		// Log the error and return
		return
	}

	c.cache[key] = compressedValue
}

// Delete deletes the cached result for the given key.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.cache, key)
}

// InvalidateShardCache invalidates the cache for a given shard.
func (c *Cache) InvalidateShardCache(shardID uint64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Find all cache entries for the given shard and delete them
	for key := range c.cache {
		if sharding.ShardRouter.GetShardID(key) == shardID {
			delete(c.cache, key)
		}
	}
}