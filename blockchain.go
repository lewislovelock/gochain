// Package main is the entry point for our blockchain application
package main

import (
	"fmt"     // for printing
	"strconv" // for converting bool to string
	"time"    // for block timestamps
)

// Block represents each 'item' in the blockchain
type Block struct {
	Timestamp     int64  // when the block was created
	Data          []byte // the actual data/transactions in the block
	PrevBlockHash []byte // the hash of the previous block
	Hash          []byte // the hash of the current block
	ValidatorID   []byte // ID of miner (PoW) or validator (PoS)
}

// Blockchain is a series of validated Blocks
type Blockchain struct {
	blocks        []*Block      // slice of pointers to Block
	consensusType ConsensusType // type of consensus mechanism to use
}

// NewBlock creates and returns a new Block
func NewBlock(data string, prevBlockHash []byte, consensusType ConsensusType) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Data:          []byte(data),
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{},
		ValidatorID:   []byte{},
	}

	// Create consensus mechanism and run it
	consensus := NewConsensus(consensusType, block)
	validatorID, hash := consensus.Run()

	// Set the block's hash and validator ID
	block.Hash = hash
	block.ValidatorID = validatorID

	return block
}

// NewGenesisBlock creates and returns the genesis Block
func NewGenesisBlock(consensusType ConsensusType) *Block {
	return NewBlock("Genesis Block", []byte{}, consensusType)
}

// NewBlockchain creates a new Blockchain with genesis Block
func NewBlockchain(consensusType ConsensusType) *Blockchain {
	return &Blockchain{
		blocks:        []*Block{NewGenesisBlock(consensusType)},
		consensusType: consensusType,
	}
}

// AddBlock adds a new block to the blockchain
func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash, bc.consensusType)
	bc.blocks = append(bc.blocks, newBlock)
}

// SwitchConsensus changes the consensus mechanism
func (bc *Blockchain) SwitchConsensus(newType ConsensusType) {
	bc.consensusType = newType
}

func main() {
	// Create new blockchain with PoW
	fmt.Println("Creating blockchain with Proof of Work...")
	bc := NewBlockchain(POW)

	fmt.Println("Mining block 1 with PoW...")
	bc.AddBlock("Send 50 BTC to John")

	// Switch to PoS
	fmt.Println("\nSwitching to Proof of Stake...")
	bc.SwitchConsensus(POS)

	fmt.Println("Creating block 2 with PoS...")
	bc.AddBlock("Send 30 BTC to Jane")

	// Print all blocks in the blockchain
	for i, block := range bc.blocks {
		fmt.Printf("\nBlock %d:\n", i)
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("Validator ID: %s\n", block.ValidatorID)

		// Validate the block
		consensus := NewConsensus(bc.consensusType, block)
		fmt.Printf("Valid: %s\n", strconv.FormatBool(consensus.Validate()))
	}
}
