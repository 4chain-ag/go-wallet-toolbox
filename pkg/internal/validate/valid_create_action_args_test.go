package validate_test

import (
	"bytes"
	"testing"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/fixtures"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/utils"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/validate"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
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
		args wdk.ValidCreateActionArgs
	}{
		"IsSendWith is set even though there is no 'send-with' txs in options": {
			args: func() wdk.ValidCreateActionArgs {
				args := fixtures.DefaultValidCreateActionArgs()
				args.IsSendWith = true
				args.Options.SendWith = []wdk.TXIDHexString{}
				return args
			}(),
		},
		"IsRemixChange is set even though there are some inputs or outputs": {
			args: func() wdk.ValidCreateActionArgs {
				args := fixtures.DefaultValidCreateActionArgs()
				args.IsRemixChange = true
				return args
			}(),
		},
		"IsNewTx is set even though there are no inputs or outputs": {
			args: func() wdk.ValidCreateActionArgs {
				args := fixtures.DefaultValidCreateActionArgs()
				args.IsNewTx = true
				args.Inputs = []wdk.ValidCreateActionInput{}
				args.Outputs = []wdk.ValidCreateActionOutput{}
				return args
			}(),
		},
		"IsSignAction is set even though there are no nil unlocking scripts": {
			args: func() wdk.ValidCreateActionArgs {
				args := fixtures.DefaultValidCreateActionArgs()
				args.IsSignAction = true
				return args
			}(),
		},
		"IsDelayed is set even though options.AcceptDelayedBroadcast is false": {
			args: func() wdk.ValidCreateActionArgs {
				args := fixtures.DefaultValidCreateActionArgs()
				args.IsDelayed = true
				args.Options.AcceptDelayedBroadcast = false
				return args
			}(),
		},
		"IsNoSend is set even though options.NoSend is false": {
			args: func() wdk.ValidCreateActionArgs {
				args := fixtures.DefaultValidCreateActionArgs()
				args.IsNoSend = true
				args.Options.NoSend = false
				return args
			}(),
		},
		"Description too short": {
			args: func() wdk.ValidCreateActionArgs {
				args := fixtures.DefaultValidCreateActionArgs()
				args.Description = "sh"
				return args
			}(),
		},
		"Description too long": {
			args: func() wdk.ValidCreateActionArgs {
				args := fixtures.DefaultValidCreateActionArgs()
				args.Description = wdk.String5to2000Bytes(bytes.Repeat([]byte{'a'}, 2001))
				return args
			}(),
		},
		"Label empty": {
			args: func() wdk.ValidCreateActionArgs {
				args := fixtures.DefaultValidCreateActionArgs()
				args.Labels = []wdk.IdentifierStringUnder300{""}
				return args
			}(),
		},
		"Label too long": {
			args: func() wdk.ValidCreateActionArgs {
				args := fixtures.DefaultValidCreateActionArgs()
				args.Labels = []wdk.IdentifierStringUnder300{wdk.IdentifierStringUnder300(bytes.Repeat([]byte{'a'}, 301))}
				return args
			}(),
		},
		"Output's locking script not in hex format": {
			args: func() wdk.ValidCreateActionArgs {
				args := fixtures.DefaultValidCreateActionArgs()
				args.Outputs[0].LockingScript = "not-hex"
				return args
			}(),
		},
		"Output's Satoshis value too high": {
			args: func() wdk.ValidCreateActionArgs {
				args := fixtures.DefaultValidCreateActionArgs()
				args.Outputs[0].Satoshis = 2100000000000001
				return args
			}(),
		},
		"Output's description too short": {
			args: func() wdk.ValidCreateActionArgs {
				args := fixtures.DefaultValidCreateActionArgs()
				args.Outputs[0].OutputDescription = "sh"
				return args
			}(),
		},
		"Output's description too long": {
			args: func() wdk.ValidCreateActionArgs {
				args := fixtures.DefaultValidCreateActionArgs()
				args.Outputs[0].OutputDescription = wdk.String5to2000Bytes(bytes.Repeat([]byte{'a'}, 2001))
				return args
			}(),
		},
		"Output's basket too long": {
			args: func() wdk.ValidCreateActionArgs {
				args := fixtures.DefaultValidCreateActionArgs()
				args.Outputs[0].Basket = utils.Ptr(wdk.IdentifierStringUnder300(bytes.Repeat([]byte{'a'}, 301)))
				return args
			}(),
		},
		"Output's basket empty": {
			args: func() wdk.ValidCreateActionArgs {
				args := fixtures.DefaultValidCreateActionArgs()
				args.Outputs[0].Basket = utils.Ptr(wdk.IdentifierStringUnder300(""))
				return args
			}(),
		},
		"Output's tag too long": {
			args: func() wdk.ValidCreateActionArgs {
				args := fixtures.DefaultValidCreateActionArgs()
				args.Outputs[0].Tags = []wdk.IdentifierStringUnder300{wdk.IdentifierStringUnder300(bytes.Repeat([]byte{'a'}, 301))}
				return args
			}(),
		},
		"Output's tag empty": {
			args: func() wdk.ValidCreateActionArgs {
				args := fixtures.DefaultValidCreateActionArgs()
				args.Outputs[0].Tags = []wdk.IdentifierStringUnder300{wdk.IdentifierStringUnder300("")}
				return args
			}(),
		},
		"Input's unlockingScript & unlockingScriptLength not provided": {
			args: func() wdk.ValidCreateActionArgs {
				args := fixtures.DefaultValidCreateActionArgs()
				args.Inputs = []wdk.ValidCreateActionInput{{}}
				return args
			}(),
		},
		"Input's unlockingScript not in hex format": {
			args: func() wdk.ValidCreateActionArgs {
				args := fixtures.DefaultValidCreateActionArgs()
				args.Inputs = []wdk.ValidCreateActionInput{{
					UnlockingScript: utils.Ptr(wdk.HexString("not-hex")),
				}}
				return args
			}(),
		},
		"Input's unlockingScript length doesn't match unlockingScriptLength": {
			args: func() wdk.ValidCreateActionArgs {
				args := fixtures.DefaultValidCreateActionArgs()
				args.Inputs = []wdk.ValidCreateActionInput{{
					UnlockingScript:       utils.Ptr(wdk.HexString("00")),
					UnlockingScriptLength: utils.Ptr(wdk.PositiveInteger(2)),
				}}
				return args
			}(),
		},
		"Input's description too short": {
			args: func() wdk.ValidCreateActionArgs {
				args := fixtures.DefaultValidCreateActionArgs()
				args.Inputs = []wdk.ValidCreateActionInput{{
					UnlockingScript:  utils.Ptr(wdk.HexString("00")),
					InputDescription: "sh",
				}}
				return args
			}(),
		},
		"Input's description too long": {
			args: func() wdk.ValidCreateActionArgs {
				args := fixtures.DefaultValidCreateActionArgs()
				args.Inputs = []wdk.ValidCreateActionInput{{
					UnlockingScript:  utils.Ptr(wdk.HexString("00")),
					InputDescription: wdk.String5to2000Bytes(bytes.Repeat([]byte{'a'}, 2001)),
				}}
				return args
			}(),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			err := validate.ValidCreateActionArgs(&test.args)

			// then:
			require.Error(t, err)
		})
	}
}
