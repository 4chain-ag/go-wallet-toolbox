package integrationtests

import (
	"context"
	"testing"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/fixtures"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/actions/funder/errfunder"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/testabilities"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/testabilities/testusers"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/testabilities/testutils"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk/primitives"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInternalizePlusCreate(t *testing.T) {
	given := testabilities.Given(t)
	activeStorage := given.Provider().GORM()

	var internalizedTxID string

	t.Run("Internalize", func(t *testing.T) {
		// given:
		args := fixtures.DefaultInternalizeActionArgs(t, wdk.WalletPaymentProtocol)

		// when:
		result, err := activeStorage.InternalizeAction(
			context.Background(),
			testusers.Alice.AuthID(),
			args,
		)

		// then:
		require.NoError(t, err)
		require.Equal(t, true, result.Accepted)

		// update:
		internalizedTxID = result.TxID
	})

	t.Run("Create", func(t *testing.T) {
		// given:
		args := fixtures.DefaultValidCreateActionArgs()
		args.Outputs[0].Satoshis = 1

		// when:
		result, err := activeStorage.CreateAction(
			context.Background(),
			testusers.Alice.AuthID(),
			args,
		)

		// then:
		require.NoError(t, err)

		providedOutput, _ := testutils.FindOutput(t, result.Outputs, testutils.ProvidedByYouCondition)
		assert.Equal(t, primitives.SatoshiValue(1), providedOutput.Satoshis)

		changeValue := testutils.SumOutputsWithCondition(t, result.Outputs, testutils.SatoshiValue, testutils.ProvidedByStorageCondition)
		assert.Equal(t, primitives.SatoshiValue(fixtures.ExpectedValueToInternalize-1-1), changeValue)

		require.Equal(t, 1, len(result.Inputs))
		allocatedUTXO := result.Inputs[0]
		assert.Equal(t, internalizedTxID, allocatedUTXO.SourceTxid)
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

func TestInternalizeBasketInsertionThenCreat(t *testing.T) {
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
