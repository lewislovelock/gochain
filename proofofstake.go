// Package main implements proof-of-stake system
package main

import (
	"bytes"         // for comparing and combining byte slices
	"crypto/sha256" // for hashing

	// for converting to binary
	"fmt"       // for printing
	"math/big"  // for working with large integers
	"math/rand" // for random number generation
	"time"      // for timestamps
)

// Validator represents a participant in the PoS system
type Validator struct {
	Address []byte // validator's address
	Stake   uint64 // amount of coins staked
	Balance uint64 // total balance including stake
}

// ProofOfStake represents a proof-of-stake system
type ProofOfStake struct {
	block      *Block       // pointer to the block being validated
	validators []*Validator // list of validators
	threshold  *big.Int     // threshold for valid blocks (similar to PoW target)
}

// NewProofOfStake builds and returns a ProofOfStake
func NewProofOfStake(b *Block) *ProofOfStake {
	// In a real implementation, validators would be loaded from a persistent store
	// Here we create some mock validators for demonstration
	validators := createMockValidators()

	// Set a threshold for valid blocks (simplified version)
	threshold := big.NewInt(1)
	threshold.Lsh(threshold, 255) // Very easy threshold for demonstration

	pos := &ProofOfStake{
		block:      b,
		validators: validators,
		threshold:  threshold,
	}
	return pos
}

// createMockValidators creates test validators
func createMockValidators() []*Validator {
	return []*Validator{
		{Address: []byte("validator1"), Stake: 1000, Balance: 5000},
		{Address: []byte("validator2"), Stake: 2000, Balance: 8000},
		{Address: []byte("validator3"), Stake: 3000, Balance: 10000},
	}
}

// selectValidator chooses a validator based on their stake
func (pos *ProofOfStake) selectValidator() *Validator {
	// Calculate total stake
	var totalStake uint64
	for _, v := range pos.validators {
		totalStake += v.Stake
	}

	// Random number between 0 and totalStake
	rand.Seed(time.Now().UnixNano())
	selection := rand.Uint64() % totalStake

	// Select validator based on stake weight
	var accumulator uint64
	for _, v := range pos.validators {
		accumulator += v.Stake
		if selection <= accumulator {
			return v
		}
	}

	return pos.validators[0] // fallback
}

// prepareData combines block fields for hashing
func (pos *ProofOfStake) prepareData(validator *Validator) []byte {
	return bytes.Join(
		[][]byte{
			pos.block.PrevBlockHash,
			pos.block.Data,
			IntToHex(pos.block.Timestamp),
			validator.Address,
			IntToHex(int64(validator.Stake)),
		},
		[]byte{},
	)
}

// Run performs the proof-of-stake consensus
// Returns validator address and resulting hash
func (pos *ProofOfStake) Run() ([]byte, []byte) {
	fmt.Printf("Selecting validator for new block...")

	// Select validator based on stake
	validator := pos.selectValidator()

	// Prepare and hash the block data
	data := pos.prepareData(validator)
	hash := sha256.Sum256(data)

	fmt.Printf("\nBlock forged by validator with stake: %d\n", validator.Stake)

	return validator.Address, hash[:]
}

// Validate verifies the proof-of-stake
func (pos *ProofOfStake) Validate() bool {
	// In a real implementation, we would:
	// 1. Verify the validator's signature
	// 2. Check if the validator has sufficient stake
	// 3. Verify the validator hasn't forged another block recently
	// 4. Check for double-spending

	// For demonstration, we'll do a simplified validation
	for _, v := range pos.validators {
		data := pos.prepareData(v)
		hash := sha256.Sum256(data)

		// Convert hash to big integer
		var hashInt big.Int
		hashInt.SetBytes(hash[:])

		// Check if hash is below threshold
		if hashInt.Cmp(pos.threshold) == -1 {
			return true
		}
	}

	return false
}
