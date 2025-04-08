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
