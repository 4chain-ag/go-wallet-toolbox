package fixtures

import (
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/utils"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
)

func DefaultValidCreateActionArgs() wdk.ValidCreateActionArgs {
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
		Labels:   []wdk.IdentifierStringUnder300{"outputbrc29"},
		Options: wdk.ValidCreateActionOptions{
			ValidProcessActionOptions: wdk.ValidProcessActionOptions{
				AcceptDelayedBroadcast: utils.Ptr[wdk.BooleanDefaultTrue](false),
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
