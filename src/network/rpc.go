package network

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"sync"

	"github.com/AetherDB/aetherdb/raft"
)

// RPCRequest represents a request sent over the network
type RPCRequest struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// RPCResponse represents a response sent over the network
type RPCResponse struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
	Err  string       `json:"err"`
}

// RPCServer represents an RPC server
type RPCServer struct {
	listener net.Listener
	raftNode  *raft.RaftNode
	mu        sync.Mutex
}

// NewRPCServer creates a new RPC server
func NewRPCServer(raftNode *raft.RaftNode) (*RPCServer, error) {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		return nil, err
	}
	return &RPCServer{listener: listener, raftNode: raftNode}, nil
}

// Start starts the RPC server
func (s *RPCServer) Start() {
	go func() {
		for {
			conn, err := s.listener.Accept()
			if err != nil {
				fmt.Println(err)
				continue
			}
			go s.handleConnection(conn)
		}
	}()
}

// handleConnection handles a new connection
func (s *RPCServer) handleConnection(conn net.Conn) {
	defer conn.Close()
	var req RPCRequest
	err := json.NewDecoder(conn).Decode(&req)
	if err != nil {
		fmt.Println(err)
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	switch req.Type {
	case "get":
		s.handleGet(conn, req.Data)
	case "put":
		s.handlePut(conn, req.Data)
	case "delete":
		s.handleDelete(conn, req.Data)
	default:
		s.sendError(conn, "unknown request type")
	}
}

// handleGet handles a get request
func (s *RPCServer) handleGet(conn net.Conn, data interface{}) {
	key, ok := data.(string)
	if !ok {
		s.sendError(conn, "invalid data type")
		return
	}
	value, err := s.raftNode.Get(key)
	if err != nil {
		s.sendError(conn, err.Error())
		return
	}
	s.sendResponse(conn, value)
}

// handlePut handles a put request
func (s *RPCServer) handlePut(conn net.Conn, data interface{}) {
	var pair struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	err := json.Unmarshal([]byte(fmt.Sprintf("%v", data)), &pair)
	if err != nil {
		s.sendError(conn, "invalid data type")
		return
	}
	err = s.raftNode.Put(pair.Key, pair.Value)
	if err != nil {
		s.sendError(conn, err.Error())
		return
	}
	s.sendResponse(conn, nil)
}

// handleDelete handles a delete request
func (s *RPCServer) handleDelete(conn net.Conn, data interface{}) {
	key, ok := data.(string)
	if !ok {
		s.sendError(conn, "invalid data type")
		return
	}
	err := s.raftNode.Delete(key)
	if err != nil {
		s.sendError(conn, err.Error())
		return
	}
	s.sendResponse(conn, nil)
}

// sendResponse sends a response over the network
func (s *RPCServer) sendResponse(conn net.Conn, data interface{}) {
	resp := RPCResponse{Type: "response", Data: data}
	err := json.NewEncoder(conn).Encode(resp)
	if err != nil {
		fmt.Println(err)
	}
}

// sendError sends an error over the network
func (s *RPCServer) sendError(conn net.Conn, err string) {
	resp := RPCResponse{Type: "error", Err: err}
	err = json.NewEncoder(conn).Encode(resp)
	if err != nil {
		fmt.Println(err)
	}
}

// RPCClient represents an RPC client
type RPCClient struct {
	conn net.Conn
}

// NewRPCClient creates a new RPC client
func NewRPCClient(addr string) (*RPCClient, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &RPCClient{conn: conn}, nil
}

// Get sends a get request over the network
func (c *RPCClient) Get(key string) (string, error) {
	req := RPCRequest{Type: "get", Data: key}
	err := json.NewEncoder(c.conn).Encode(req)
	if err != nil {
		return "", err
	}
	var resp RPCResponse
	err = json.NewDecoder(c.conn).Decode(&resp)
	if err != nil {
		return "", err
	}
	if resp.Err != "" {
		return "", errors.New(resp.Err)
	}
	return resp.Data.(string), nil
}

// Put sends a put request over the network
func (c *RPCClient) Put(key string, value string) error {
	req := RPCRequest{Type: "put", Data: struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}{Key: key, Value: value}}
	err := json.NewEncoder(c.conn).Encode(req)
	if err != nil {
		return err
	}
	var resp RPCResponse
	err = json.NewDecoder(c.conn).Decode(&resp)
	if err != nil {
		return err
	}
	if resp.Err != "" {
		return errors.New(resp.Err)
	}
	return nil
}

// Delete sends a delete request over the network
func (c *RPCClient) Delete(key string) error {
	req := RPCRequest{Type: "delete", Data: key}
	err := json.NewEncoder(c.conn).Encode(req)
	if err != nil {
		return err
	}
	var resp RPCResponse
	err = json.NewDecoder(c.conn).Decode(&resp)
	if err != nil {
		return err
	}
	if resp.Err != "" {
		return errors.New(resp.Err)
	}
	return nil
}