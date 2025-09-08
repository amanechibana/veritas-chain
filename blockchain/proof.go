package blockchain

import (
	"crypto/ecdsa"
	"errors"

	"github.com/amanechibana/veritas-chain/identity"
)

type ProofOfAuthority struct {
	Block     *Block
	PublicKey ecdsa.PublicKey
}

// AuthorizedUniversities contains the addresses of authorized universities
// These addresses correspond to identities in the identity system
var AuthorizedUniversities = map[string]string{
	"harvard":  "", // Will be set when university identities are created
	"mit":      "", // Will be set when university identities are created
	"stanford": "", // Will be set when university identities are created
	"yale":     "", // Will be set when university identities are created
}

// UniversityIdentities holds the identity system for managing university identities
var UniversityIdentities *identity.Identities

// InitializeUniversityIdentities initializes the university identity system
func InitializeUniversityIdentities() error {
	var err error
	UniversityIdentities, err = identity.CreateIdentities()
	if err != nil {
		return err
	}
	return nil
}

// CreateUniversityIdentity creates a new identity for a university
func CreateUniversityIdentity(universityName string) (string, error) {
	if UniversityIdentities == nil {
		if err := InitializeUniversityIdentities(); err != nil {
			return "", err
		}
	}

	address := UniversityIdentities.AddIdentity()
	AuthorizedUniversities[universityName] = address
	UniversityIdentities.SaveFile()

	return address, nil
}

func NewProof(block *Block, publicKey ecdsa.PublicKey) *ProofOfAuthority {
	return &ProofOfAuthority{Block: block, PublicKey: publicKey}
}

func (p *ProofOfAuthority) Run() error {
	// 1. Check if the public key is authorized
	if !p.IsAuthorizedUniversity() {
		return errors.New("university public key is not authorized")
	}

	// 2. Verify the block's signature
	if !p.Block.Verify(p.PublicKey) {
		return errors.New("block signature verification failed")
	}

	return nil
}

func (p *ProofOfAuthority) IsAuthorizedUniversity() bool {
	if UniversityIdentities == nil {
		return false
	}

	// Check if the public key matches any of the authorized universities
	for _, address := range AuthorizedUniversities {
		if address == "" {
			continue // Skip uninitialized universities
		}

		universityIdentity := UniversityIdentities.GetIdentity(address)
		if p.PublicKeysEqual(p.PublicKey, universityIdentity.PrivateKey.PublicKey) {
			return true
		}
	}
	return false
}

// PublicKeysEqual compares two ECDSA public keys for equality
func (p *ProofOfAuthority) PublicKeysEqual(key1, key2 ecdsa.PublicKey) bool {
	return key1.Curve == key2.Curve &&
		key1.X.Cmp(key2.X) == 0 &&
		key1.Y.Cmp(key2.Y) == 0
}

// GetUniversityName returns the name of the university for a given public key
func (p *ProofOfAuthority) GetUniversityName() string {
	if UniversityIdentities == nil {
		return "unknown"
	}

	for name, address := range AuthorizedUniversities {
		if address == "" {
			continue // Skip uninitialized universities
		}

		universityIdentity := UniversityIdentities.GetIdentity(address)
		if p.PublicKeysEqual(p.PublicKey, universityIdentity.PrivateKey.PublicKey) {
			return name
		}
	}
	return "unknown"
}

// GetUniversityIdentity returns the identity for a given university name
func GetUniversityIdentity(universityName string) (*identity.Identity, error) {
	if UniversityIdentities == nil {
		return nil, errors.New("university identities not initialized")
	}

	address, exists := AuthorizedUniversities[universityName]
	if !exists || address == "" {
		return nil, errors.New("university not found or not initialized")
	}

	universityIdentity := UniversityIdentities.GetIdentity(address)
	return &universityIdentity, nil
}

// AddAuthorizedUniversity adds a new university to the authorized list
func AddAuthorizedUniversity(name string) (string, error) {
	return CreateUniversityIdentity(name)
}

// RemoveAuthorizedUniversity removes a university from the authorized list
func RemoveAuthorizedUniversity(name string) {
	delete(AuthorizedUniversities, name)
}

// GetAuthorizedUniversities returns a list of authorized university names
func GetAuthorizedUniversities() []string {
	var universities []string
	for name, address := range AuthorizedUniversities {
		if address != "" { // Only include initialized universities
			universities = append(universities, name)
		}
	}
	return universities
}

// GetUniversityAddress returns the address for a given university name
func GetUniversityAddress(universityName string) (string, error) {
	address, exists := AuthorizedUniversities[universityName]
	if !exists || address == "" {
		return "", errors.New("university not found or not initialized")
	}
	return address, nil
}
