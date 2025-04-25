package methodtests

import (
	"context"
	"testing"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/fixtures"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/testabilities"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/testabilities/testusers"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk/primitives"
	"github.com/go-softwarelab/common/pkg/to"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInternalizeActionNilAuth(t *testing.T) {
	given := testabilities.Given(t)

	// given:
	activeStorage := given.Provider().GORM()

	// when:
	_, err := activeStorage.InternalizeAction(context.Background(), wdk.AuthID{UserID: nil}, fixtures.DefaultInternalizeActionArgs(t, wdk.WalletPaymentProtocol))

	// then:
	require.Error(t, err)
}

func TestInternalizeActionWalletPaymentHappyPath(t *testing.T) {
	given := testabilities.Given(t)

	// given:
	activeStorage := given.Provider().GORM()

	// and:
	args := fixtures.DefaultInternalizeActionArgs(t, wdk.WalletPaymentProtocol)

	// when:
	result, err := activeStorage.InternalizeAction(
		context.Background(),
		wdk.AuthID{UserID: to.Ptr(testusers.Alice.ID)},
		args,
	)

	// then:
	require.NoError(t, err)

	assert.Equal(t, true, result.Accepted)
	assert.Equal(t, false, result.IsMerge)
	assert.Equal(t, primitives.SatoshiValue(999), result.Satoshis)
	assert.Equal(t, "a24745add717b4222d1869b3a71ad5228a3468c12f3b2bd40ce5ec84e20bf97c", result.TxID)
}

func TestInternalizeActionBasketInsertionHappyPath(t *testing.T) {
	given := testabilities.Given(t)

	// given:
	activeStorage := given.Provider().GORM()

	// and:
	args := fixtures.DefaultInternalizeActionArgs(t, wdk.BasketInsertionProtocol)

	// when:
	result, err := activeStorage.InternalizeAction(
		context.Background(),
		wdk.AuthID{UserID: to.Ptr(testusers.Alice.ID)},
		args,
	)

	// then:
	require.NoError(t, err)

	assert.Equal(t, true, result.Accepted)
	assert.Equal(t, false, result.IsMerge)
	assert.Equal(t, primitives.SatoshiValue(0), result.Satoshis)
	assert.Equal(t, "a24745add717b4222d1869b3a71ad5228a3468c12f3b2bd40ce5ec84e20bf97c", result.TxID)
}

func TestInternalizeActionErrorCases(t *testing.T) {
	tests := map[string]struct {
		modifier func(args wdk.InternalizeActionArgs) wdk.InternalizeActionArgs
	}{
		"Wrong beef": {
			modifier: func(args wdk.InternalizeActionArgs) wdk.InternalizeActionArgs {
				args.Tx = []byte{0, 1, 2, 3}
				return args
			},
		},
		"Output index out of range of provided tx": {
			modifier: func(args wdk.InternalizeActionArgs) wdk.InternalizeActionArgs {
				args.Outputs[0].OutputIndex = 999
				return args
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			given := testabilities.Given(t)

			// given:
			activeStorage := given.Provider().GORM()

			// and:
			args := test.modifier(fixtures.DefaultInternalizeActionArgs(t, wdk.WalletPaymentProtocol))

			// when:
			_, err := activeStorage.InternalizeAction(
				context.Background(),
				wdk.AuthID{UserID: to.Ptr(testusers.Alice.ID)},
				args,
			)

			// then:
			require.Error(t, err)
		})
	}
}

func TestInternalizeActionForStoredTransaction(t *testing.T) {
	given := testabilities.Given(t)

	// given:
	activeStorage := given.Provider().GORM()

	// and:
	ownedTxSpec, _ := given.Faucet(activeStorage, testusers.Alice).TopUp(100_000)
	ownedAtomicBeef, _ := ownedTxSpec.TX().AtomicBEEF(false)

	// and:
	args := fixtures.DefaultInternalizeActionArgs(t, wdk.WalletPaymentProtocol)
	args.Tx = ownedAtomicBeef

	// then:
	require.Panics(t, func() {
		// when:
		_, _ = activeStorage.InternalizeAction(
			context.Background(),
			wdk.AuthID{UserID: to.Ptr(testusers.Alice.ID)},
			args,
		)
	})

}
