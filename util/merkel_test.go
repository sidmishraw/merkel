package util

import (
	"testing"
)

// tests the sha256 hash function implementation
func testHash(t *testing.T) {
	hashString := Hash("hello")
	if hashString != "2CF24DBA5FB0A30E26E83B2AC5B9E29E1B161E5C1FA7425E73043362938B9824" {
		t.Error("The hash is not correct!")
	}
}
