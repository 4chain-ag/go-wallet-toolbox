package integrationtests

import (
	"context"
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/fixtures"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/randomizer"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/actions/funder/errfunder"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/integrationtests/tsgenerated"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/testabilities"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/testabilities/testusers"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk/primitives"
	"github.com/go-softwarelab/common/pkg/to"
	"github.com/stretchr/testify/require"
)

//go:embed tsgenerated/create_action_result.json
var createActionResultJSON string

func TestInternalizePlusCreate(t *testing.T) {
	given := testabilities.Given(t)
	activeStorage := given.Provider().
		WithRandomizer(randomizer.NewTestRandomizer()).
		GORM()

	t.Run("Internalize", func(t *testing.T) {
		// given:
		args := wdk.InternalizeActionArgs{
			Tx: tsgenerated.AtomicBeefToInternalize(t),
			Outputs: []*wdk.InternalizeOutput{
				{
					OutputIndex: 0,
					Protocol:    wdk.WalletPaymentProtocol,
					PaymentRemittance: &wdk.WalletPayment{
						DerivationPrefix:  fixtures.DerivationPrefix,
						DerivationSuffix:  fixtures.DerivationSuffix,
						SenderIdentityKey: fixtures.AnyoneIdentityKey,
					},
				},
			},
			Labels: []primitives.StringUnder300{
				"label1", "label2",
			},
			Description:    "description",
			SeekPermission: nil,
		}

		// when:
		result, err := activeStorage.InternalizeAction(
			context.Background(),
			testusers.Alice.AuthID(),
			args,
		)

		// then:
		require.NoError(t, err)

		// when:
		resultJSON, err := json.Marshal(result)

		// then:
		require.NoError(t, err)

		require.JSONEq(t, `{
		  "accepted": true,
		  "isMerge": false,
		  "txid": "db8ee15998a415f69144483cf9a755d05f2a7c44e3569c8fe750720a26f90fe7",
		  "satoshis": 99904
		}`, string(resultJSON))
	})

	t.Run("Create", func(t *testing.T) {
		// given:
		args := wdk.ValidCreateActionArgs{
			Description: "outputBRC29",
			Inputs:      []wdk.ValidCreateActionInput{},
			Outputs: []wdk.ValidCreateActionOutput{
				{
					LockingScript:      "76a9144b0d6cbef5a813d2d12dcec1de2584b250dc96a388ac",
					Satoshis:           1000,
					OutputDescription:  "outputBRC29",
					CustomInstructions: to.Ptr(`{"derivationPrefix":"Pr==","derivationSuffix":"Su==","type":"BRC29"}`),
				},
			},
			LockTime: 0,
			Version:  1,
			Labels:   []primitives.StringUnder300{"outputbrc29"},
			Options: wdk.ValidCreateActionOptions{
				ValidProcessActionOptions: wdk.ValidProcessActionOptions{
					AcceptDelayedBroadcast: to.Ptr[primitives.BooleanDefaultTrue](false),
					SendWith:               []primitives.TXIDHexString{},
				},
				SignAndProcess:   to.Ptr(primitives.BooleanDefaultTrue(true)),
				KnownTxids:       []primitives.TXIDHexString{},
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

		// when:
		result, err := activeStorage.CreateAction(
			context.Background(),
			testusers.Alice.AuthID(),
			args,
		)

		// then:
		require.NoError(t, err)

		// when:
		resultJSON, err := json.Marshal(result)

		// then:
		require.NoError(t, err)
		require.JSONEq(t, createActionResultJSON, string(resultJSON))
	})
}

func TestInternalizePlusTooHighCreate(t *testing.T) {
	given := testabilities.Given(t)
	activeStorage := given.Provider().GORM()

	t.Run("Internalize", func(t *testing.T) {
		// given:
		args := fixtures.DefaultInternalizeActionArgs(t, wdk.BasketInsertionProtocol)

		// when:
		result, err := activeStorage.InternalizeAction(
			context.Background(),
			testusers.Alice.AuthID(),
			args,
		)

		// then:
		require.NoError(t, err)
		require.Equal(t, true, result.Accepted)
	})

	t.Run("Create", func(t *testing.T) {
		// given:
		args := fixtures.DefaultValidCreateActionArgs()
		args.Outputs[0].Satoshis = 2 * fixtures.ExpectedValueToInternalize

		// when:
		_, err := activeStorage.CreateAction(
			context.Background(),
			testusers.Alice.AuthID(),
			args,
		)

		// then:
		require.ErrorIs(t, err, errfunder.NotEnoughFunds)
	})
}

func TestInternalizeBasketInsertionThenCreate(t *testing.T) {
	given := testabilities.Given(t)
	activeStorage := given.Provider().GORM()

	t.Run("Internalize", func(t *testing.T) {
		// given:
		args := fixtures.DefaultInternalizeActionArgs(t, wdk.BasketInsertionProtocol)

		// when:
		result, err := activeStorage.InternalizeAction(
			context.Background(),
			testusers.Alice.AuthID(),
			args,
		)

		// then:
		require.NoError(t, err)
		require.Equal(t, true, result.Accepted)
	})

	t.Run("Create", func(t *testing.T) {
		// given:
		args := fixtures.DefaultValidCreateActionArgs()
		args.Outputs[0].Satoshis = 1

		// when:
		_, err := activeStorage.CreateAction(
			context.Background(),
			testusers.Alice.AuthID(),
			args,
		)

		// then:
		require.ErrorIs(t, err, errfunder.NotEnoughFunds)
	})
}
