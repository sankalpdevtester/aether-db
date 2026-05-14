package raft

import (
	"math/rand"
	"sync"
	"time"

	"github.com/AetherDB/aetherdb/src/feature"
	"github.com/AetherDB/aetherdb/src/utils"
)

// RaftNode represents a node in the Raft cluster
type RaftNode struct {
	id        string
	peers     map[string]*RaftNode
	state     RaftState
	term      int
	votedFor  string
	log       []LogEntry
	commitIndex int
	lastApplied int
	electionTimeout time.Duration
	heartbeatTimeout time.Duration
	rand *rand.Rand
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
	Term    int
	Command string
}

// NewRaftNode creates a new Raft node
func NewRaftNode(id string, peers map[string]*RaftNode) *RaftNode {
	node := &RaftNode{
		id:        id,
		peers:     peers,
		state:     Follower,
		term:      0,
		votedFor:  "",
		log:       make([]LogEntry, 0),
		commitIndex: 0,
		lastApplied: 0,
		electionTimeout: 150 * time.Millisecond,
		heartbeatTimeout: 50 * time.Millisecond,
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	return node
}

// Start starts the Raft node
func (n *RaftNode) Start() {
	n.mu.Lock()
	n.state = Follower
	n.term = 0
	n.votedFor = ""
	n.log = make([]LogEntry, 0)
	n.commitIndex = 0
	n.lastApplied = 0
	n.mu.Unlock()

	// Start election timer
	go n.startElectionTimer()

	// Start heartbeat timer
	go n.startHeartbeatTimer()
}

// startElectionTimer starts the election timer
func (n *RaftNode) startElectionTimer() {
	for {
		// Wait for election timeout
		time.Sleep(n.electionTimeout + time.Duration(n.rand.Intn(100))*time.Millisecond)

		n.mu.Lock()
		if n.state == Follower {
			// Become candidate and start election
			n.state = Candidate
			n.term++
			n.votedFor = n.id
			n.mu.Unlock()

			// Request votes from peers
			n.requestVotes()
		} else {
			n.mu.Unlock()
		}
	}
}

// startHeartbeatTimer starts the heartbeat timer
func (n *RaftNode) startHeartbeatTimer() {
	for {
		// Wait for heartbeat timeout
		time.Sleep(n.heartbeatTimeout)

		n.mu.Lock()
		if n.state == Leader {
			// Send heartbeat to peers
			n.sendHeartbeat()
		}
		n.mu.Unlock()
	}
}

// requestVotes requests votes from peers
func (n *RaftNode) requestVotes() {
	for _, peer := range n.peers {
		// Send request vote message to peer
		peer.handleRequestVote(n.id, n.term)
	}
}

// sendHeartbeat sends a heartbeat to peers
func (n *RaftNode) sendHeartbeat() {
	for _, peer := range n.peers {
		// Send heartbeat message to peer
		peer.handleHeartbeat(n.id, n.term)
	}
}

// handleRequestVote handles a request vote message from a peer
func (n *RaftNode) handleRequestVote(peerID string, term int) {
	n.mu.Lock()
	defer n.mu.Unlock()

	if term > n.term {
		// Update term and vote for peer
		n.term = term
		n.votedFor = peerID
	}
}

// handleHeartbeat handles a heartbeat message from a peer
func (n *RaftNode) handleHeartbeat(peerID string, term int) {
	n.mu.Lock()
	defer n.mu.Unlock()

	if term >= n.term {
		// Update term and reset election timer
		n.term = term
		n.state = Follower
	}
}

// AppendEntries appends entries to the log
func (n *RaftNode) AppendEntries(entries []LogEntry) {
	n.mu.Lock()
	defer n.mu.Unlock()

	// Append entries to log
	n.log = append(n.log, entries...)
}

// Commit commits the log up to the given index
func (n *RaftNode) Commit(index int) {
	n.mu.Lock()
	defer n.mu.Unlock()

	// Commit log up to index
	n.commitIndex = index
}

// Apply applies the committed log entries
func (n *RaftNode) Apply() {
	n.mu.Lock()
	defer n.mu.Unlock()

	// Apply committed log entries
	for i := n.lastApplied + 1; i <= n.commitIndex; i++ {
		// Apply log entry
		feature.ApplyLogEntry(n.log[i-1])
		n.lastApplied = i
	}
}