// Package feature provides additional features for AetherDB.
package feature

import (
	"encoding/binary"
	"fmt"
	"log"
	"sync"

	"github.com/aetherdb/aetherdb/src/raft"
	"github.com/aetherdb/aetherdb/src/network"
)

// TimeSeriesIndex represents a time series index.
type TimeSeriesIndex struct {
	mu      sync.RWMutex
	index   map[uint64][]byte
	raftNode *raft.RaftNode
}

// NewTimeSeriesIndex returns a new time series index.
func NewTimeSeriesIndex(raftNode *raft.RaftNode) *TimeSeriesIndex {
	return &TimeSeriesIndex{
		index:   make(map[uint64][]byte),
		raftNode: raftNode,
	}
}

// Add adds a new entry to the time series index.
func (tsi *TimeSeriesIndex) Add(key uint64, value []byte) {
	tsi.mu.Lock()
	defer tsi.mu.Unlock()
	tsi.index[key] = value
}

// Get returns the value associated with the given key.
func (tsi *TimeSeriesIndex) Get(key uint64) []byte {
	tsi.mu.RLock()
	defer tsi.mu.RUnlock()
	return tsi.index[key]
}

// RangeQuery returns all values within the given range.
func (tsi *TimeSeriesIndex) RangeQuery(start uint64, end uint64) [][]byte {
	tsi.mu.RLock()
	defer tsi.mu.RUnlock()
	var values [][]byte
	for key := start; key <= end; key++ {
		value := tsi.index[key]
		if value != nil {
			values = append(values, value)
		}
	}
	return values
}

// Sync syncs the time series index with the Raft node.
func (tsi *TimeSeriesIndex) Sync() {
	tsi.mu.RLock()
	defer tsi.mu.RUnlock()
	for key, value := range tsi.index {
		tsi.raftNode.Apply(network.NewAppendRequest(key, value))
	}
}

// HandleAppendRequest handles an append request from the Raft node.
func (tsi *TimeSeriesIndex) HandleAppendRequest(req *network.AppendRequest) {
	tsi.mu.Lock()
	defer tsi.mu.Unlock()
	tsi.index[req.Key] = req.Value
}

func main() {
	// Example usage:
	raftNode := raft.NewRaftNode()
	tsi := NewTimeSeriesIndex(raftNode)
	tsi.Add(1, []byte("value1"))
	tsi.Add(2, []byte("value2"))
	values := tsi.RangeQuery(1, 2)
	for _, value := range values {
		log.Println(string(value))
	}
}