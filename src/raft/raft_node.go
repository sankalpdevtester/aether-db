```go
package raft

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

// RaftNode represents a node in the Raft cluster
type RaftNode struct {
	id        string
	peers     map[string]*RaftNode
	state     RaftState
	currentTerm uint64
	votedFor   string
	log        []LogEntry
	commitIndex uint64
	lastApplied uint64
	electionTimeout time.Duration
	heartbeatTimeout time.Duration
	mu sync.RWMutex
}

// RaftState represents the state of a Raft node
type RaftState int

const (
	Follower RaftState = iota
	Candidate
	Leader
)

// LogEntry represents a log entry in the Raft log
type LogEntry struct {
	Term    uint64
	Command string
}

// NewRaftNode returns a new Raft node
func NewRaftNode(id string, peers map[string]*RaftNode) *RaftNode {
	return &RaftNode{
		id:        id,
		peers:     peers,
		state:     Follower,
		currentTerm: 0,
		votedFor:   "",
		log:        make([]LogEntry, 0),
		commitIndex: 0,
		lastApplied: 0,
		electionTimeout: time.Duration(rand.Intn(150)) * time.Millisecond,
		heartbeatTimeout: 50 * time.Millisecond,
	}
}

// Start starts the Raft node
func (rn *RaftNode) Start() {
	rn.mu.Lock()
	rn.state = Follower
	rn.mu.Unlock()
	go rn.runElectionTimer()
}

// runElectionTimer runs the election timer
func (rn *RaftNode) runElectionTimer() {
	for {
		time.Sleep(rn.electionTimeout)
		rn.mu.Lock()
		if rn.state == Follower {
			rn.startElection()
		}
		rn.mu.Unlock()
	}
}

// startElection starts a new election
func (rn *RaftNode) startElection() {
	rn.currentTerm++
	rn.votedFor = rn.id
	rn.state = Candidate
	votes := 1
	for _, peer := range rn.peers {
		if peer.voteFor(rn.id, rn.currentTerm) {
			votes++
		}
	}
	if votes > len(rn.peers)/2 {
		rn.state = Leader
		rn.runLeader()
	}
}

// voteFor votes for a candidate
func (rn *RaftNode) voteFor(candidate string, term uint64) bool {
	rn.mu.Lock()
	defer rn.mu.Unlock()
	if term < rn.currentTerm {
		return false
	}
	if term > rn.currentTerm {
		rn.currentTerm = term
		rn.votedFor = ""
	}
	if rn.votedFor == "" || rn.votedFor == candidate {
		rn.votedFor = candidate
		return true
	}
	return false
}

// runLeader runs the leader
func (rn *RaftNode) runLeader() {
	go rn.runHeartbeatTimer()
}

// runHeartbeatTimer runs the heartbeat timer
func (rn *RaftNode) runHeartbeatTimer() {
	for {
		time.Sleep(rn.heartbeatTimeout)
		rn.mu.Lock()
		if rn.state == Leader {
			rn.sendHeartbeats()
		}
		rn.mu.Unlock()
	}
}

// sendHeartbeats sends heartbeats to all peers
func (rn *RaftNode) sendHeartbeats() {
	for _, peer := range rn.peers {
		peer.receiveHeartbeat(rn.id, rn.currentTerm)
	}
}

// receiveHeartbeat receives a heartbeat from the leader
func (rn *RaftNode) receiveHeartbeat(leader string, term uint64) {
	rn.mu.Lock()
	defer rn.mu.Unlock()
	if term < rn.currentTerm {
		return
	}
	if term > rn.currentTerm {
		rn.currentTerm = term
		rn.votedFor = ""
	}
	rn.state = Follower
	rn.electionTimeout = time.Duration(rand.Intn(150)) * time.Millisecond
}

func main() {
	peers := make(map[string]*RaftNode)
	peers["node1"] = NewRaftNode("node1", peers)
	peers["node2"] = NewRaftNode("node2", peers)
	peers["node3"] = NewRaftNode("node3", peers)
	peers["node1"].Start()
	peers["node2"].Start()
	peers["node3"].Start()
	log.Println("Raft nodes started")
	select {}
}
```