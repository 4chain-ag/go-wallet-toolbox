package randomizer

import (
	"encoding/base64"
	"fmt"
	"slices"

	"github.com/go-softwarelab/common/pkg/must"
)

// TestRandomizer is a test implementation of the Randomizer interface.
// It provides deterministic outputs for testing purposes.
type TestRandomizer struct{}

// NewTestRandomizer creates and returns a new instance of TestRandomizer.
func NewTestRandomizer() *TestRandomizer {
	return &TestRandomizer{}
}

// Base64 generates a deterministic base64-encoded string of the specified length.
// The content of the string is a repeated sequence of the character 'a'.
func (t *TestRandomizer) Base64(length uint64) (string, error) {
	if length == 0 {
		return "", fmt.Errorf("length cannot be zero")
	}

	randomBytes := slices.Repeat([]byte{'a'}, must.ConvertToIntFromUnsigned(length))
	return base64.StdEncoding.EncodeToString(randomBytes), nil
}

// Shuffle performs a deterministic shuffle operation on a slice of size n.
// It calls the provided swap function twice for each pair of indices to preserve the original order.
func (t *TestRandomizer) Shuffle(n int, swap func(i int, j int)) {
	for i := 0; i < n-1; i++ {
		swap(i, i+1)
		swap(i, i+1)
	}
}

// Uint64 returns a deterministic uint64 value, which is always 0 in this implementation.
func (t *TestRandomizer) Uint64(max uint64) uint64 {
	return 0
}
