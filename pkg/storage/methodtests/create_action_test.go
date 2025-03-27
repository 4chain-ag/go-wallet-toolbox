package methodtests

import (
	"testing"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/utils"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/testabilities"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	"github.com/stretchr/testify/require"
)

func defaultArgs() wdk.ValidCreateActionArgs {
	return wdk.ValidCreateActionArgs{
		Description: "outputBRC29",
		Inputs:      []wdk.ValidCreateActionInput{},
		Outputs: []wdk.ValidCreateActionOutput{
			{
				LockingScript:      "76a914dbc0a7c84983c5bf199b7b2d41b3acf0408ee5aa88ac",
				Satoshis:           42000,
				OutputDescription:  "outputBRC29",
				CustomInstructions: utils.Ptr(`{"derivationPrefix":"bPRI9FYwsIo=","derivationSuffix":"FdjLdpnLnJM=","type":"BRC29"}`),
			},
		},
		LockTime: 0,
		Version:  1,
		Labels:   []string{"outputbrc29"},
		Options: wdk.ValidCreateActionOptions{
			ValidProcessActionOptions: wdk.ValidProcessActionOptions{
				AcceptDelayedBroadcast: false,
				NoSend:                 false,
				ReturnTXIDOnly:         false,
				SendWith:               []wdk.TXIDHexString{},
			},
			SignAndProcess:   true,
			KnownTxids:       []wdk.TXIDHexString{},
			NoSendChange:     []wdk.OutPoint{},
			RandomizeOutputs: false,
		},
		IsSendWith:                   false,
		IsDelayed:                    false,
		IsNoSend:                     false,
		IsNewTx:                      true,
		IsRemixChange:                false,
		IsSignAction:                 false,
		IncludeAllSourceTransactions: true,
	}
}

func TestInconsistentArgs(t *testing.T) {
	given := testabilities.Given(t)

	// given:
	activeStorage := given.GormProvider()

	// and:
	userID := 1
	authID := wdk.AuthID{UserID: &userID}

	tests := map[string]struct {
		args wdk.ValidCreateActionArgs
	}{
		"IsSendWith is set even though there is no 'send-with' txs in options": {
			args: func() wdk.ValidCreateActionArgs {
				args := defaultArgs()
				args.IsSendWith = true
				args.Options.SendWith = []wdk.TXIDHexString{}
				return args
			}(),
		},
		"IsRemixChange is set even though there are some inputs or outputs": {
			args: func() wdk.ValidCreateActionArgs {
				args := defaultArgs()
				args.IsRemixChange = true
				return args
			}(),
		},
		"IsNewTx is set even though there are no inputs or outputs": {
			args: func() wdk.ValidCreateActionArgs {
				args := defaultArgs()
				args.IsNewTx = true
				args.Inputs = []wdk.ValidCreateActionInput{}
				args.Outputs = []wdk.ValidCreateActionOutput{}
				return args
			}(),
		},
		"IsSignAction is set even though there are no nil unlocking scripts": {
			args: func() wdk.ValidCreateActionArgs {
				args := defaultArgs()
				args.IsSignAction = true
				return args
			}(),
		},
		"IsDelayed is set even though options.AcceptDelayedBroadcast is false": {
			args: func() wdk.ValidCreateActionArgs {
				args := defaultArgs()
				args.IsDelayed = true
				args.Options.AcceptDelayedBroadcast = false
				return args
			}(),
		},
		"IsNoSend is set even though options.NoSend is false": {
			args: func() wdk.ValidCreateActionArgs {
				args := defaultArgs()
				args.IsNoSend = true
				args.Options.NoSend = false
				return args
			}(),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// when:
			_, err := activeStorage.CreateAction(authID, test.args)

			// then:
			require.Error(t, err)
		})
	}
}

func TestNilAuth(t *testing.T) {
	given := testabilities.Given(t)

	// given:
	activeStorage := given.GormProvider()

	// when:
	_, err := activeStorage.CreateAction(wdk.AuthID{UserID: nil}, defaultArgs())

	// then:
	require.Error(t, err)
}
