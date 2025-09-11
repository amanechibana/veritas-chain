package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/amanechibana/veritas-chain/identity"
)

// Block represents a simple block in the Veritas Chain
type Block struct {
	Timestamp         int64    `json:"timestamp"`
	Hash              []byte   `json:"hash"`
	PrevHash          []byte   `json:"prev_hash"`
	Height            int      `json:"height"`
	CertificateHashes []string `json:"certificate_hashes"` // Hashed certificate IDs
	Signature         []byte   `json:"signature"`          // Digital signature of the block
	MerkleRoot        []byte   `json:"merkle_root"`        // Merkle tree of the block
	UniversityAddress []byte   `json:"university_address"` // University address that created this block
}

// NewBlock creates a new block with certificate hashes
func NewBlock(certificateIDs []string, prevHash []byte, height int, signer identity.Signer) *Block {

	block := &Block{
		Timestamp:         time.Now().Unix(),
		Hash:              []byte{},
		PrevHash:          prevHash,
		Height:            height,
		CertificateHashes: hashCertificateIDs(certificateIDs),
		MerkleRoot:        BuildMerkleTree(certificateIDs).Root.Data,
		UniversityAddress: signer.Address(),
	}

	// Sign the block with the provided signer
	err := block.SignWithSigner(signer)
	if err != nil {
		log.Panic(err)
	}

	pow := NewProof(block, signer.PublicKey())
	err = pow.Run()

	if err != nil {
		log.Panic(err)
	}

	block.Hash = block.CalculateHash()

	return block
}

// Genesis creates the first block in the blockchain
func Genesis(signer identity.Signer) *Block {
	return NewBlock([]string{}, []byte{}, 0, signer)
}

func (b *Block) Serialize() []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	if err := encoder.Encode(b); err != nil {
		log.Panic(err)
	}
	return buffer.Bytes()
}

func Deserialize(data []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))

	if err := decoder.Decode(&block); err != nil {
		log.Panic(err)
	}

	return &block
}

// hashCertificateIDs takes certificate IDs and returns their SHA-256 hashes as hex strings
func hashCertificateIDs(certificateIDs []string) []string {
	var hashes []string
	for _, id := range certificateIDs {
		hash := sha256.Sum256([]byte(id))
		hashes = append(hashes, hex.EncodeToString(hash[:]))
	}
	return hashes
}

func BuildMerkleTree(certificateIDs []string) *MerkleTree {
	return NewMerkleTree(certificateIDs)
}

// CalculateHash calculates the hash of the block (including signature)
func (b *Block) CalculateHash() []byte {
	data := bytes.Join(
		[][]byte{
			b.PrevHash,
			b.HashCertificates(),
			b.MerkleRoot,
			ToHex(int64(b.Timestamp)),
			ToHex(int64(b.Height)),
			b.Signature,
		},
		[]byte{},
	)

	hash := sha256.Sum256(data)
	return hash[:]
}

// CalculateHashForSigning calculates the hash of the block for signing (excluding signature)
func (b *Block) CalculateHashForSigning() []byte {
	data := bytes.Join(
		[][]byte{
			b.PrevHash,
			b.HashCertificates(),
			b.MerkleRoot,
			ToHex(int64(b.Timestamp)),
			ToHex(int64(b.Height)),
		},
		[]byte{},
	)

	hash := sha256.Sum256(data)
	return hash[:]
}

func (b *Block) HashCertificates() []byte {
	var certHashes [][]byte
	for _, certHash := range b.CertificateHashes {
		certHashes = append(certHashes, []byte(certHash))
	}
	return bytes.Join(certHashes, []byte{})
}

func ToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

// VerifyCertificate checks if a certificate ID exists in this block
func (b *Block) VerifyCertificate(certificateID string) bool {
	targetHash := sha256.Sum256([]byte(certificateID))
	targetHashStr := hex.EncodeToString(targetHash[:])

	for _, hash := range b.CertificateHashes {
		if hash == targetHashStr {
			return true
		}
	}
	return false
}

