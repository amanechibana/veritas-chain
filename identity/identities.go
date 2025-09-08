package identity

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"os"
)

const identityFile = "./tmp/identities.data"

type Identities struct {
	Identities map[string]*Identity
}

func CreateIdentities() (*Identities, error) {
	identities := Identities{}
	identities.Identities = make(map[string]*Identity)

	// Try to load existing identities, but don't fail if file doesn't exist
	err := identities.LoadFile()
	if err != nil {
		// If file doesn't exist, that's okay - we'll start with empty identities
		// Only return error if it's a different type of error
		if _, ok := err.(*os.PathError); !ok {
			return &identities, err
		}
		// PathError means file doesn't exist, which is fine for first run
		err = nil
	}

	return &identities, err
}

func (ws *Identities) AddIdentity() string {
	identity := MakeIdentity()
	address := string(identity.Address())

	ws.Identities[address] = identity

	return address
}

func (ws *Identities) GetAllAddresses() []string {
	var addresses []string

	for address := range ws.Identities {
		addresses = append(addresses, address)
	}

	return addresses
}

func (ws Identities) GetIdentity(address string) Identity {
	return *ws.Identities[address]
}

func (ws *Identities) LoadFile() error {
	// Try to load from JSON file first (new format)
	identities, err := LoadIdentitiesFromFile(identityFile)
	if err == nil {
		ws.Identities = identities
		return nil
	}

	// Fallback to gob format (legacy) - only if file exists but JSON failed
	if _, statErr := os.Stat(identityFile); statErr == nil {
		var wallets Identities
		fileContent, readErr := os.ReadFile(identityFile)
		if readErr != nil {
			return readErr
		}

		gob.Register(elliptic.P256())
		decoder := gob.NewDecoder(bytes.NewReader(fileContent))
		gobErr := decoder.Decode(&wallets)
		if gobErr != nil {
			return gobErr
		}

		ws.Identities = wallets.Identities
		return nil
	}

	// File doesn't exist, which is fine for first run
	return err
}

func (ws *Identities) SaveFile() {
	err := SaveIdentitiesToFile(ws.Identities, identityFile)
	if err != nil {
		fmt.Printf("Failed to save identities: %v\n", err)
		return
	}
	fmt.Println("Identities saved successfully")
}
