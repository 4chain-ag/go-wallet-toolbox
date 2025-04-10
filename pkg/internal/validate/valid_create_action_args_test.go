package validate_test

import (
	"bytes"
	"testing"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/fixtures"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/validate"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk/primitives"
	"github.com/go-softwarelab/common/pkg/to"
	"github.com/stretchr/testify/require"
)

func TestForDefaultValidCreateActionArgs(t *testing.T) {
	// given:
	args := fixtures.DefaultValidCreateActionArgs()

	// when:
	err := validate.ValidCreateActionArgs(&args)

	// then:
	require.NoError(t, err)
}

func TestWrongCreateActionArgs(t *testing.T) {
	tests := map[string]struct {
		modifier func(args wdk.ValidCreateActionArgs) wdk.ValidCreateActionArgs
	}{
		"IsSendWith is set even though there is no 'send-with' txs in options": {
			modifier: func(args wdk.ValidCreateActionArgs) wdk.ValidCreateActionArgs {
				args.IsSendWith = true
				args.Options.SendWith = []primitives.TXIDHexString{}
				return args
			},
		},
		"IsRemixChange is set even though there are some inputs or outputs": {
			modifier: func(args wdk.ValidCreateActionArgs) wdk.ValidCreateActionArgs {
				args.IsRemixChange = true
				return args
			},
		},
		"IsNewTx is set even though there are no inputs or outputs": {
			modifier: func(args wdk.ValidCreateActionArgs) wdk.ValidCreateActionArgs {
				args.IsNewTx = true
				args.Inputs = []wdk.ValidCreateActionInput{}
				args.Outputs = []wdk.ValidCreateActionOutput{}
				return args
			},
		},
		"IsSignAction is set even though there are no nil unlocking scripts": {
			modifier: func(args wdk.ValidCreateActionArgs) wdk.ValidCreateActionArgs {
				args.IsSignAction = true
				return args
			},
		},
		"IsDelayed is set even though options.AcceptDelayedBroadcast is false": {
			modifier: func(args wdk.ValidCreateActionArgs) wdk.ValidCreateActionArgs {
				args.IsDelayed = true
				args.Options.AcceptDelayedBroadcast = to.Ptr[primitives.BooleanDefaultTrue](false)
				return args
			},
		},
		"IsNoSend is set even though options.NoSend is false": {
			modifier: func(args wdk.ValidCreateActionArgs) wdk.ValidCreateActionArgs {
				args.IsNoSend = true
				args.Options.NoSend = to.Ptr[primitives.BooleanDefaultFalse](false)
				return args
			},
		},
		"Description too short": {
			modifier: func(args wdk.ValidCreateActionArgs) wdk.ValidCreateActionArgs {
				args.Description = "sh"
				return args
			},
		},
		"Description too long": {
			modifier: func(args wdk.ValidCreateActionArgs) wdk.ValidCreateActionArgs {
				args.Description = primitives.String5to2000Bytes(bytes.Repeat([]byte{'a'}, 2001))
				return args
			},
		},
		"Label empty": {
			modifier: func(args wdk.ValidCreateActionArgs) wdk.ValidCreateActionArgs {
				args.Labels = []primitives.StringUnder300{""}
				return args
			},
		},
		"Label too long": {
			modifier: func(args wdk.ValidCreateActionArgs) wdk.ValidCreateActionArgs {
				args.Labels = []primitives.StringUnder300{primitives.StringUnder300(bytes.Repeat([]byte{'a'}, 301))}
				return args
			},
		},
		"Output's locking script not in hex format": {
			modifier: func(args wdk.ValidCreateActionArgs) wdk.ValidCreateActionArgs {
				args.Outputs[0].LockingScript = "not-hex"
				return args
			},
		},
		"Output's Satoshis value too high": {
			modifier: func(args wdk.ValidCreateActionArgs) wdk.ValidCreateActionArgs {
				args.Outputs[0].Satoshis = 2100000000000001
				return args
			},
		},
		"Output's description too short": {
			modifier: func(args wdk.ValidCreateActionArgs) wdk.ValidCreateActionArgs {
				args.Outputs[0].OutputDescription = "sh"
				return args
			},
		},
		"Output's description too long": {
			modifier: func(args wdk.ValidCreateActionArgs) wdk.ValidCreateActionArgs {
				args.Outputs[0].OutputDescription = primitives.String5to2000Bytes(bytes.Repeat([]byte{'a'}, 2001))
				return args
			},
		},
		"Output's basket too long": {
			modifier: func(args wdk.ValidCreateActionArgs) wdk.ValidCreateActionArgs {
				args.Outputs[0].Basket = to.Ptr(primitives.StringUnder300(bytes.Repeat([]byte{'a'}, 301)))
				return args
			},
		},
		"Output's basket empty": {
			modifier: func(args wdk.ValidCreateActionArgs) wdk.ValidCreateActionArgs {
				args.Outputs[0].Basket = to.Ptr(primitives.StringUnder300(""))
				return args
			},
		},
		"Output's tag too long": {
			modifier: func(args wdk.ValidCreateActionArgs) wdk.ValidCreateActionArgs {
				args.Outputs[0].Tags = []primitives.StringUnder300{primitives.StringUnder300(bytes.Repeat([]byte{'a'}, 301))}
				return args
			},
		},
		"Output's tag empty": {
			modifier: func(args wdk.ValidCreateActionArgs) wdk.ValidCreateActionArgs {
				args.Outputs[0].Tags = []primitives.StringUnder300{""}
				return args
			},
		},
		"Input's unlockingScript & unlockingScriptLength not provided": {
			modifier: func(args wdk.ValidCreateActionArgs) wdk.ValidCreateActionArgs {
				args.Inputs = []wdk.ValidCreateActionInput{{}}
				return args
			},
		},
		"Input's unlockingScript not in hex format": {
			modifier: func(args wdk.ValidCreateActionArgs) wdk.ValidCreateActionArgs {
				args.Inputs = []wdk.ValidCreateActionInput{{
					UnlockingScript: to.Ptr(primitives.HexString("not-hex")),
				}}
				return args
			},
		},
		"Input's unlockingScript length doesn't match unlockingScriptLength": {
			modifier: func(args wdk.ValidCreateActionArgs) wdk.ValidCreateActionArgs {
				args.Inputs = []wdk.ValidCreateActionInput{{
					UnlockingScript:       to.Ptr(primitives.HexString("00")),
					UnlockingScriptLength: to.Ptr(primitives.PositiveInteger(2)),
				}}
				return args
			},
		},
		"Input's description too short": {
			modifier: func(args wdk.ValidCreateActionArgs) wdk.ValidCreateActionArgs {
				args.Inputs = []wdk.ValidCreateActionInput{{
					UnlockingScript:  to.Ptr(primitives.HexString("00")),
					InputDescription: "sh",
				}}
				return args
			},
		},
		"Input's description too long": {
			modifier: func(args wdk.ValidCreateActionArgs) wdk.ValidCreateActionArgs {
				args.Inputs = []wdk.ValidCreateActionInput{{
					UnlockingScript:  to.Ptr(primitives.HexString("00")),
					InputDescription: primitives.String5to2000Bytes(bytes.Repeat([]byte{'a'}, 2001)),
				}}
				return args
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			defaultArgs := fixtures.DefaultValidCreateActionArgs()
			err := validate.ValidCreateActionArgs(to.Ptr(test.modifier(defaultArgs)))

			// then:
			require.Error(t, err)
		})
	}
}
