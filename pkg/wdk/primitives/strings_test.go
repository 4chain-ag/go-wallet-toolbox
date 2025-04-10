package primitives_test

import (
	"testing"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk/primitives"
	"github.com/go-softwarelab/common/pkg/seq"
	"github.com/stretchr/testify/require"
)

func TestString5to2000Bytes(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		// when:
		err := primitives.String5to2000Bytes("valid string").Validate()

		// then:
		require.NoError(t, err)
	})

	errorcases := map[string]struct {
		value primitives.String5to2000Bytes
	}{
		"too short": {
			value: "1234",
		},
		"too long": {
			value: primitives.String5to2000Bytes(seq.Collect(seq.Repeat('a', 2001))),
		},
	}
	for name, test := range errorcases {
		t.Run(name, func(t *testing.T) {
			// when:
			err := test.value.Validate()

			// then:
			require.Error(t, err)
		})
	}
}

func TestBase64String(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		// when:
		err := primitives.Base64String("SGVsbG8gV29ybGQ=").Validate()

		// then:
		require.NoError(t, err)
	})

	errorcases := map[string]struct {
		value primitives.Base64String
	}{
		"invalid length": {
			value: "SGVsbG8gV29ybGQ",
		},
		"invalid padding": {
			value: "SGVsbG8gV29ybGQ===",
		},
		"invalid characters": {
			value: "SGVsbG8!V29ybGQ=",
		},
	}
	for name, test := range errorcases {
		t.Run(name, func(t *testing.T) {
			// when:
			err := test.value.Validate()

			// then:
			require.Error(t, err)
		})
	}
}

func TestStringUnder50Bytes(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		// when:
		err := primitives.StringUnder50Bytes("valid string").Validate()

		// then:
		require.NoError(t, err)
	})

	errorcases := map[string]struct {
		value primitives.StringUnder50Bytes
	}{
		"empty": {
			value: "",
		},
		"too long": {
			value: primitives.StringUnder50Bytes(seq.Collect(seq.Repeat('a', 51))),
		},
	}
	for name, test := range errorcases {
		t.Run(name, func(t *testing.T) {
			// when:
			err := test.value.Validate()

			// then:
			require.Error(t, err)
		})
	}
}

func TestStringUnder300(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		// when:
		err := primitives.StringUnder300("valid string").Validate()

		// then:
		require.NoError(t, err)
	})

	errorcases := map[string]struct {
		value primitives.StringUnder300
	}{
		"empty": {
			value: "",
		},
		"too long": {
			value: primitives.StringUnder300(seq.Collect(seq.Repeat('a', 301))),
		},
	}
	for name, test := range errorcases {
		t.Run(name, func(t *testing.T) {
			// when:
			err := test.value.Validate()

			// then:
			require.Error(t, err)
		})
	}
}

func TestHexString(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		// when:
		err := primitives.HexString("48656c6c6f").Validate()

		// then:
		require.NoError(t, err)
	})

	errorcases := map[string]struct {
		value primitives.HexString
	}{
		"odd length": {
			value: "48656c6",
		},
		"invalid characters": {
			value: "48656g6c6c",
		},
	}
	for name, test := range errorcases {
		t.Run(name, func(t *testing.T) {
			// when:
			err := test.value.Validate()

			// then:
			require.Error(t, err)
		})
	}
}

func TestPubKeyHex(t *testing.T) {
	t.Run("valid compressed", func(t *testing.T) {
		// when:
		hex := primitives.PubKeyHex(seq.Collect(seq.Repeat('a', 66)))
		err := hex.Validate()

		// then:
		require.NoError(t, err)
	})

	t.Run("valid uncompressed", func(t *testing.T) {
		// when:
		hex := primitives.PubKeyHex(seq.Collect(seq.Repeat('a', 130)))
		err := hex.Validate()

		// then:
		require.NoError(t, err)
	})

	errorcases := map[string]struct {
		value primitives.PubKeyHex
	}{
		"invalid length": {
			value: primitives.PubKeyHex(seq.Collect(seq.Repeat('a', 64))),
		},
		"invalid hex characters": {
			value: primitives.PubKeyHex(seq.Collect(seq.Repeat('z', 66))),
		},
	}
	for name, test := range errorcases {
		t.Run(name, func(t *testing.T) {
			// when:
			err := test.value.Validate()

			// then:
			require.Error(t, err)
		})
	}
}