// Sign signs the block with the provided private key
func (block *Block) Sign(privateKey ecdsa.PrivateKey) error {
	// 1. Create a hash of the block data (excluding signature)
	blockHash := block.CalculateHashForSigning()

	// 2. Sign the hash
	r, s, err := ecdsa.Sign(rand.Reader, &privateKey, blockHash)

	if err != nil {
		return err
	}

	// 3. Combine r and s into a single signature
	signature := append(r.Bytes(), s.Bytes()...)

	// 4. Store the signature
	block.Signature = signature

	return nil
}

// SignWithSigner signs the block using the provided signer abstraction.
func (block *Block) SignWithSigner(signer identity.Signer) error {
	blockHash := block.CalculateHashForSigning()
	sig, err := signer.Sign(blockHash)
	if err != nil {
		return err
	}
	block.Signature = sig
	return nil
}

// Verify verifies the block's signature using the provided public key
func (b *Block) Verify(publicKey ecdsa.PublicKey) bool {
	// 1. Check if signature exists
	if len(b.Signature) == 0 {
		return false
	}

	// 2. Create the same hash that was signed
	blockHash := b.CalculateHashForSigning()

	// 3. Split the signature back into r and s components
	sigLen := len(b.Signature)
	if sigLen%2 != 0 {
		return false // Signature should have even length (r + s)
	}

	halfLen := sigLen / 2
	rBytes := b.Signature[:halfLen]
	sBytes := b.Signature[halfLen:]

	// 4. Convert bytes back to big.Int
	r := new(big.Int).SetBytes(rBytes)
	s := new(big.Int).SetBytes(sBytes)

	// 5. Verify the signature
	return ecdsa.Verify(&publicKey, blockHash, r, s)
}

// GetCertificateCount returns the number of certificates in this block
func (b *Block) GetCertificateCount() int {
	return len(b.CertificateHashes)
}

// Validate checks if a block is valid
func (b *Block) Validate() error {
	// Check if hash is correct
	calculatedHash := b.CalculateHash()
	if !bytes.Equal(b.Hash, calculatedHash) {
		return fmt.Errorf("invalid block hash: expected %x, got %x", calculatedHash, b.Hash)
	}

	// Check if height is non-negative
	if b.Height < 0 {
		return fmt.Errorf("invalid block height: %d", b.Height)
	}

	// Check if timestamp is reasonable (not in the future)
	currentTime := time.Now().Unix()
	if b.Timestamp > currentTime+3600 { // Allow 1 hour in the future for clock skew
		return fmt.Errorf("block timestamp is too far in the future: %d", b.Timestamp)
	}

	// Check if certificate hashes are valid hex strings
	for i, certHash := range b.CertificateHashes {
		if len(certHash) != 64 { // SHA-256 hex string should be 64 characters
			return fmt.Errorf("invalid certificate hash at index %d: expected 64 hex chars, got %d", i, len(certHash))
		}
		// Try to decode to verify it's valid hex
		if _, err := hex.DecodeString(certHash); err != nil {
			return fmt.Errorf("invalid certificate hash at index %d: not valid hex", i)
		}
	}

	return nil
}

// GenerateCertificateProof builds a Merkle proof for a given certID using this block's leaves
func (b *Block) GenerateCertificateProof(certID string) (MerkleProof, bool) {
	if len(b.CertificateHashes) == 0 || len(b.MerkleRoot) == 0 {
		return MerkleProof{}, false
	}
	// Build leaves from stored hex hashes (stable order)
	leaves := make([][]byte, 0, len(b.CertificateHashes))
	idx := -1

	target := sha256.Sum256([]byte(certID))
	targetHex := hex.EncodeToString(target[:])

	for i, h := range b.CertificateHashes {
		hb, err := hex.DecodeString(h)
		if err != nil {
			return MerkleProof{}, false
		}
		leaves = append(leaves, hb)
		if h == targetHex {
			idx = i
		}
	}
	if idx == -1 {
		return MerkleProof{}, false
	}

	proof := GenerateProof(leaves, idx)
	return proof, true
}

// VerifyCertificateWithProof verifies a certID against this block's MerkleRoot using a provided proof
func (b *Block) VerifyCertificateWithProof(certID string, proof MerkleProof) bool {
	return VerifyProof([]byte(certID), proof, b.MerkleRoot)
}
