package validate_test

import (
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/fixtures"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/validate"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInconsistentValidCreateActionArgs(t *testing.T) {
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
