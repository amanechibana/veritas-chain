package identity

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/json"
	"log"
	"math/big"
	"os"

	"github.com/mr-tron/base58"
)

func Base58Encode(input []byte) string {
	return base58.Encode(input)
}

func Base58Decode(input string) []byte {
	decode, err := base58.Decode(input)
	if err != nil {
		log.Panic(err)
	}

	return decode
}

// SerializableIdentity represents an identity that can be JSON serialized
type SerializableIdentity struct {
	PrivateKeyD *big.Int `json:"private_key_d"`
	PublicKeyX  *big.Int `json:"public_key_x"`
	PublicKeyY  *big.Int `json:"public_key_y"`
}

// SerializableIdentities represents identities that can be JSON serialized
type SerializableIdentities struct {
	Identities map[string]*SerializableIdentity `json:"identities"`
}

// ToSerializable converts an Identity to a SerializableIdentity
func (i *Identity) ToSerializable() *SerializableIdentity {
	return &SerializableIdentity{
		PrivateKeyD: i.PrivateKey.D,
		PublicKeyX:  i.PrivateKey.PublicKey.X,
		PublicKeyY:  i.PrivateKey.PublicKey.Y,
	}
}

// FromSerializable converts a SerializableIdentity back to an Identity
func (si *SerializableIdentity) FromSerializable() *Identity {
	curve := elliptic.P256()

	privateKey := &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: curve,
			X:     si.PublicKeyX,
			Y:     si.PublicKeyY,
		},
		D: si.PrivateKeyD,
	}

	publicKey := append(si.PublicKeyX.Bytes(), si.PublicKeyY.Bytes()...)

	return &Identity{
		PrivateKey: *privateKey,
		PublicKey:  publicKey,
	}
}

// SaveIdentitiesToFile saves identities to a JSON file
func SaveIdentitiesToFile(identities map[string]*Identity, filename string) error {
	// Convert to serializable format
	serializable := &SerializableIdentities{
		Identities: make(map[string]*SerializableIdentity),
	}

	for address, identity := range identities {
		serializable.Identities[address] = identity.ToSerializable()
	}

	// Encode to JSON
	jsonData, err := json.MarshalIndent(serializable, "", "  ")
	if err != nil {
		return err
	}

	// Write to file
	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}

// LoadIdentitiesFromFile loads identities from a JSON file
func LoadIdentitiesFromFile(filename string) (map[string]*Identity, error) {
	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, err
	}

	fileContent, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Try to decode as JSON
	var serializable SerializableIdentities
	err = json.Unmarshal(fileContent, &serializable)
	if err != nil {
		return nil, err
	}

	// Convert back to Identity format
	identities := make(map[string]*Identity)
	for address, serializableIdentity := range serializable.Identities {
		identities[address] = serializableIdentity.FromSerializable()
	}

	return identities, nil
}
