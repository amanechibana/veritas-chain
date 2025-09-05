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

	fmt.Println("--------------------------------")
	fmt.Printf("Blockchain: %+v\n", blockchain)
}
