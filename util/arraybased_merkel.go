//
//  BSD 3-Clause License
//
// Copyright (c) 2018, Sidharth Mishra
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// * Redistributions of source code must retain the above copyright notice, this
//  list of conditions and the following disclaimer.
//
// * Redistributions in binary form must reproduce the above copyright notice,
//  this list of conditions and the following disclaimer in the documentation
//  and/or other materials provided with the distribution.
//
// * Neither the name of the copyright holder nor the names of its
//  contributors may be used to endorse or promote products derived from
//  this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// merkel-sid.go
// @author Sidharth Mishra
// @created Sat Oct 06 2018 20:38:39 GMT-0700 (PDT)
// @last-modified Fri Oct 12 2018 01:48:51 GMT-0700 (PDT)
//

package util

// ----
// An array based implementation of the Merkel tree.
// ----

// ArrayBasedMerkle is a merkle tree that uses an array for storing the entire tree,
// instead of the usual node based tree structure.
//
type ArrayBasedMerkle struct {
	leaves []string
}

// NewArrayBasedMerkle creates a new array based merkle tree instance with 0 nodes
func NewArrayBasedMerkle() *ArrayBasedMerkle {

	return &ArrayBasedMerkle{
		leaves: make([]string, 0),
	}
}

// recursiveConstruct takes a list of leaf nodes and recursive hash + merges them till it
// reaches the merkle root!
//
func recursiveConstruct(leaves []string, roundLeafStartIndex int) []string {
	leafCount := len(leaves)
	index := roundLeafStartIndex

	// base case, when the root node is reached, the diff from leafCount will be 1
	//
	if leafCount-roundLeafStartIndex == 1 {
		return leaves
	}

	// if the round has odd number of nodes, pair the last node with itself to even it out.
	// Since it needs to pair for computing the hash till only one node remains (the root),
	// the number of nodes in each round needs to be even –– expect the root node!
	//
	if leafCount%2 != 0 {
		leaves = append(leaves, leaves[leafCount-1])
	}

	// Compute the hashes by stringing pairs
	//
	for ; index < leafCount; index += 2 {
		leaves = append(leaves, MHash(leaves[index], leaves[index+1]))
	}

	return recursiveConstruct(leaves, index)
}

// ComputeTree computes the tree for the given leaf nodes
// Assumption 1: Since, the name is `hashes`, I'm assuming that these
// entities are already hashed (leaves).
//
func (merkle *ArrayBasedMerkle) ComputeTree(hashes []Hashable) {
	// Add the leaf nodes to the merkle tree
	//
	merkle.leaves = make([]string, len(hashes))
	for index := 0; index < len(merkle.leaves); index++ {
		merkle.leaves[index] = hashes[index].GetHash()
	}

	// compute the tree recursively, the last node is the root
	//
	merkle.leaves = recursiveConstruct(merkle.leaves, 0)
}

// GetRoot gets the hashvalue of the root of the merkle tree. It is the last element of the array.
//
func (merkle *ArrayBasedMerkle) GetRoot() string {

	if len(merkle.leaves) > 0 {

		return merkle.leaves[len(merkle.leaves)-1]
	}

	return ""
}

// GetTree gets the total merkle tree.
//
func (merkle *ArrayBasedMerkle) GetTree() []string {

	return merkle.leaves
}

// GetPath provides the path for the given hashable entity.
//
func (merkle *ArrayBasedMerkle) GetPath(hash Hashable) *MTPath {
	hashHash := hash.GetHash()

	leafIndex := -1

	for index, hashValue := range merkle.leaves {
		if hashValue == hashHash {
			leafIndex = index
			break
		}
	}

	if leafIndex == -1 {
		// no match found in this merkle tree
		//
		return nil
	}

	// If the leafIndex is even, I know that it is going to be the left side of a pair
	// or the root. It is easy to verify that it is a root (in my case, the root is the last index)
	// element of the array –– since the array size if fixed, I don't need to worry about unpredictability.
	//
	// Similarly, if the leadIndex is odd, it is going to be the right side of the pair, and it cannot be
	// the root.
	//
	merklePath := &MTPath{
		Nodes:     make([]string, 0),
		LeafIndex: leafIndex,
	}

	if leafIndex%2 == 0 {
		handleEvenScenario(merkle.leaves, merklePath)
	} else {
		handleOddScenario(merkle.leaves, merklePath)
	}

	return merklePath
}

func handleEvenScenario(treeNodes []string, merklePath *MTPath) {

}

func handleOddScenario(treeNodes []string, merklePath *MTPath) {

}
