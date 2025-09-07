package main

import (
	"fmt"

	"github.com/amanechibana/veritas-chain/blockchain"
)

func main() {
	// Create the genesis block
	blockchain := blockchain.InitBlockchain()
	fmt.Printf("Genesis block and blockchain created: %+v\n", blockchain)
	fmt.Println("--------------------------------")

	// Create the next block
	blockchain.AddBlock([]string{"CERT-001", "CERT-002"})

	fmt.Printf("Block 1 created: %x\n", blockchain.LastHash)
	fmt.Println("--------------------------------")
	fmt.Printf("Block 1 created: %+v\n", blockchain.Blocks[1])

	// Create another block
	blockchain.AddBlock([]string{"CERT-003", "CERT-004", "CERT-005"})
	fmt.Println("--------------------------------")
	fmt.Printf("Block 2 created: %x\n", blockchain.LastHash)
	fmt.Printf("Block 2 has %d certificates\n", blockchain.Blocks[2].GetCertificateCount())

	fmt.Println("--------------------------------")

	// Test chain validation
	if err := blockchain.ValidateChain(); err != nil {
		fmt.Printf("Chain validation failed: %v\n", err)
	} else {
		fmt.Println("Chain validation passed!")
	}

	fmt.Printf("Blockchain: %+v\n", blockchain)
}
