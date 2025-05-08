package arc_test

import (
	"context"
	"net/http"
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
				Result: results.ResultStatusSuccess,
				TxID:   tx.TxID().String(),
				Data:   given.ARC().TxInfoJSON(txID),
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
				Result: results.ResultStatusSuccess,
				TxID:   tx.TxID().String(),
				Data:   given.ARC().TxInfoJSON(txID),
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
				Result: results.ResultStatusSuccess,
				TxID:   parentTx.TxID().String(),
				Data:   given.ARC().TxInfoJSON(parentTxID),
			},
			{
				Result: results.ResultStatusSuccess,
				TxID:   childTxID,
				Data:   given.ARC().TxInfoJSON(childTxID),
			},
		}

		require.ElementsMatch(t, exp, res.TxIDResults)
	})

	invalidBEEFTestCases := map[string]struct {
		BEEF func(t testing.TB) *sdk.Beef
	}{
		"return error on nil beef": {
			BEEF: func(t testing.TB) *sdk.Beef {
				return nil
			},
		},
		"return error on empty beef": {
			BEEF: func(t testing.TB) *sdk.Beef {
				return sdk.NewBeefV2()
			},
		},
		"return error on beef v2 with txID only": {
			BEEF: func(t testing.TB) *sdk.Beef {
				tx := txtestabilities.GivenTX().WithInput(100).WithP2PKHOutput(99).TX()
				beef, err := sdk.NewBeefFromTransaction(tx)
				require.NoError(t, err)

				beefTxIDOnly, err := beef.TxidOnly()
				require.NoError(t, err)
				return beefTxIDOnly
			},
		},
		"return error on beef v2 with multiple subject transactions (TEMPORARY)": {
			BEEF: func(t testing.TB) *sdk.Beef {
				tx1 := txtestabilities.GivenTX().WithInput(100).WithP2PKHOutput(99).TX()
				beef, err := sdk.NewBeefFromTransaction(tx1)
				require.NoError(t, err)

				tx2 := txtestabilities.GivenTX().WithInput(200).WithP2PKHOutput(100).TX()

				_, err = beef.MergeTransaction(tx2)
				require.NoError(t, err)

				return beef
			},
		},
	}
	for name, test := range invalidBEEFTestCases {
		t.Run(name, func(t *testing.T) {
			// given:
			given := testabilities.Given(t)

			// setup arc server
			given.ARC().IsUpAndRunning()

			// and:
			service := given.Services().NewArcService()

			// when:
			res, err := service.PostBeef(context.Background(), test.BEEF(t), nil)

			// then:
			assert.Error(t, err)
			assert.Nil(t, res)
		})
	}

	arcFailingTestCases := map[string]struct {
		setupARC func(testabilities.ArcFixture)
	}{
		"return error when arc is unreachable": {
			setupARC: func(testabilities.ArcFixture) {},
		},
		"return error when arc returns unauthorized": {
			setupARC: func(arc testabilities.ArcFixture) {
				arc.WillAlwaysReturnStatus(http.StatusUnauthorized)
			},
		},
		"return error when arc returns forbidden": {
			setupARC: func(arc testabilities.ArcFixture) {
				arc.WillAlwaysReturnStatus(http.StatusForbidden)
			},
		},
		"return error when arc returns not found": {
			setupARC: func(arc testabilities.ArcFixture) {
				arc.WillAlwaysReturnStatus(http.StatusNotFound)
			},
		},
		"return error when arc returns internal server error": {
			setupARC: func(arc testabilities.ArcFixture) {
				arc.WillAlwaysReturnStatus(http.StatusInternalServerError)
			},
		},
	}
	for name, test := range arcFailingTestCases {
		t.Run(name, func(t *testing.T) {
			// given:
			given := testabilities.Given(t)

			// setup arc server
			test.setupARC(given.ARC())

			// and:
			service := given.Services().NewArcService()

			// and:
			tx := txtestabilities.GivenTX().WithInput(100).WithP2PKHOutput(99).TX()
			beef, err := sdk.NewBeefFromTransaction(tx)
			require.NoError(t, err)

			// when:
			res, err := service.PostBeef(context.Background(), beef, nil)

			// then:
			assert.Error(t, err)
			assert.Nil(t, res)
		})
	}
}
