package main

import (
	"fmt"

	"github.com/amanechibana/veritas-chain/blockchain"
)

func main() {
	// Create the genesis block
	blockchain := blockchain.InitBlockchain()
	fmt.Printf("Genesis block and blockchain created: LastHash=%x\n", blockchain.LastHash)
	fmt.Println("--------------------------------")

	// Create the next block
	block1 := blockchain.AddBlock([]string{"CERT-001", "CERT-002"})
	fmt.Printf("Block 1 created: Height=%d, Hash=%x\n", block1.Height, blockchain.LastHash)
	fmt.Println("--------------------------------")

	// Create another block
	block2 := blockchain.AddBlock([]string{"CERT-003", "CERT-004", "CERT-005"})
	fmt.Printf("Block 2 created: Height=%d, Hash=%x\n", block2.Height, blockchain.LastHash)
	fmt.Println("--------------------------------")

	// Print all blocks using iterator
	fmt.Println("All blocks in chain:")
	iter := blockchain.Iterator()
	blockCount := 0
	for {
		block := iter.Next()
		fmt.Printf("Block %d: Height=%d, Hash=%x, Certificates=%d\n",
			blockCount, block.Height, block.Hash, block.GetCertificateCount())
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
