package sharding

import (
	"context"
	"fmt"
	"log"

	"github.com/AetherDB/AetherDB/src/raft/consensus"
	"github.com/AetherDB/AetherDB/src/raft/raft_node"
)

// ShardReplicator is responsible for replicating data across shards
type ShardReplicator struct {
	shardManager *ShardManager
	consensus    *consensus.Consensus
}

// NewShardReplicator returns a new instance of ShardReplicator
func NewShardReplicator(shardManager *ShardManager, consensus *consensus.Consensus) *ShardReplicator {
	return &ShardReplicator{
		shardManager: shardManager,
		consensus:    consensus,
	}
}

// ReplicateData replicates data across shards
func (sr *ShardReplicator) ReplicateData(ctx context.Context, data []byte) error {
	shards := sr.shardManager.shards

	for _, shard := range shards {
		for _, node := range shard.Nodes {
			if err := node.ReplicateData(ctx, data); err != nil {
				return err
			}
		}
	}

	return nil
}

// Start starts the shard replicator
func (sr *ShardReplicator) Start(ctx context.Context) error {
	for _, shard := range sr.shardManager.shards {
		if err := shard.Start(ctx); err != nil {
			return err
		}
	}

	return nil
}

// Stop stops the shard replicator
func (sr *ShardReplicator) Stop(ctx context.Context) error {
	for _, shard := range sr.shardManager.shards {
		if err := shard.Stop(ctx); err != nil {
			return err
		}
	}

	return nil
}