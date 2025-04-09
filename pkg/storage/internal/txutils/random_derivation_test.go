package txutils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandomDerivation(t *testing.T) {
	// when:
	derivation, err := RandomDerivation(16)

	// then:
	require.NoError(t, err)
	require.NotEmpty(t, derivation)
}

func TestRandomDerivationUniqueness(t *testing.T) {
	// when:
	derivation1, err := RandomDerivation(16)
	require.NoError(t, err)

	derivation2, err := RandomDerivation(16)
	require.NoError(t, err)

	// then:
	require.NotEqual(t, derivation1, derivation2)
}

func TestRandomDerivationOnZeroLength(t *testing.T) {
	// when:
	_, err := RandomDerivation(0)

	// then:
	require.Error(t, err)
}

func TestRandomDerivationLengths(t *testing.T) {
	for length := uint64(1); length <= 100; length++ {
		// when:
		derivation, err := RandomDerivation(length)

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
		require.Equal(t, expectedBase64Length, uint64(len(derivation)))
	}
}
