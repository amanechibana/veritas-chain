package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
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

// Genesis creates the first block in the blockchain
func Genesis() *Block {
	return NewBlock([]string{}, []byte{}, 0, "genesis")
}
