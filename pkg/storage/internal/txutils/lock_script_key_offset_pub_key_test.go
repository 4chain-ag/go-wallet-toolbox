package txutils

import (
	primitives "github.com/bsv-blockchain/go-sdk/primitives/ec"
	"github.com/bsv-blockchain/go-sdk/script"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLockScriptWithKeyOffsetFromPubKey(t *testing.T) {
	// given:
	generator := NewLockingScriptWithKeyOffset()

	// and:
	offsetPrivKey := "L1Yz9NvuvDru3Z7f8Kh8CK34U6UoKa6niynPEvqgHpm3AaNZig6z"
	pubKey := "02f40c35f798e2ece03ae1ebf749545336db8402eb7e620bfe04d50da8ca8b06cc"

	// and:
	generator.offsetPrivGenerator = func() (*primitives.PrivateKey, error) {
		return primitives.PrivateKeyFromWif(offsetPrivKey)
	}

	// when:
	lockingScript, keyOffset, err := generator.Generate(pubKey)

	// then:
	require.NoError(t, err)

	// NOTE: these values are cross-checked with the original/TS code
	assert.Equal(t, "76a914b95556849619ac10419b6a591b6920cb6deef47b88ac", lockingScript)
	assert.Equal(t, offsetPrivKey, keyOffset)
}

func TestLockScriptWithKeyOffset_Uniqueness(t *testing.T) {
	// given:
	generator := NewLockingScriptWithKeyOffset()

	// and:
	pubKey := "02f40c35f798e2ece03ae1ebf749545336db8402eb7e620bfe04d50da8ca8b06cc"

	lockingScripts := make(map[string]struct{})
	keyOffsets := make(map[string]struct{})

	iterations := 100

	// when:
	for range iterations {
		lockingScript, keyOffset, err := generator.Generate(pubKey)
		require.NoError(t, err)

		lockingScripts[lockingScript] = struct{}{}
		keyOffsets[keyOffset] = struct{}{}

		_, err = script.DecodeScriptHex(lockingScript)
		require.NoError(t, err)

		_, err = primitives.PrivateKeyFromWif(keyOffset)
		require.NoError(t, err)
	}

	// then:
	assert.Equal(t, iterations, len(lockingScripts), "Locking script should be unique")
	assert.Equal(t, iterations, len(keyOffsets), "Key offset should be unique")
}
