package blockchain

import (
	"bytes"
	"crypto/sha256"
	"log"
)

type MerkleTree struct {
	Root *MerkleNode
}

type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Data  []byte
}

type MerkleProof struct {
	Siblings   [][]byte
	Directions []bool
}

func NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	node := &MerkleNode{}

	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		node.Data = hash[:]
	} else {
		prevHashes := append(left.Data, right.Data...)
		hash := sha256.Sum256(prevHashes)
		node.Data = hash[:]
	}

	node.Left = left
	node.Right = right

	return node
}

func NewMerkleTree(certificateIDs []string) *MerkleTree {
	var nodes []MerkleNode

	if len(certificateIDs) == 0 {
		log.Panic("No certificate IDs provided")
	}

	if len(certificateIDs)%2 != 0 {
		certificateIDs = append(certificateIDs, certificateIDs[len(certificateIDs)-1])
	}

	for _, data := range certificateIDs {
		nodes = append(nodes, *NewMerkleNode(nil, nil, []byte(data)))
	}

	for len(nodes) > 1 {
		var level []MerkleNode
		for j := 0; j < len(nodes); j += 2 {
			left := &nodes[j]
			right := left
			if j+1 < len(nodes) {
				right = &nodes[j+1]
			}
			parent := NewMerkleNode(left, right, nil)
			level = append(level, *parent)
		}
		nodes = level
	}

	return &MerkleTree{&nodes[0]}
}

func GenerateProof(leaves [][]byte, leafIndex int) MerkleProof {
	if len(leaves) == 0 || leafIndex < 0 || leafIndex >= len(leaves) {
		return MerkleProof{}
	}
	level := make([][]byte, len(leaves))
	copy(level, leaves)

	idx := leafIndex
	var siblings [][]byte
	var dirs []bool

	for len(level) > 1 {
		if len(level)%2 == 1 {
			level = append(level, level[len(level)-1])
		}
		var next [][]byte
		for i := 0; i < len(level); i += 2 {
			left, right := level[i], level[i+1]

			if i == idx || i+1 == idx {
				if i == idx {
					siblings = append(siblings, right)
					dirs = append(dirs, true) // sibling on right
				} else {
					siblings = append(siblings, left)
					dirs = append(dirs, false) // sibling on left
				}
				idx = len(next)
			}

			parent := sha256.Sum256(append(left, right...))
			next = append(next, parent[:])
		}
		level = next
	}
	return MerkleProof{Siblings: siblings, Directions: dirs}
}

func VerifyProof(leafData []byte, proof MerkleProof, root []byte) bool {
	h := sha256.Sum256(leafData)
	curr := h[:]
	for i := range proof.Siblings {
		sib := proof.Siblings[i]
		if proof.Directions[i] {
			sum := sha256.Sum256(append(curr, sib...)) // curr || sib
			curr = sum[:]
		} else {
			sum := sha256.Sum256(append(sib, curr...)) // sib || curr
			curr = sum[:]
		}
	}
	return bytes.Equal(curr, root)
}
