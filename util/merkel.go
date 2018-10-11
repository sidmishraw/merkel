package util

import (
	"crypto/sha256"
	"log"
)

// Hashable - anything that can provide it's hash
type Hashable interface {
	GetHash() string
	GetHashBytes() []byte
}

// MerkleTreeI - a merkle tree interface required for constructing and providing verification
type MerkleTreeI interface {
	//API to create a tree from leaf nodes
	ComputeTree(hashes []Hashable)
	GetRoot() string
	GetTree() []string

	// API for verification when the leaf node is known
	GetPath(hash Hashable) *MTPath               // Server needs to provide this
	VerifyPath(hash Hashable, path *MTPath) bool //This is only required by a client but useful for testing

	// API for random verification when the leaf node is uknown
	// (verification of the data to hash used as leaf node is outside this API)
	GetPathByIndex(idx int) *MTPath
}

// MTPath - The merkle tree path
type MTPath struct {
	Nodes     []string `json:"nodes"`
	LeafIndex int      `json:"leaf_index"`
}

// Hash - the hashing used for the merkle tree construction
func Hash(text string) string {
	sha := sha256.New()
	_, err := sha.Write([]byte(text))
	if err != nil {
		log.Fatalln("failed to hash the value!")
	}
	return string(sha.Sum(nil))
}

// MHash - merkle hashing of a pair of child hashes
func MHash(h1 string, h2 string) string {
	return Hash(h1 + h2)
}
