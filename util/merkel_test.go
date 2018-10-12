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
