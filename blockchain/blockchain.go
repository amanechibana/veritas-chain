package blockchain

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/amanechibana/veritas-chain/identity"
	"github.com/dgraph-io/badger/v4"
)

// Blockchain is a handle to the on-disk chain state
type Blockchain struct {
	LastHash []byte
	Database *badger.DB
}

type BlockchainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

type BlockchainStats struct {
	BlockCount       int
	CertificateCount int
}

// DBExists checks for Badger MANIFEST to determine if DB exists at given path
func DBExists(dbPath string) bool {
	manifest := filepath.Join(dbPath, "MANIFEST")
	if _, err := os.Stat(manifest); os.IsNotExist(err) {
		return false
	}
	return true
}

func ContinueBlockchain(dbPath string) *Blockchain {
	if !DBExists(dbPath) {
		fmt.Println("No blockchain found")
		runtime.Goexit()
	}

	var lastHash []byte

	opts := badger.DefaultOptions(dbPath)
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	db, err := badger.Open(opts)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		if err != nil {
			log.Panic(err)
		}
		return item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})
	})

	if err != nil {
		log.Panic(err)
	}

	chain := Blockchain{lastHash, db}
	return &chain
}

func InitBlockchain(dbPath string, signer identity.Signer) *Blockchain {
	opts := badger.DefaultOptions(dbPath)
	opts.Dir = dbPath
	opts.ValueDir = dbPath

	// Ensure directory exists
	_ = os.MkdirAll(dbPath, 0o755)

	db, err := badger.Open(opts)
	if err != nil {
		log.Panic(err)
	}

	// Check if blockchain already exists
	if DBExists(dbPath) {
		// Try to load existing blockchain
		var lastHash []byte
		err = db.View(func(txn *badger.Txn) error {
			item, err := txn.Get([]byte("lh"))
			if err != nil {
				return err
			}
			return item.Value(func(val []byte) error {
				lastHash = val
				return nil
			})
		})
		if err != nil {
			// If we can't load the existing blockchain (corrupted/incomplete), recreate
			fmt.Println("Existing blockchain is corrupted or incomplete, recreating...")
			db.Close()
			os.RemoveAll(dbPath)
			_ = os.MkdirAll(dbPath, 0o755)
			db, err = badger.Open(opts)
			if err != nil {
				log.Panic(err)
			}
		} else {
			fmt.Println("Loaded existing blockchain")
			return &Blockchain{lastHash, db}
		}
	}

	// Create new blockchain with genesis block
	var lastHash []byte
	err = db.Update(func(txn *badger.Txn) error {
		genesis := Genesis(signer)
		encodedBlock := genesis.Serialize()
		if err := txn.Set(genesis.Hash, encodedBlock); err != nil {
			return err
		}
		if err := txn.Set([]byte("lh"), genesis.Hash); err != nil {
			return err
		}
		lastHash = genesis.Hash
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	fmt.Println("Created new blockchain with genesis block")
	return &Blockchain{lastHash, db}
}

func (chain *Blockchain) AddBlock(certificateIDs []string, signer identity.Signer) (*Block, error) {
	var lastHash []byte
	var prevBlock *Block

	// Get the previous block to determine height
	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		if err != nil {
			return err
		}
		if err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		}); err != nil {
			return err
		}
		// Get the previous block to calculate height
		item, err = txn.Get(lastHash)
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			prevBlock = Deserialize(val)
			return nil
		})
	})
	if err != nil {
		return nil, err
	}

	// Calculate height: previous block height + 1
	newHeight := prevBlock.Height + 1
	newBlock := NewBlock(certificateIDs, lastHash, newHeight, signer)

	err = chain.Database.Update(func(txn *badger.Txn) error {
		if err := txn.Set(newBlock.Hash, newBlock.Serialize()); err != nil {
			return err
		}
		if err := txn.Set([]byte("lh"), newBlock.Hash); err != nil {
			return err
		}
		chain.LastHash = newBlock.Hash
		return nil
	})
	if err != nil {
		return nil, err
	}
	return newBlock, nil
}

// ValidateChain checks if the entire blockchain is valid
func (bc *Blockchain) ValidateChain() error {
	// Check if blockchain is empty
	if len(bc.LastHash) == 0 {
		return fmt.Errorf("blockchain is empty")
	}

	// Load all blocks into memory for validation (we need to validate in order)
	var blocks []*Block
	currentHash := append([]byte{}, bc.LastHash...)

	// Walk backwards from last block to genesis
	for {
		var data []byte
		err := bc.Database.View(func(txn *badger.Txn) error {
			item, err := txn.Get(currentHash)
			if err != nil {
				return err
			}
			return item.Value(func(val []byte) error {
				data = append([]byte{}, val...)
				return nil
			})
		})
		if err != nil {
			return fmt.Errorf("failed to load block: %v", err)
		}
		block := Deserialize(data)
		blocks = append(blocks, block)
		if len(block.PrevHash) == 0 { // reached genesis
			break
		}
		currentHash = block.PrevHash
	}

	// Reverse to get oldest->newest order
	for i, j := 0, len(blocks)-1; i < j; i, j = i+1, j-1 {
		blocks[i], blocks[j] = blocks[j], blocks[i]
	}

	// Validate genesis block
	genesis := blocks[0]
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
	for i := 1; i < len(blocks); i++ {
		block := blocks[i]
		prevBlock := blocks[i-1]

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
	lastBlock := blocks[len(blocks)-1]
	if !bytes.Equal(bc.LastHash, lastBlock.Hash) {
		return fmt.Errorf("LastHash mismatch: expected %x, got %x", lastBlock.Hash, bc.LastHash)
	}

	return nil
}

func (chain *Blockchain) GetStats() BlockchainStats {
	var blockCount int
	var certificateCount int

	// Count blocks and certificates by iterating through the chain
	iter := chain.Iterator()
	for {
		block := iter.Next()
		blockCount++
		certificateCount += len(block.CertificateHashes)

		// Stop when we reach the genesis block (PrevHash is empty)
		if len(block.PrevHash) == 0 {
			break
		}
	}

	return BlockchainStats{
		BlockCount:       blockCount,
		CertificateCount: certificateCount,
	}
}

// Iterator creates a new blockchain iterator
func (bc *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{
		CurrentHash: bc.LastHash,
		Database:    bc.Database,
	}
}

// Next returns the next block in the chain (newest to oldest)
func (iter *BlockchainIterator) Next() *Block {
	var block *Block
	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			block = Deserialize(val)
			return nil
		})
	})
	if err != nil {
		log.Panic(err)
	}

	// Update CurrentHash to the previous block's hash
	iter.CurrentHash = block.PrevHash
	return block
}

// Close closes the underlying database
func (bc *Blockchain) Close() error {
	if bc.Database != nil {
		return bc.Database.Close()
	}
	return nil
}
