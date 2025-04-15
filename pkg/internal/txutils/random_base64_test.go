package txutils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandomBase64(t *testing.T) {
	// when:
	randomized, err := RandomBase64(16)

	// then:
	require.NoError(t, err)
	require.NotEmpty(t, randomized)
}

func TestRandomBase64Uniqueness(t *testing.T) {
	// when:
	randomized1, err := RandomBase64(16)
	require.NoError(t, err)

	randomized2, err := RandomBase64(16)
	require.NoError(t, err)

	// then:
	require.NotEqual(t, randomized1, randomized2)
}

func TestRandomBase64OnZeroLength(t *testing.T) {
	// when:
	_, err := RandomBase64(0)

	// then:
	require.Error(t, err)
}

func TestRandomBase64Lengths(t *testing.T) {
	for length := uint64(1); length <= 100; length++ {
		// when:
		randomized, err := RandomBase64(length)

		// then:
		require.NoError(t, err)

		// NOTE: Base64 encoding adds padding, so the length sequence is as follows:
		// Length -> Base64 Length
		// 1 -> 4
		// 2 -> 4
		// 3 -> 4
		// 4 -> 8
		// 5 -> 8
		// 6 -> 8
		// 7 -> 12
		// ...
		expectedBase64Length := ((length-1)/3 + 1) * 4
		require.Equal(t, expectedBase64Length, uint64(len(randomized)))
	}
}
