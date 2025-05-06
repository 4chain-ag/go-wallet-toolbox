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

	t.Run("broadcasting happy path", func(t *testing.T) {
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

		var txids = []string{tx.TxID().String()}

		// when:
		res, err := service.PostBeef(context.Background(), beef, txids)

		// then:
		assert.NoError(t, err)
		assert.NotNil(t, res)

		exp := []results.PostTxID{
			{
				TxID: tx.TxID().String(),
				Data: nil, // TODO
			},
		}

		require.ElementsMatch(t, exp, res.TxIDResults)

	})

}
