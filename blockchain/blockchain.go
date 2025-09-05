package blockchain

type Blockchain struct {
	LastHash []byte
	Blocks   []*Block
}

func (bc *Blockchain) AddBlock(certificateIDs []string) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(certificateIDs, prevBlock.Hash, len(bc.Blocks)+1, "harvard")
	bc.Blocks = append(bc.Blocks, newBlock)
	bc.LastHash = newBlock.Hash
}

func InitBlockchain() *Blockchain {
	return &Blockchain{Blocks: []*Block{Genesis()}}
}
