```go
// Package feature provides secondary indexing for AetherDB.
package feature

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/AetherDB/aetherdb/src/main"
	"github.com/AetherDB/aetherdb/src/utils"
)

// SecondaryIndex represents a secondary index in AetherDB.
type SecondaryIndex struct {
	indexName string
	key       string
	values    map[string][]string
	mu        sync.RWMutex
}

// NewSecondaryIndex returns a new secondary index.
func NewSecondaryIndex(indexName, key string) *SecondaryIndex {
	return &SecondaryIndex{
		indexName: indexName,
		key:       key,
		values:    make(map[string][]string),
	}
}

// Add adds a new value to the secondary index.
func (si *SecondaryIndex) Add(value string, primaryKey string) error {
	si.mu.Lock()
	defer si.mu.Unlock()

	if _, ok := si.values[value]; !ok {
		si.values[value] = make([]string, 0)
	}

	si.values[value] = append(si.values[value], primaryKey)
	return nil
}

// Remove removes a value from the secondary index.
func (si *SecondaryIndex) Remove(value string, primaryKey string) error {
	si.mu.Lock()
	defer si.mu.Unlock()

	if _, ok := si.values[value]; !ok {
		return errors.New("value not found in index")
	}

	var newValues []string
	for _, v := range si.values[value] {
		if v != primaryKey {
			newValues = append(newValues, v)
		}
	}

	si.values[value] = newValues
	return nil
}

// Get returns all primary keys associated with a given value.
func (si *SecondaryIndex) Get(value string) ([]string, error) {
	si.mu.RLock()
	defer si.mu.RUnlock()

	if _, ok := si.values[value]; !ok {
		return nil, errors.New("value not found in index")
	}

	return si.values[value], nil
}

// IndexQuery represents a query on a secondary index.
type IndexQuery struct {
	indexName string
	key       string
	value     string
}

// NewIndexQuery returns a new index query.
func NewIndexQuery(indexName, key, value string) *IndexQuery {
	return &IndexQuery{
		indexName: indexName,
		key:       key,
		value:     value,
	}
}

// Execute executes the index query.
func (iq *IndexQuery) Execute(db *main.AetherDB) ([]string, error) {
	index, ok := db.SecondaryIndexes[iq.indexName]
	if !ok {
		return nil, errors.New("index not found")
	}

	return index.Get(iq.value)
}

// RegisterSecondaryIndex registers a new secondary index with AetherDB.
func RegisterSecondaryIndex(db *main.AetherDB, indexName, key string) error {
	index := NewSecondaryIndex(indexName, key)
	db.SecondaryIndexes[indexName] = index
	return nil
}

// UnregisterSecondaryIndex unregisters a secondary index from AetherDB.
func UnregisterSecondaryIndex(db *main.AetherDB, indexName string) error {
	delete(db.SecondaryIndexes, indexName)
	return nil
}

func init() {
	main.RegisterFeature("secondary_index", RegisterSecondaryIndex, UnregisterSecondaryIndex)
}
```