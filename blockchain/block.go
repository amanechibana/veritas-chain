package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"time"
)

// Block represents a simple block in the Veritas Chain
type Block struct {
	Timestamp         int64    `json:"timestamp"`
	Hash              []byte   `json:"hash"`
	PrevHash          []byte   `json:"prev_hash"`
	Height            int      `json:"height"`
	CertificateHashes []string `json:"certificate_hashes"` // Hashed certificate IDs
	UniversityAddress string   `json:"university_address"`
}

// NewBlock creates a new block with certificate hashes
func NewBlock(certificateIDs []string, prevHash []byte, height int, universityAddress string) *Block {
	block := &Block{
		Timestamp:         time.Now().Unix(),
		Hash:              []byte{},
		PrevHash:          prevHash,
		Height:            height,
		UniversityAddress: universityAddress,
		CertificateHashes: hashCertificateIDs(certificateIDs),
	}

	pow := NewProof(block, universityAddress)
	err := pow.Run()

	if err != nil {
		log.Panic(err)
	}

	block.Hash = block.CalculateHash()

	return block
}

// Genesis creates the first block in the blockchain
func Genesis() *Block {
	return NewBlock([]string{}, []byte{}, 0, "genesis")
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

// CalculateHash calculates the hash of the block
func (b *Block) CalculateHash() []byte {
	data := bytes.Join(
		[][]byte{
			b.PrevHash,
			b.HashCertificates(),
			ToHex(int64(b.Timestamp)),
			ToHex(int64(b.Height)),
			[]byte(b.UniversityAddress),
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

	// Check if university address is not empty (except for genesis)
	if b.Height > 0 && b.UniversityAddress == "" {
		return fmt.Errorf("university address cannot be empty for non-genesis blocks")
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
