package sharding

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/AetherDB/aetherdb/src/raft/consensus"
	"github.com/AetherDB/aetherdb/src/raft/raft_node"
)

// ShardManager is responsible for managing shards in the distributed database.
type ShardManager struct {
	shards       map[uint64]*Shard
	shardLock    sync.RWMutex
	consensusSvc *consensus.ConsensusService
	raftNode     *raft_node.RaftNode
}

// NewShardManager returns a new instance of ShardManager.
func NewShardManager(consensusSvc *consensus.ConsensusService, raftNode *raft_node.RaftNode) *ShardManager {
	return &ShardManager{
		shards:       make(map[uint64]*Shard),
		consensusSvc: consensusSvc,
		raftNode:     raftNode,
	}
}

// AddShard adds a new shard to the shard manager.
func (sm *ShardManager) AddShard(shardID uint64, shard *Shard) {
	sm.shardLock.Lock()
	defer sm.shardLock.Unlock()
	sm.shards[shardID] = shard
}

// GetShard returns a shard by its ID.
func (sm *ShardManager) GetShard(shardID uint64) *Shard {
	sm.shardLock.RLock()
	defer sm.shardLock.RUnlock()
	return sm.shards[shardID]
}

// RemoveShard removes a shard from the shard manager.
func (sm *ShardManager) RemoveShard(shardID uint64) {
	sm.shardLock.Lock()
	defer sm.shardLock.Unlock()
	delete(sm.shards, shardID)
}

// Start starts the shard manager.
func (sm *ShardManager) Start(ctx context.Context) {
	log.Println("Shard manager started")
	sm.raftNode.Start(ctx)
	sm.consensusSvc.Start(ctx)
}

// Stop stops the shard manager.
func (sm *ShardManager) Stop(ctx context.Context) {
	log.Println("Shard manager stopped")
	sm.raftNode.Stop(ctx)
	sm.consensusSvc.Stop(ctx)
}

// Shard represents a shard in the distributed database.
type Shard struct {
	ID        uint64
	Replicas  []*Replica
	Leader    *Replica
	Followers []*Replica
}

// NewShard returns a new instance of Shard.
func NewShard(id uint64) *Shard {
	return &Shard{
		ID: id,
	}
}

// AddReplica adds a replica to the shard.
func (s *Shard) AddReplica(replica *Replica) {
	s.Replicas = append(s.Replicas, replica)
}

// RemoveReplica removes a replica from the shard.
func (s *Shard) RemoveReplica(replica *Replica) {
	for i, r := range s.Replicas {
		if r == replica {
			s.Replicas = append(s.Replicas[:i], s.Replicas[i+1:]...)
			return
		}
	}
}

// Replica represents a replica in the shard.
type Replica struct {
	ID        uint64
	Shard     *Shard
	Endpoint string
}

// NewReplica returns a new instance of Replica.
func NewReplica(id uint64, shard *Shard, endpoint string) *Replica {
	return &Replica{
		ID:        id,
		Shard:     shard,
		Endpoint: endpoint,
	}
}