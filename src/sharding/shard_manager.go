package sharding

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/AetherDB/AetherDB/src/raft/consensus"
	"github.com/AetherDB/AetherDB/src/raft/raft_node"
)

// ShardManager is responsible for managing shards in the distributed database
type ShardManager struct {
	shards    map[uint64]*Shard
	mu        sync.RWMutex
	consensus *consensus.Consensus
}

// NewShardManager returns a new instance of ShardManager
func NewShardManager(consensus *consensus.Consensus) *ShardManager {
	return &ShardManager{
		shards:    make(map[uint64]*Shard),
		consensus: consensus,
	}
}

// AddShard adds a new shard to the shard manager
func (sm *ShardManager) AddShard(shardID uint64, nodes []*raft_node.RaftNode) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	shard := &Shard{
		ID:    shardID,
		Nodes: nodes,
	}

	sm.shards[shardID] = shard
}

// GetShard returns a shard by its ID
func (sm *ShardManager) GetShard(shardID uint64) *Shard {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	return sm.shards[shardID]
}

// RemoveShard removes a shard from the shard manager
func (sm *ShardManager) RemoveShard(shardID uint64) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	delete(sm.shards, shardID)
}

// Shard represents a shard in the distributed database
type Shard struct {
	ID    uint64
	Nodes []*raft_node.RaftNode
}

// Start starts the shard
func (s *Shard) Start(ctx context.Context) error {
	for _, node := range s.Nodes {
		if err := node.Start(ctx); err != nil {
			return err
		}
	}

	return nil
}

// Stop stops the shard
func (s *Shard) Stop(ctx context.Context) error {
	for _, node := range s.Nodes {
		if err := node.Stop(ctx); err != nil {
			return err
		}
	}

	return nil
}