package identity

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"math/big"
	"os"
)

// Signer defines the minimal interface required to sign blocks and expose identity metadata.
type Signer interface {
	PublicKey() ecdsa.PublicKey
	Address() []byte
	Sign(message []byte) ([]byte, error)
}

// IdentitySigner adapts the existing Identity type to the Signer interface.
type IdentitySigner struct {
	identity *Identity
}

func NewIdentitySigner(id *Identity) *IdentitySigner {
	return &IdentitySigner{identity: id}
}

func (s *IdentitySigner) PublicKey() ecdsa.PublicKey {
	return s.identity.PrivateKey.PublicKey
}

func (s *IdentitySigner) Address() []byte {
	return s.identity.Address()
}

// Sign returns a raw ECDSA signature as r||s bytes for the given message digest.
func (s *IdentitySigner) Sign(message []byte) ([]byte, error) {
	r, ecdsaS, err := ecdsa.Sign(rand.Reader, &s.identity.PrivateKey, message)
	if err != nil {
		return nil, err
	}
	signature := append(r.Bytes(), ecdsaS.Bytes()...)
	return signature, nil
}

// SplitSignatureRS splits a concatenated r||s signature back to big.Int components.
func SplitSignatureRS(sig []byte) (r, s *big.Int) {
	half := len(sig) / 2
	r = new(big.Int).SetBytes(sig[:half])
	s = new(big.Int).SetBytes(sig[half:])
	return r, s
}

// NewP256SignerFromHexD constructs an IdentitySigner from a hex-encoded private scalar D (P-256).
func NewP256SignerFromHexD(hexD string) (*IdentitySigner, error) {
	bytesD, err := hex.DecodeString(hexD)
	if err != nil {
		return nil, err
	}
	curve := elliptic.P256()
	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = curve
	priv.D = new(big.Int).SetBytes(bytesD)
	priv.PublicKey.X, priv.PublicKey.Y = curve.ScalarBaseMult(priv.D.Bytes())

	pubBytes := append(priv.PublicKey.X.Bytes(), priv.PublicKey.Y.Bytes()...)
	id := &Identity{PrivateKey: *priv, PublicKey: pubBytes}
	return &IdentitySigner{identity: id}, nil
}

// LoadSignerFromEnv tries to load a signer from environment variables.
// Supported:
//
//	SIGNER_PRIVATE_KEY_HEX: hex of the P-256 private scalar D
//
// If not set, returns nil, nil indicating caller should fall back to generated signer.
func LoadSignerFromEnv() (*IdentitySigner, error) {
	hexD := os.Getenv("SIGNER_PRIVATE_KEY_HEX")
	if hexD != "" {
		return NewP256SignerFromHexD(hexD)
	}
	return nil, nil
}
