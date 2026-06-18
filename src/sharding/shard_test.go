package sharding

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/AetherDB/AetherDB/src/raft/consensus"
	"github.com/AetherDB/AetherDB/src/raft/raft_node"
)

func TestShardManager(t *testing.T) {
	consensus := &consensus.Consensus{}
	shardManager := NewShardManager(consensus)

	shardID := uint64(1)
	nodes := []*raft_node.RaftNode{
		&raft_node.RaftNode{},
		&raft_node.RaftNode{},
	}

	shardManager.AddShard(shardID, nodes)

	shard := shardManager.GetShard(shardID)

	if shard == nil {
		t.Errorf("shard not found")
	}

	if err := shard.Start(context.Background()); err != nil {
		t.Errorf("failed to start shard: %v", err)
	}

	if err := shard.Stop(context.Background()); err != nil {
		t.Errorf("failed to stop shard: %v", err)
	}
}

func TestShardReplicator(t *testing.T) {
	consensus := &consensus.Consensus{}
	shardManager := NewShardManager(consensus)
	shardReplicator := NewShardReplicator(shardManager, consensus)

	shardID := uint64(1)
	nodes := []*raft_node.RaftNode{
		&raft_node.RaftNode{},
		&raft_node.RaftNode{},
	}

	shardManager.AddShard(shardID, nodes)

	data := []byte("test data")

	if err := shardReplicator.ReplicateData(context.Background(), data); err != nil {
		t.Errorf("failed to replicate data: %v", err)
	}
}

func TestShardRouter(t *testing.T) {
	consensus := &consensus.Consensus{}
	shardManager := NewShardManager(consensus)
	rpcServer := &rpc.RPCServer{}
	shardRouter := NewShardRouter(shardManager, consensus, rpcServer)

	shardID := uint64(1)
	nodes := []*raft_node.RaftNode{
		&raft_node.RaftNode{},
		&raft_node.RaftNode{},
	}

	shardManager.AddShard(shardID, nodes)

	request := &rpc.Request{
		ShardID: shardID,
	}

	if _, err := shardRouter.RouteRequest(context.Background(), request); err != nil {
		t.Errorf("failed to route request: %v", err)
	}
}