package sharding

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/AetherDB/aetherdb/src/raft/consensus"
	"github.com/AetherDB/aetherdb/src/sharding/shard_manager"
)

// ReplicationManager is responsible for managing replication in the distributed database.
type ReplicationManager struct {
	shardManager *shard_manager.ShardManager
	consensusSvc *consensus.ConsensusService
	replicas     map[uint64]*Replica
	replicaLock   sync.RWMutex
}

// NewReplicationManager returns a new instance of ReplicationManager.
func NewReplicationManager(shardManager *shard_manager.ShardManager, consensusSvc *consensus.ConsensusService) *ReplicationManager {
	return &ReplicationManager{
		shardManager: shardManager,
		consensusSvc: consensusSvc,
		replicas:     make(map[uint64]*Replica),
	}
}

// AddReplica adds a replica to the replication manager.
func (rm *ReplicationManager) AddReplica(replica *Replica) {
	rm.replicaLock.Lock()
	defer rm.replicaLock.Unlock()
	rm.replicas[replica.ID] = replica
}

// RemoveReplica removes a replica from the replication manager.
func (rm *ReplicationManager) RemoveReplica(replica *Replica) {
	rm.replicaLock.Lock()
	defer rm.replicaLock.Unlock()
	delete(rm.replicas, replica.ID)
}

// StartReplication starts replication for a given shard.
func (rm *ReplicationManager) StartReplication(ctx context.Context, shardID uint64) {
	shard := rm.shardManager.GetShard(shardID)
	if shard == nil {
		log.Printf("shard not found: %d", shardID)
		return
	}
	leader := shard.Leader
	if leader == nil {
		log.Printf("leader not found for shard: %d", shardID)
		return
	}
	followers := shard.Followers
	for _, follower := range followers {
		rm.startFollowerReplication(ctx, leader, follower)
	}
}

// startFollowerReplication starts replication for a given follower.
func (rm *ReplicationManager) startFollowerReplication(ctx context.Context, leader *Replica, follower *Replica) {
	// For demonstration purposes, assume the replication is done by calling a function on the follower.
	follower.StartReplication(ctx, leader)
}

// StopReplication stops replication for a given shard.
func (rm *ReplicationManager) StopReplication(ctx context.Context, shardID uint64) {
	shard := rm.shardManager.GetShard(shardID)
	if shard == nil {
		log.Printf("shard not found: %d", shardID)
		return
	}
	leader := shard.Leader
	if leader == nil {
		log.Printf("leader not found for shard: %d", shardID)
		return
	}
	followers := shard.Followers
	for _, follower := range followers {
		rm.stopFollowerReplication(ctx, leader, follower)
	}
}

// stopFollowerReplication stops replication for a given follower.
func (rm *ReplicationManager) stopFollowerReplication(ctx context.Context, leader *Replica, follower *Replica) {
	// For demonstration purposes, assume the replication is stopped by calling a function on the follower.
	follower.StopReplication(ctx, leader)
}