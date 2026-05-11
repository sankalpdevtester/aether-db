package storage

import (
	"sync"
	"time"
)

type Store struct {
	data  map[string]Entry
	mutex sync.RWMutex
}

type Entry struct {
	Value     string
	Timestamp time.Time
	TTL       int64 // seconds, 0 = no expiry
}

func NewStore() *Store {
	return &Store{
		data: make(map[string]Entry),
	}
}

func (s *Store) Set(key, value string, ttl int64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	s.data[key] = Entry{
		Value:     value,
		Timestamp: time.Now(),
		TTL:       ttl,
	}
}

func (s *Store) Get(key string) (string, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	entry, exists := s.data[key]
	if !exists {
		return "", false
	}
	
	// Check if expired
	if entry.TTL > 0 {
		elapsed := time.Since(entry.Timestamp).Seconds()
		if elapsed > float64(entry.TTL) {
			return "", false
		}
	}
	
	return entry.Value, true
}

func (s *Store) Delete(key string) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	_, exists := s.data[key]
	if exists {
		delete(s.data, key)
		return true
	}
	return false
}

func (s *Store) Keys() []string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	keys := make([]string, 0, len(s.data))
	for k := range s.data {
		keys = append(keys, k)
	}
	return keys
}

func (s *Store) Close() {
	// Cleanup if needed
}
