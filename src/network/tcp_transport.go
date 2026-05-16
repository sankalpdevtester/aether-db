package network

import (
	"context"
	"fmt"
	"net"
	"sync"

	"google.golang.org/grpc"

	"aetherdb/src/raft"
	"aetherdb/src/rpc"
)

// TCPTransport represents a TCP transport layer for RPC communication.
type TCPTransport struct {
	// address is the TCP address to listen on.
	address string

	// server is the gRPC server instance.
	server *grpc.Server

	// pool is a connection pool for client connections.
	pool sync.Pool

	// raftNode is the Raft node instance.
	raftNode *raft.RaftNode
}

// NewTCPTransport creates a new TCP transport instance.
func NewTCPTransport(address string, raftNode *raft.RaftNode) *TCPTransport {
	return &TCPTransport{
		address:  address,
		raftNode:  raftNode,
		pool:      sync.Pool{},
		server:    grpc.NewServer(),
	}
}

// Start starts the TCP transport server.
func (t *TCPTransport) Start() error {
	// Register the RPC service.
	rpc.RegisterAetherDBServer(t.server, t)

	// Listen on the TCP address.
	listener, err := net.Listen("tcp", t.address)
	if err != nil {
		return err
	}

	// Start the gRPC server.
	go func() {
		if err := t.server.Serve(listener); err != nil {
			fmt.Printf("failed to start TCP transport server: %v\n", err)
		}
	}()

	return nil
}

// Stop stops the TCP transport server.
func (t *TCPTransport) Stop() {
	t.server.Stop()
}

// GetConnection returns a connection from the pool or creates a new one.
func (t *TCPTransport) GetConnection(address string) (net.Conn, error) {
	// Check if there's a connection available in the pool.
	conn := t.pool.Get()
	if conn != nil {
		return conn.(net.Conn), nil
	}

	// Create a new connection.
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// PutConnection returns a connection to the pool.
func (t *TCPTransport) PutConnection(conn net.Conn) {
	t.pool.Put(conn)
}

// AetherDB is the RPC service implementation.
func (t *TCPTransport) AetherDB(ctx context.Context, req *rpc.AetherDBRequest) (*rpc.AetherDBResponse, error) {
	// Handle the RPC request.
	switch req.Operation {
	case rpc.Operation_GET:
		return t.handleGet(ctx, req)
	case rpc.Operation_PUT:
		return t.handlePut(ctx, req)
	case rpc.Operation_DELETE:
		return t.handleDelete(ctx, req)
	default:
		return nil, fmt.Errorf("unknown operation: %s", req.Operation)
	}
}

// handleGet handles a GET operation.
func (t *TCPTransport) handleGet(ctx context.Context, req *rpc.AetherDBRequest) (*rpc.AetherDBResponse, error) {
	// Get the value from the Raft node.
	value, err := t.raftNode.Get(req.Key)
	if err != nil {
		return nil, err
	}

	// Return the response.
	return &rpc.AetherDBResponse{
		Value: value,
	}, nil
}

// handlePut handles a PUT operation.
func (t *TCPTransport) handlePut(ctx context.Context, req *rpc.AetherDBRequest) (*rpc.AetherDBResponse, error) {
	// Put the value into the Raft node.
	if err := t.raftNode.Put(req.Key, req.Value); err != nil {
		return nil, err
	}

	// Return the response.
	return &rpc.AetherDBResponse{}, nil
}

// handleDelete handles a DELETE operation.
func (t *TCPTransport) handleDelete(ctx context.Context, req *rpc.AetherDBRequest) (*rpc.AetherDBResponse, error) {
	// Delete the value from the Raft node.
	if err := t.raftNode.Delete(req.Key); err != nil {
		return nil, err
	}

	// Return the response.
	return &rpc.AetherDBResponse{}, nil
}