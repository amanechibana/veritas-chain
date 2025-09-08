package main

import (
	"fmt"
	"os"

	"github.com/amanechibana/veritas-chain/blockchain"
)

func main() {
	// Ensure tmp directory exists for identity storage
	if err := os.MkdirAll("./tmp", 0755); err != nil {
		fmt.Printf("Failed to create tmp directory: %v\n", err)
		return
	}

	// Initialize university identity system
	if err := blockchain.InitializeUniversityIdentities(); err != nil {
		fmt.Printf("Failed to initialize university identities: %v\n", err)
		return
	}

	// Create an authorized university identity for blockchain operations
	universityName := "harvard"
	address, err := blockchain.CreateUniversityIdentity(universityName)
	if err != nil {
		fmt.Printf("Failed to create university identity: %v\n", err)
		return
	}
	fmt.Printf("University identity created for %s: Address=%s\n", universityName, address)
	fmt.Println("--------------------------------")

	// Get the university identity for blockchain operations
	universityIdentity, err := blockchain.GetUniversityIdentity(universityName)
	if err != nil {
		fmt.Printf("Failed to get university identity: %v\n", err)
		return
	}

	// Show university identity information
	fmt.Printf("University Identity Details:\n")
	fmt.Printf("  Name: %s\n", universityName)
	fmt.Printf("  Address: %s\n", string(universityIdentity.Address()))
	fmt.Printf("  Public Key X: %x\n", universityIdentity.PrivateKey.PublicKey.X.Bytes())
	fmt.Printf("  Public Key Y: %x\n", universityIdentity.PrivateKey.PublicKey.Y.Bytes())
	fmt.Printf("  Private Key D: %x\n", universityIdentity.PrivateKey.D.Bytes())
	fmt.Println("--------------------------------")

	// Create the genesis block
	blockchain := blockchain.InitBlockchain(universityIdentity)
	fmt.Printf("Genesis block and blockchain created: LastHash=%x\n", blockchain.LastHash)
	fmt.Println("--------------------------------")

	// Create the next block
	block1 := blockchain.AddBlock([]string{"CERT-001", "CERT-002"}, universityIdentity)
	fmt.Printf("Block 1 created: Height=%d, Hash=%x\n", block1.Height, blockchain.LastHash)
	fmt.Printf("Block 1 signature: %x\n", block1.Signature)
	fmt.Printf("Block 1 signature length: %d bytes\n", len(block1.Signature))

	// Verify the signature
	if block1.Verify(universityIdentity.PrivateKey.PublicKey) {
		fmt.Println("Block 1 signature verification: VALID")
	} else {
		fmt.Println("Block 1 signature verification: INVALID")
	}
	fmt.Println("--------------------------------")

	// Create another block
	block2 := blockchain.AddBlock([]string{"CERT-003", "CERT-004", "CERT-005"}, universityIdentity)
	fmt.Printf("Block 2 created: Height=%d, Hash=%x\n", block2.Height, blockchain.LastHash)
	fmt.Printf("Block 2 signature: %x\n", block2.Signature)
	fmt.Printf("Block 2 signature length: %d bytes\n", len(block2.Signature))

	// Verify the signature
	if block2.Verify(universityIdentity.PrivateKey.PublicKey) {
		fmt.Println("Block 2 signature verification: VALID")
	} else {
		fmt.Println("Block 2 signature verification: INVALID")
	}
	fmt.Println("--------------------------------")

	// Print all blocks using iterator
	fmt.Println("All blocks in chain:")
	iter := blockchain.Iterator()
	blockCount := 0
	for {
		block := iter.Next()
		fmt.Printf("Block %d: Height=%d, Hash=%x, Certificates=%d\n",
			blockCount, block.Height, block.Hash, block.GetCertificateCount())

		// Show signature information
		if len(block.Signature) > 0 {
			fmt.Printf("  Signature: %x (length: %d bytes)\n", block.Signature, len(block.Signature))

			// Verify signature with the university identity
			if block.Verify(universityIdentity.PrivateKey.PublicKey) {
				fmt.Printf("	Signature verification: VALID (signed by %s)\n", universityName)
			} else {
				fmt.Printf("	Signature verification: INVALID\n")
			}
		} else {
			fmt.Printf("  ⚠️  No signature found\n")
		}

		blockCount++
		if len(block.PrevHash) == 0 { // reached genesis
			break
		}
	}

	fmt.Println("--------------------------------")

	// Test chain validation
	if err := blockchain.ValidateChain(); err != nil {
		fmt.Printf("Chain validation failed: %v\n", err)
	} else {
		fmt.Println("Chain validation passed!")
	}

	// Close the database
	defer blockchain.Close()
}
