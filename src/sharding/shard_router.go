package sharding

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/AetherDB/aetherdb/src/network/rpc"
	"github.com/AetherDB/aetherdb/src/sharding/shard_manager"
)

// ShardRouter is responsible for routing requests to the correct shard.
type ShardRouter struct {
	shardManager *shard_manager.ShardManager
	rpcSvc       *rpc.RPCService
}

// NewShardRouter returns a new instance of ShardRouter.
func NewShardRouter(shardManager *shard_manager.ShardManager, rpcSvc *rpc.RPCService) *ShardRouter {
	return &ShardRouter{
		shardManager: shardManager,
		rpcSvc:       rpcSvc,
	}
}

// Route routes a request to the correct shard.
func (sr *ShardRouter) Route(ctx context.Context, req *http.Request) (*http.Response, error) {
	shardID := sr.getShardID(req)
	shard := sr.shardManager.GetShard(shardID)
	if shard == nil {
		return nil, fmt.Errorf("shard not found")
	}
	leader := shard.Leader
	if leader == nil {
		return nil, fmt.Errorf("leader not found")
	}
	endpoint := leader.Endpoint
	resp, err := sr.rpcSvc.Call(ctx, endpoint, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// getShardID returns the shard ID for a given request.
func (sr *ShardRouter) getShardID(req *http.Request) uint64 {
	// For demonstration purposes, assume the shard ID is in the request header.
	shardIDStr := req.Header.Get("X-Shard-ID")
	shardID, err := strconv.ParseUint(shardIDStr, 10, 64)
	if err != nil {
		log.Printf("error parsing shard ID: %v", err)
		return 0
	}
	return shardID
}

// Start starts the shard router.
func (sr *ShardRouter) Start(ctx context.Context) {
	log.Println("Shard router started")
	sr.rpcSvc.Start(ctx)
}

// Stop stops the shard router.
func (sr *ShardRouter) Stop(ctx context.Context) {
	log.Println("Shard router stopped")
	sr.rpcSvc.Stop(ctx)
}