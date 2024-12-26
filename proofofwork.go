// Package main implements proof-of-work system
package main

import (
	"bytes"           // for comparing and combining byte slices
	"crypto/sha256"   // for hashing
	"encoding/binary" // for converting to binary
	"fmt"             // for printing
	"math"            // for math operations
	"math/big"        // for working with large integers
)

// Difficulty of mining. In Bitcoin this is adjusted dynamically.
// Here we'll keep it constant for simplicity
const targetBits = 16

// ProofOfWork represents a proof-of-work system
type ProofOfWork struct {
	block  *Block   // pointer to the block for which we're calculating proof-of-work
	target *big.Int // target threshold below which hash must be
}

// NewProofOfWork builds and returns a ProofOfWork
func NewProofOfWork(b *Block) *ProofOfWork {
	// Initialize a big integer as 1
	target := big.NewInt(1)
	// Left shift it by (256 - targetBits)
	// This sets our target threshold: any hash below this is valid
	target.Lsh(target, uint(256-targetBits))

	pow := &ProofOfWork{b, target}
	return pow
}

// prepareData combines block fields with nonce and targetBits for hashing
func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.Data,
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)
	return data
}

// Run performs the proof-of-work computation
// Returns miner ID (nonce as bytes) and resulting hash
func (pow *ProofOfWork) Run() ([]byte, []byte) {
	var hashInt big.Int // holds the integer representation of our hash
	var hash [32]byte   // holds the actual hash bytes
	nonce := 0          // counter we'll increment while mining

	fmt.Printf("Mining a new block...")

	// Essentially infinite loop until we find a valid hash
	for nonce < math.MaxInt64 {
		// Prepare data for hashing
		data := pow.prepareData(nonce)
		// Calculate hash of the data
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash) // Show mining progress

		// Convert hash to big integer
		hashInt.SetBytes(hash[:])

		// Compare with target
		// If hash is less than target, we found a valid proof-of-work
		if hashInt.Cmp(pow.target) == -1 {
			fmt.Printf("\nBlock mined! Nonce: %d\n", nonce)
			break
		} else {
			nonce++
		}
	}

	// Convert nonce to bytes to match Consensus interface
	return IntToHex(int64(nonce)), hash[:]
}

// Validate verifies the proof-of-work
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	// Convert ValidatorID (which contains the nonce) back to int
	nonce := int(binary.BigEndian.Uint64(pow.block.ValidatorID))

	data := pow.prepareData(nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1
	return isValid
}

// IntToHex converts an int64 to a byte array
func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		panic(err)
	}
	return buff.Bytes()
}
