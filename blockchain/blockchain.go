package blockchain

import (
	"bytes"
	"fmt"
)

type Blockchain struct {
	LastHash []byte
	Blocks   []*Block
}

func (bc *Blockchain) AddBlock(certificateIDs []string) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(certificateIDs, prevBlock.Hash, len(bc.Blocks), "harvard")
	bc.Blocks = append(bc.Blocks, newBlock)
	bc.LastHash = newBlock.Hash
}

func InitBlockchain() *Blockchain {
	genesis := Genesis()
	return &Blockchain{
		LastHash: genesis.Hash,
		Blocks:   []*Block{genesis},
	}
}

// ValidateChain checks if the entire blockchain is valid
func (bc *Blockchain) ValidateChain() error {
	if len(bc.Blocks) == 0 {
		return fmt.Errorf("blockchain is empty")
	}

	// Validate genesis block
	genesis := bc.Blocks[0]
	if genesis.Height != 0 {
		return fmt.Errorf("first block must be genesis block with height 0, got %d", genesis.Height)
	}
	if len(genesis.PrevHash) != 0 {
		return fmt.Errorf("genesis block should have empty PrevHash")
	}
	if err := genesis.Validate(); err != nil {
		return fmt.Errorf("genesis block validation failed: %v", err)
	}

	// Validate all other blocks
	for i := 1; i < len(bc.Blocks); i++ {
		block := bc.Blocks[i]
		prevBlock := bc.Blocks[i-1]

		// Validate individual block
		if err := block.Validate(); err != nil {
			return fmt.Errorf("block %d validation failed: %v", i, err)
		}

		// Check height sequence
		if block.Height != i {
			return fmt.Errorf("block %d has incorrect height: expected %d, got %d", i, i, block.Height)
		}

		// Check previous hash linking
		if !bytes.Equal(block.PrevHash, prevBlock.Hash) {
			return fmt.Errorf("block %d has incorrect PrevHash: expected %x, got %x",
				i, prevBlock.Hash, block.PrevHash)
		}

		// Check timestamp ordering (blocks should be in chronological order)
		if block.Timestamp < prevBlock.Timestamp {
			return fmt.Errorf("block %d timestamp (%d) is before previous block timestamp (%d)",
				i, block.Timestamp, prevBlock.Timestamp)
		}
	}

	// Check if LastHash matches the last block
	lastBlock := bc.Blocks[len(bc.Blocks)-1]
	if !bytes.Equal(bc.LastHash, lastBlock.Hash) {
		return fmt.Errorf("LastHash mismatch: expected %x, got %x", lastBlock.Hash, bc.LastHash)
	}

	return nil
}
