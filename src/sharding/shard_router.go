package sharding

import (
	"context"
	"fmt"
	"log"

	"github.com/AetherDB/AetherDB/src/network/rpc"
	"github.com/AetherDB/AetherDB/src/raft/consensus"
)

// ShardRouter is responsible for routing requests to the correct shard
type ShardRouter struct {
	shardManager *ShardManager
	consensus    *consensus.Consensus
	rpcServer    *rpc.RPCServer
}

// NewShardRouter returns a new instance of ShardRouter
func NewShardRouter(shardManager *ShardManager, consensus *consensus.Consensus, rpcServer *rpc.RPCServer) *ShardRouter {
	return &ShardRouter{
		shardManager: shardManager,
		consensus:    consensus,
		rpcServer:    rpcServer,
	}
}

// RouteRequest routes a request to the correct shard
func (sr *ShardRouter) RouteRequest(ctx context.Context, request *rpc.Request) (*rpc.Response, error) {
	shardID := request.ShardID

	shard := sr.shardManager.GetShard(shardID)

	if shard == nil {
		return nil, fmt.Errorf("shard not found")
	}

	for _, node := range shard.Nodes {
		if err := node.RouteRequest(ctx, request); err != nil {
			return nil, err
		}
	}

	return &rpc.Response{}, nil
}

// Start starts the shard router
func (sr *ShardRouter) Start(ctx context.Context) error {
	if err := sr.rpcServer.Start(ctx); err != nil {
		return err
	}

	return nil
}

// Stop stops the shard router
func (sr *ShardRouter) Stop(ctx context.Context) error {
	if err := sr.rpcServer.Stop(ctx); err != nil {
		return err
	}

	return nil
}