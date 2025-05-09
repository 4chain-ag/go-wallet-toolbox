package validate_test

import (
	"testing"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/fixtures"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/validate"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk/primitives"
	"github.com/stretchr/testify/require"
)

func TestForDefaultValidListCertificatesArgs(t *testing.T) {
	// given:
	args := fixtures.DefaultValidListCertificatesArgs()

	// when:
	err := validate.ListCertificatesArgs(args)

	// then:
	require.NoError(t, err)
}

func TestWrongListCertificatesArgs(t *testing.T) {
	tests := map[string]struct {
		modifier func(args *wdk.ListCertificatesArgs) *wdk.ListCertificatesArgs
	}{
		"Invalid Certifier in Certifiers list": {
			modifier: func(args *wdk.ListCertificatesArgs) *wdk.ListCertificatesArgs {
				args.Certifiers = []primitives.PubKeyHex{"invalid!"}
				return args
			},
		},
		"Certifier with odd length hex": {
			modifier: func(args *wdk.ListCertificatesArgs) *wdk.ListCertificatesArgs {
				args.Certifiers = []primitives.PubKeyHex{"abc"}
				return args
			},
		},
		"Invalid Type in Types list (non-base64)": {
			modifier: func(args *wdk.ListCertificatesArgs) *wdk.ListCertificatesArgs {
				args.Types = []primitives.Base64String{"not@base64!"}
				return args
			},
		},
		"Limit above maximum (10001)": {
			modifier: func(args *wdk.ListCertificatesArgs) *wdk.ListCertificatesArgs {
				args.Limit = 10001
				return args
			},
		},
		"Partial with invalid Certifier hex": {
			modifier: func(args *wdk.ListCertificatesArgs) *wdk.ListCertificatesArgs {
				invalid := primitives.PubKeyHex("zzzz")
				args.Partial = &wdk.ListCertificatesArgsPartial{Certifier: &invalid}
				return args
			},
		},
		"Partial with invalid Type encoding": {
			modifier: func(args *wdk.ListCertificatesArgs) *wdk.ListCertificatesArgs {
				invalid := primitives.Base64String("not-base64")
				args.Partial = &wdk.ListCertificatesArgsPartial{Type: &invalid}
				return args
			},
		},
		"Partial with invalid SerialNumber format": {
			modifier: func(args *wdk.ListCertificatesArgs) *wdk.ListCertificatesArgs {
				invalid := primitives.Base64String("invalid!")
				args.Partial = &wdk.ListCertificatesArgsPartial{SerialNumber: &invalid}
				return args
			},
		},
		"Partial with malformed RevocationOutpoint": {
			modifier: func(args *wdk.ListCertificatesArgs) *wdk.ListCertificatesArgs {
				invalid := primitives.OutpointString("missing.index")
				args.Partial = &wdk.ListCertificatesArgsPartial{RevocationOutpoint: &invalid}
				return args
			},
		},
		"Partial with invalid Signature length": {
			modifier: func(args *wdk.ListCertificatesArgs) *wdk.ListCertificatesArgs {
				invalid := primitives.HexString("abc") // Odd length
				args.Partial = &wdk.ListCertificatesArgsPartial{Signature: &invalid}
				return args
			},
		},
		"Partial with non-hex Signature": {
			modifier: func(args *wdk.ListCertificatesArgs) *wdk.ListCertificatesArgs {
				invalid := primitives.HexString("zzzz")
				args.Partial = &wdk.ListCertificatesArgsPartial{Signature: &invalid}
				return args
			},
		},
		"Partial with invalid Subject format": {
			modifier: func(args *wdk.ListCertificatesArgs) *wdk.ListCertificatesArgs {
				invalid := primitives.PubKeyHex("ghij")
				args.Partial = &wdk.ListCertificatesArgsPartial{Subject: &invalid}
				return args
			},
		},
		"Partial with numeric Outpoint index": {
			modifier: func(args *wdk.ListCertificatesArgs) *wdk.ListCertificatesArgs {
				invalid := primitives.OutpointString("deadbeef.12x") // Non-numeric index
				args.Partial = &wdk.ListCertificatesArgsPartial{RevocationOutpoint: &invalid}
				return args
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			defaultArgs := fixtures.DefaultValidListCertificatesArgs()
			modifiedArgs := test.modifier(defaultArgs)

			// when:
			err := validate.ListCertificatesArgs(modifiedArgs)

			// then:
			require.Error(t, err)
		})
	}
}
