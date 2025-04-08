package txutils

import (
	primitives "github.com/bsv-blockchain/go-sdk/primitives/ec"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLockScriptWithKeyOffsetFromPubKey(t *testing.T) {
	// given:
	generator := NewLockingScriptWithKeyOffset()

	// and:
	generator.offsetPrivGenerator = func() (*primitives.PrivateKey, error) {
		return primitives.PrivateKeyFromHex("8143f5ed6c5b41c3d084d39d49e161d8dde4b50b0685a4e4ac23959d3b8a319b")
	}

	// when:
	lockingScript, keyOffset, err := generator.Generate("02f40c35f798e2ece03ae1ebf749545336db8402eb7e620bfe04d50da8ca8b06cc")

	// then:
	require.NoError(t, err)
	t.Logf("lockingScript: %s", lockingScript)
	t.Logf("keyOffset: %s", keyOffset)
}
