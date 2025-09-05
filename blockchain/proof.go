package blockchain

import "errors"

type ProofOfAuthority struct {
	Block        *Block
	UniversityID string
}

func NewProof(block *Block, UniversityID string) *ProofOfAuthority {
	return &ProofOfAuthority{Block: block, UniversityID: UniversityID}
}

func (p *ProofOfAuthority) Run() error {
	if !p.IsAuthorizedUniversity() {
		return errors.New("university is not authorized")
	}

	return nil
}

func (p *ProofOfAuthority) IsAuthorizedUniversity() bool {
	authorized := []string{"harvard", "mit", "stanford", "yale", "genesis"}
	for _, auth := range authorized {
		if p.UniversityID == auth {
			return true
		}
	}
	return false
}
