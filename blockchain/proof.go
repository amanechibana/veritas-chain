package blockchain

import (
	"crypto/ecdsa"
	"errors"
)

type ProofOfAuthority struct {
	Block     *Block
	PublicKey ecdsa.PublicKey
}

// NewProof creates a new ProofOfAuthority verifier
func NewProof(block *Block, publicKey ecdsa.PublicKey) *ProofOfAuthority {
	return &ProofOfAuthority{Block: block, PublicKey: publicKey}
}

func (p *ProofOfAuthority) Run() error {
	// 1) Authorization check (placeholder):
	// In a production system, verify PublicKey against a configured set of authorized signer keys.
	// For now, assume authorization is handled elsewhere and proceed.

	// 2) Verify the block's signature
	if !p.Block.Verify(p.PublicKey) {
		return errors.New("block signature verification failed")
	}
	return nil
}
