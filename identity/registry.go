package identity

import (
	"encoding/json"
	"errors"
	"os"
)

// AuthorizedSigners represents a mapping of university/organization name to address.
type AuthorizedSigners map[string]string

// LoadAuthorizedSigners loads a JSON file mapping names to addresses.
// Example file content:
//
//	{
//	  "harvard": "1HW5zUskrWwHW7owJExCd5uDMb8Qm8foUG",
//	  "mit": "..."
//	}
func LoadAuthorizedSigners(path string) (AuthorizedSigners, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var m AuthorizedSigners
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return m, nil
}

// ResolveNameByAddress returns the first name whose address matches the provided address.
func (a AuthorizedSigners) ResolveNameByAddress(address string) (string, error) {
	for name, addr := range a {
		if addr == address {
			return name, nil
		}
	}
	return "", errors.New("address not found in authorized signers")
}
