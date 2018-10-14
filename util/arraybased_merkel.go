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

// ArrayBasedAutoEvenerMerkle is a merkle tree that uses an array for storing the entire tree,
// instead of the usual node based tree structure. It also automatically concatenates odd'th node with itself,
// when computing the hashes for next level.
//
type ArrayBasedAutoEvenerMerkle struct {
	nodes []string
	// the number of leaf nodes is fixed at the time of construction of this merkle
	// tree, and it is always even
	nbrLeafNodes uint
}

// NewArrayBasedMerkle creates a new array based merkle tree instance with 0 nodes
func NewArrayBasedMerkle() *ArrayBasedAutoEvenerMerkle {

	return &ArrayBasedAutoEvenerMerkle{
		nodes: make([]string, 0),
	}
}

// GetNumberOfLeafNodes gets the number of leaf nodes of the merkle tree.
//
func (merkle *ArrayBasedAutoEvenerMerkle) GetNumberOfLeafNodes() uint {

	if merkle.nbrLeafNodes%2 != 0 {
		return merkle.nbrLeafNodes + 1
	}

	return merkle.nbrLeafNodes
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
func (merkle *ArrayBasedAutoEvenerMerkle) ComputeTree(hashes []Hashable) {
	// Add the leaf nodes to the merkle tree
	//
	merkle.nodes = make([]string, len(hashes))
	merkle.nbrLeafNodes = uint(len(hashes))
	for index := 0; index < len(merkle.nodes); index++ {
		merkle.nodes[index] = hashes[index].GetHash()
	}

	// compute the tree recursively, the last node is the root
	//
	merkle.nodes = recursiveConstruct(merkle.nodes, 0)
}

// GetRoot gets the hashvalue of the root of the merkle tree. It is the last element of the array.
//
func (merkle *ArrayBasedAutoEvenerMerkle) GetRoot() string {

	if len(merkle.nodes) > 0 {

		return merkle.nodes[len(merkle.nodes)-1]
	}

	return ""
}

// GetTree gets the total merkle tree.
//
func (merkle *ArrayBasedAutoEvenerMerkle) GetTree() []string {

	return merkle.nodes
}

// GetPath provides the merkle path for the given hashable entity. Returns `nil` if the
// hash doesn't exist in the merkle tree.
//
func (merkle *ArrayBasedAutoEvenerMerkle) GetPath(hash Hashable) *MTPath {
	hashHash := hash.GetHash()

	leafIndex := -1

	nbrLeafNodes := Evenify(merkle.nbrLeafNodes)

	// find the leaf index in the merkle tree
	//
	for index := 0; index < int(nbrLeafNodes); index++ {
		if merkle.nodes[index] == hashHash {
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

	// Compute the merkle path for the leafIndex.
	//
	merkle.computeMerklePath(uint(leafIndex), &merklePath.Nodes, 0, 0)

	return merklePath
}

// computeMerklePath computes the merkle path for the given leaf-node. The `pathNodes` is the reference
// to the array of nodes that make up the path -- all the way upto the root from the leaf node -- the `prevLevelLeafCount`
// is the number of leaf nodes from the previous level, and `totalLeavesConsumed` refers to the total number
// of nodes preceeding the leafNode in the current level or iteration.
//
// Disclaimer: this is a recursive approach!
//
func (merkle *ArrayBasedAutoEvenerMerkle) computeMerklePath(leafIndex uint, pathNodes *[]string, prevLevelLeafCount uint, totalLeavesConsumed uint) {

	*pathNodes = append(*pathNodes, merkle.nodes[leafIndex])

	if leafIndex == uint(len(merkle.nodes)-1) {
		// reached the root and the path is complete
		//
		return
	}

	// eg - 6 for index 7
	// eg - 0 for index 1
	//
	leafIndex = Leftify(leafIndex)

	// eg - (8 - 6) / 2 = 1 (0-indexed)
	// eg - (10 - 6 - 4) / 2 = 0 (0 - indexed)
	//
	bucketNbr := (leafIndex - totalLeavesConsumed) / 2

	var thisLevelLeafCount uint
	if prevLevelLeafCount == 0 {
		thisLevelLeafCount = Evenify(merkle.nbrLeafNodes) // eg - 6
	} else {
		thisLevelLeafCount = Evenify(prevLevelLeafCount / 2) // eg - 6/2 + 1 = 4
	}

	// eg - 6 + 4 = 10
	//
	nextLevelLeafIndex := leafIndex + (thisLevelLeafCount - bucketNbr)

	merkle.computeMerklePath(nextLevelLeafIndex, pathNodes, thisLevelLeafCount, prevLevelLeafCount+thisLevelLeafCount)
}

// Evenify evens out a number, if it is odd, it returns the next even number
// else just returns the same number as is.
//
func Evenify(n uint) uint {
	if n%2 == 0 {
		return n
	}
	return n + 1
}

// Leftify computes the left node index of the pair of the given node.
//
func Leftify(n uint) uint {
	if n%2 == 0 {
		return n
	}
	return n - 1
}

// VerifyPath verifies a given leaf node and the path by computing the merkle-path for the leaf node and
// diffing it against the given path.
//
func (merkle *ArrayBasedAutoEvenerMerkle) VerifyPath(hash Hashable, path *MTPath) bool {
	calculatedPath := merkle.GetPath(hash)
	if calculatedPath == nil {
		return false
	}
	if calculatedPath.LeafIndex != path.LeafIndex {
		return false
	}
	for i, v := range calculatedPath.Nodes {
		if path.Nodes[i] != v {
			return false
		}
	}
	return true
}

// GetPathByIndex gets the merkle path for the leaf node given the index.
// Disclaimer: the data verification of the leaf node is out of the scope for this implementation.
//
func (merkle *ArrayBasedAutoEvenerMerkle) GetPathByIndex(idx int) *MTPath {
	if idx >= int(Evenify(merkle.nbrLeafNodes)) || idx < 0 {
		return nil
	}
	merklePath := &MTPath{
		Nodes:     make([]string, 0),
		LeafIndex: idx,
	}
	// since the leaf index is know here, we simply go ahead and compute the
	// merkle path for it.
	//
	// Data verification is out of the scope --
	//
	merkle.computeMerklePath(uint(idx), &merklePath.Nodes, 0, 0)
	return merklePath
}
