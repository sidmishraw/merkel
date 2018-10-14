package util

import (
	"fmt"
	"strings"
	"testing"
)

// tests the sha256 hash function implementation
func TestHash(t *testing.T) {
	hashString := Hash("hello")
	if hashString != strings.ToLower("2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824") {
		t.Error("The hash is not correct!")
	}
}

func TestRecursiveConstructOddTree(t *testing.T) {

	hashes := [3]string{"a", "b", "c"}

	leaves := recursiveConstruct(hashes[:], 0)

	fmt.Println(leaves)

	if len(leaves) != 7 {
		t.Error("Expect Merkle tree size to be 7!")
	}

	for index, val := range leaves {
		switch index {
		case 0:
			if val != "a" {
				t.Error("Expected 'a'")
			}
		case 1:
			if val != "b" {
				t.Error("Expect 'b'")
			}
		case 2:
			if val != "c" {
				t.Error("Expect 'c'")
			}
		case 3:
			if val != "c" {
				t.Error("Expect 'c'")
			}
		case 4:
			if val != MHash("a", "b") {
				t.Error(`Expect MHash("a", "b")`)
			}
		case 5:
			if val != MHash("c", "c") {
				t.Error(`Expect  MHash("c", "c")`)
			}
		case 6:
			if val != MHash(MHash("a", "b"), MHash("c", "c")) {
				t.Error(`Expect MHash(MHash("a", "b"), MHash("c", "c"))`)
			}
		}
	}
}

func TestRecursiveConstructEvenTree(t *testing.T) {

	hashes := [6]string{"a", "b", "c", "d", "e", "f"}

	leaves := recursiveConstruct(hashes[:], 0)

	fmt.Println(leaves)

	if len(leaves) != 13 {
		t.Error("Expected merkle tree to have 13 nodes")
	}
}

func TestComputeMerklePathFor_b(t *testing.T) {

	hashes := [6]Hashable{H{"a"}, H{"b"}, H{"c"}, H{"d"}, H{"e"}, H{"f"}}
	merkleTree := NewArrayBasedMerkle()

	merkleTree.ComputeTree(hashes[:])

	fmt.Printf("merkle root = %s \n", merkleTree.GetRoot())
	fmt.Printf("merkle leaf count = %d \n", merkleTree.GetNumberOfLeafNodes())
	fmt.Printf("merkle tree  = %v \n", merkleTree.GetTree())
	fmt.Printf("merkle path = %v  \n", merkleTree.GetPath(H{"b"}))

	for i, v := range merkleTree.GetPath(H{"b"}).Nodes {
		switch i {
		case 0:
			// b
			//
			if v != Hash("b") {
				t.Errorf("index 0 - Expected %s but found %s \n", Hash("b"), v)
			}
		case 1:
			// ab
			//
			if v != Hash(Hash("a")+Hash("b")) {
				t.Errorf("index 1 -Expected %s but found %s \n", Hash(Hash("a")+Hash("b")), v)
			}
		case 2:
			// abcd
			//
			if v != Hash(Hash(Hash("a")+Hash("b"))+Hash(Hash("c")+Hash("d"))) {
				t.Errorf("index 2 -Expected %s but found %s \n", Hash(Hash(Hash("a")+Hash("b"))+Hash(Hash("c")+Hash("d"))), v)
			}
		case 3:
			// abcdefef
			//
			if v != Hash(Hash(Hash(Hash("a")+Hash("b"))+Hash(Hash("c")+Hash("d")))+Hash(Hash(Hash("e")+Hash("f"))+Hash(Hash("e")+Hash("f")))) {
				t.Errorf("index 3 -Expected %s but found %s \n", Hash(Hash(Hash(Hash("a")+Hash("b"))+Hash(Hash("c")+Hash("d")))+Hash(Hash(Hash("e")+Hash("f"))+Hash(Hash("e")+Hash("f")))), v)
			}
		}
	}
}

func TestComputeMerklePathFor_c(t *testing.T) {

	hashes := [7]Hashable{H{"a"}, H{"b"}, H{"c"}, H{"d"}, H{"e"}, H{"f"}, H{"g"}}
	merkleTree := NewArrayBasedMerkle()

	merkleTree.ComputeTree(hashes[:])

	fmt.Printf("merkle root = %s \n", merkleTree.GetRoot())
	fmt.Printf("merkle leaf count = %d \n", merkleTree.GetNumberOfLeafNodes())
	fmt.Printf("merkle tree  = %v \n", merkleTree.GetTree())
	fmt.Printf("merkle path = %v  \n", merkleTree.GetPath(H{"c"}))

	for i, v := range merkleTree.GetPath(H{"c"}).Nodes {
		switch i {
		case 0:
			// c
			//
			if v != Hash("c") {
				t.Errorf("index 0 - Expected %s but found %s \n", Hash("c"), v)
			}
		case 1:
			// cd
			//
			if v != Hash(Hash("c")+Hash("d")) {
				t.Errorf("index 1 -Expected %s but found %s \n", Hash(Hash("c")+Hash("d")), v)
			}
		case 2:
			// abcd
			//
			if v != Hash(Hash(Hash("a")+Hash("b"))+Hash(Hash("c")+Hash("d"))) {
				t.Errorf("index 2 -Expected %s but found %s \n", Hash(Hash(Hash("a")+Hash("b"))+Hash(Hash("c")+Hash("d"))), v)
			}
		case 3:
			// abcdefgg
			//
			if v != Hash(Hash(Hash(Hash("a")+Hash("b"))+Hash(Hash("c")+Hash("d")))+Hash(Hash(Hash("e")+Hash("f"))+Hash(Hash("g")+Hash("g")))) {
				t.Errorf("index 3 -Expected %s but found %s \n", Hash(Hash(Hash(Hash("a")+Hash("b"))+Hash(Hash("c")+Hash("d")))+Hash(Hash(Hash("e")+Hash("f"))+Hash(Hash("g")+Hash("g")))), v)
			}
		}
	}

	t.Logf("Found merkle path to be %v", *merkleTree.GetPath(H{"c"}))
}

