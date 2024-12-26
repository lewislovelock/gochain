// Package main defines consensus mechanisms
package main

// ConsensusType represents the type of consensus mechanism
type ConsensusType int

const (
	// POW represents Proof of Work consensus
	POW ConsensusType = iota
	// POS represents Proof of Stake consensus
	POS
)

// Consensus interface defines methods that any consensus mechanism must implement
type Consensus interface {
	// Run executes the consensus algorithm and returns necessary data
	Run() ([]byte, []byte) // returns validator/miner ID and hash
	// Validate verifies the block according to consensus rules
	Validate() bool
}

// NewConsensus creates a new consensus mechanism based on the type
func NewConsensus(consensusType ConsensusType, block *Block) Consensus {
	switch consensusType {
	case POS:
		return NewProofOfStake(block)
	default:
		return NewProofOfWork(block)
	}
}
