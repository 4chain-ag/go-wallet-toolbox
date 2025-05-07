package arc_test

import (
	"context"
	"testing"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/services/internal/testabilities"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/services/results"
	sdk "github.com/bsv-blockchain/go-sdk/transaction"
	txtestabilities "github.com/bsv-blockchain/universal-test-vectors/pkg/testabilities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostBEEFWithARCService(t *testing.T) {

	t.Run("broadcast without passing txid", func(t *testing.T) {
		// given:
		given := testabilities.Given(t)

		// setup arc server
		given.ARC().IsUpAndRunning()

		// and:
		service := given.Services().NewArcService()

		// and:
		tx := txtestabilities.GivenTX().WithInput(100).WithP2PKHOutput(99).TX()
		beef, err := sdk.NewBeefFromTransaction(tx)
		require.NoError(t, err)

		txID := tx.TxID().String()

		// when:
		res, err := service.PostBeef(context.Background(), beef, nil)

		// then:
		assert.NoError(t, err)
		assert.NotNil(t, res)

		exp := []results.PostTxID{
			{
				TxID: tx.TxID().String(),
				Data: given.ARC().TxInfoJSON(txID),
			},
		}

		require.ElementsMatch(t, exp, res.TxIDResults)
	})

	t.Run("broadcast single transaction", func(t *testing.T) {
		// given:
		given := testabilities.Given(t)

		// setup arc server
		given.ARC().IsUpAndRunning()

		// and:
		service := given.Services().NewArcService()

		// and:
		tx := txtestabilities.GivenTX().WithInput(100).WithP2PKHOutput(99).TX()
		beef, err := sdk.NewBeefFromTransaction(tx)
		require.NoError(t, err)

		txID := tx.TxID().String()
		var txids = []string{txID}

		// when:
		res, err := service.PostBeef(context.Background(), beef, txids)

		// then:
		assert.NoError(t, err)
		assert.NotNil(t, res)

		exp := []results.PostTxID{
			{
				TxID: tx.TxID().String(),
				Data: given.ARC().TxInfoJSON(txID),
			},
		}

		require.ElementsMatch(t, exp, res.TxIDResults)
	})

	t.Run("broadcast multiple txids", func(t *testing.T) {
		// given:
		given := testabilities.Given(t)

		// setup arc server
		given.ARC().IsUpAndRunning()

		// and:
		service := given.Services().NewArcService()

		// and:
		parentTx := txtestabilities.GivenTX().WithInput(100).WithP2PKHOutput(99).TX()
		parentTxID := parentTx.TxID().String()

		// and:
		childTx := txtestabilities.GivenTX().WithInputFromUTXO(parentTx, 0).WithP2PKHOutput(98).TX()
		childTxID := childTx.TxID().String()
		beef, err := sdk.NewBeefFromTransaction(childTx)
		require.NoError(t, err)

		var txids = []string{parentTxID, childTxID}

		// when:
		res, err := service.PostBeef(context.Background(), beef, txids)

		// then:
		assert.NoError(t, err)
		assert.NotNil(t, res)

		exp := []results.PostTxID{
			{
				TxID: parentTx.TxID().String(),
				Data: given.ARC().TxInfoJSON(parentTxID),
			},
			{
				TxID: childTxID,
				Data: given.ARC().TxInfoJSON(childTxID),
			},
		}

		require.ElementsMatch(t, exp, res.TxIDResults)
	})

}