func TestComputeMerklePathFor_2leaves_b(t *testing.T) {

	hashes := [2]Hashable{H{"a"}, H{"b"}}
	merkleTree := NewArrayBasedMerkle()

	merkleTree.ComputeTree(hashes[:])

	fmt.Printf("merkle root = %s \n", merkleTree.GetRoot())
	fmt.Printf("merkle leaf count = %d \n", merkleTree.GetNumberOfLeafNodes())
	fmt.Printf("merkle tree  = %v \n", merkleTree.GetTree())
	fmt.Printf("merkle path = %v  \n", merkleTree.GetPath(H{"b"}))

	for i, v := range merkleTree.GetPath(H{"b"}).Nodes {
		switch i {
		case 0:
			// b
			//
			if v != Hash("b") {
				t.Errorf("index 0 - Expected %s but found %s \n", Hash("c"), v)
			}
		case 1:
			// ab
			//
			if v != Hash(Hash("a")+Hash("b")) {
				t.Errorf("index 1 -Expected %s but found %s \n", Hash(Hash("a")+Hash("b")), v)
			}
		}
	}

	t.Logf("Found merkle path to be %v", *merkleTree.GetPath(H{"b"}))
}

func TestVerifyPath(t *testing.T) {

	hashes := [2]Hashable{H{"a"}, H{"b"}}
	merkleTree := NewArrayBasedMerkle()

	merkleTree.ComputeTree(hashes[:])

	fmt.Printf("merkle root = %s \n", merkleTree.GetRoot())
	fmt.Printf("merkle leaf count = %d \n", merkleTree.GetNumberOfLeafNodes())
	fmt.Printf("merkle tree  = %v \n", merkleTree.GetTree())
	fmt.Printf("merkle path = %v  \n", merkleTree.GetPath(H{"b"}))

	nodes1 := [2]string{Hash(H{"b"}.v), Hash(Hash(H{"a"}.v) + Hash(H{"b"}.v))}

	nodes2 := [2]string{Hash(H{"c"}.v), Hash(Hash(H{"a"}.v) + Hash(H{"b"}.v))}

	myMerkelPath := &MTPath{
		LeafIndex: 1,
		Nodes:     nodes1[:],
	}
	status := merkleTree.VerifyPath(H{"b"}, myMerkelPath)
	if !status {
		t.Error("Was expecting true got false")
	}

	myMerkelPath2 := &MTPath{
		LeafIndex: -1,
		Nodes:     nodes2[:],
	}
	status2 := merkleTree.VerifyPath(H{"c"}, myMerkelPath2)
	if status2 {
		t.Error("Was expecting false got true")
	}
}

func TestGetPathByIndex_c(t *testing.T) {

	hashes := [7]Hashable{H{"a"}, H{"b"}, H{"c"}, H{"d"}, H{"e"}, H{"f"}, H{"g"}}
	merkleTree := NewArrayBasedMerkle()

	merkleTree.ComputeTree(hashes[:])

	fmt.Printf("merkle root = %s \n", merkleTree.GetRoot())
	fmt.Printf("merkle leaf count = %d \n", merkleTree.GetNumberOfLeafNodes())
	fmt.Printf("merkle tree  = %v \n", merkleTree.GetTree())
	fmt.Printf("merkle path = %v  \n", merkleTree.GetPath(H{"c"}))

	nodes := [4]string{Hash("c"), Hash(Hash("c") + Hash("d")), Hash(Hash(Hash("a")+Hash("b")) + Hash(Hash("c")+Hash("d"))), Hash(Hash(Hash(Hash("a")+Hash("b"))+Hash(Hash("c")+Hash("d"))) + Hash(Hash(Hash("e")+Hash("f"))+Hash(Hash("g")+Hash("g"))))}

	myMerklePath := &MTPath{
		LeafIndex: 2,
		Nodes:     nodes[:],
	}

	t.Logf("my merkle path = %v", *myMerklePath)

	calculatedMPath := merkleTree.GetPathByIndex(2)

	t.Logf("calculated path = %v", *calculatedMPath)

	if myMerklePath.LeafIndex != calculatedMPath.LeafIndex {
		t.Error("The leaf indices don't match")
	}

	for i, v := range myMerklePath.Nodes {
		if calculatedMPath.Nodes[i] != v {
			t.Errorf("The nodes at index %d don't match in values: got %s but expected %s", i, calculatedMPath.Nodes[i], v)
		}
	}

}

// H is a stub used for testing the merkle tree
//
type H struct {
	v string
}

// GetHash makes H conform to Hashable
//
func (h H) GetHash() string {
	return Hash(h.v)
}

// GetHashBytes makes H conform to Hashable
//
func (h H) GetHashBytes() []byte {
	return []byte(h.GetHash())
}
