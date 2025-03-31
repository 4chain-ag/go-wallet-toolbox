package funder_test

import (
	"context"
	"testing"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/actions/funder/testabilities"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/testabilities/testusers"
)

func BenchmarkFund1(b *testing.B) {
	const smallTransactionSize = 44
	const desiredUTXONumberToPreferSingleChange = 1
	const testDesiredUTXOValue = 1000

	var ctx = context.Background()

	// given:
	given, _, cleanup := testabilities.New(b)
	defer cleanup()

	// and:
	funder := given.NewFunderService()

	// and:
	// Funder is collecting utxos by 1000 rows, so we need to have more than 1000 utxos to test this case.
	for range 1600 {
		given.UTXO().OwnedBy(testusers.Alice).WithSatoshis(1).P2PKH().Stored()
	}

	b.Run("lot", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := funder.Fund(ctx, 1363, smallTransactionSize, desiredUTXONumberToPreferSingleChange, testDesiredUTXOValue, testusers.Alice.ID)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkFund2(b *testing.B) {
	const smallTransactionSize = 44
	const desiredUTXONumberToPreferSingleChange = 1
	const testDesiredUTXOValue = 1000

	var ctx = context.Background()

	// given:
	given, _, cleanup := testabilities.New(b)
	defer cleanup()

	// and:
	funder := given.NewFunderService()

	// and:
	// Funder is collecting utxos by 1000 rows, so we need to have more than 1000 utxos to test this case.
	for range 1600 {
		given.UTXO().OwnedBy(testusers.Alice).WithSatoshis(1).P2PKH().Stored()
	}

	b.Run("lot", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := funder.Fund2(ctx, 1363, smallTransactionSize, desiredUTXONumberToPreferSingleChange, testDesiredUTXOValue, testusers.Alice.ID)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
