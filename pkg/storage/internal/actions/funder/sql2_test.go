package funder_test

import (
	"testing"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/actions/funder/testabilities"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/testabilities/testusers"
)

func TestFunderSQLFund2SuccessResult(t *testing.T) {

	t.Run("target satoshis and fee are equal to the only one utxo satoshis", func(t *testing.T) {
		// given:
		given, then, cleanup := testabilities.New(t)
		defer cleanup()

		// and:
		funder := given.NewFunderService()

		// and:
		given.UTXO().OwnedBy(testusers.Alice).WithSatoshis(101).P2PKH().Stored()

		// when:
		result, err := funder.Fund2(ctx, 100, smallTransactionSize, desiredUTXONumberToPreferSingleChange, testDesiredUTXOValue, testusers.Alice.ID)

		// then:
		then.Result(result).WithoutError(err).
			HasAllocatedUTXOs().RowIndexes(0).
			HasNoChange()

	})

	t.Run("user has a lot of small utxo but they will cover the target sats and fee", func(t *testing.T) {
		// given:
		given, then, cleanup := testabilities.New(t)
		defer cleanup()

		// and:
		funder := given.NewFunderService()

		// and:
		// Funder is collecting utxos by 1000 rows, so we need to have more than 1000 utxos to test this case.
		for range 1600 {
			given.UTXO().OwnedBy(testusers.Alice).WithSatoshis(1).P2PKH().Stored()
		}

		// when:
		result, err := funder.Fund2(ctx, 1363, smallTransactionSize, desiredUTXONumberToPreferSingleChange, testDesiredUTXOValue, testusers.Alice.ID)

		// then:
		then.Result(result).WithoutError(err).
			HasAllocatedUTXOs().ForTotalAmount(1600).
			HasNoChange()

	})
}

func TestFunderSQLFund2Errors(t *testing.T) {
	t.Run("return error when user has no utxo", func(t *testing.T) {
		// given:
		given, then, cleanup := testabilities.New(t)
		defer cleanup()

		// and:
		funder := given.NewFunderService()

		// when:
		result, err := funder.Fund2(ctx, 100, smallTransactionSize, desiredUTXONumberToPreferSingleChange, testDesiredUTXOValue, testusers.Alice.ID)

		// then:
		then.Result(result).WithError(err)
	})

	t.Run("return error when user has not enough utxo", func(t *testing.T) {
		// given:
		given, then, cleanup := testabilities.New(t)
		defer cleanup()

		// and:
		funder := given.NewFunderService()

		given.UTXO().OwnedBy(testusers.Alice).WithSatoshis(10).P2PKH().Stored()

		// when:
		result, err := funder.Fund2(ctx, 100, smallTransactionSize, desiredUTXONumberToPreferSingleChange, testDesiredUTXOValue, testusers.Alice.ID)

		// then:
		then.Result(result).WithError(err)
	})

	t.Run("return error when user has not enough utxos to cover fee", func(t *testing.T) {
		// given:
		given, then, cleanup := testabilities.New(t)
		defer cleanup()

		// and:
		funder := given.NewFunderService()

		// and:
		targetSatoshis := int64(100)

		// Because apart from target satoshis we need to cover the fee also, therefore it's not enough to have only utxo for target satoshis.
		given.UTXO().OwnedBy(testusers.Alice).WithSatoshis(targetSatoshis).P2PKH().Stored()

		// when:
		result, err := funder.Fund2(ctx, targetSatoshis, smallTransactionSize, desiredUTXONumberToPreferSingleChange, testDesiredUTXOValue, testusers.Alice.ID)

		// then:
		then.Result(result).WithError(err)
	})

	t.Run("return error when user has not enough utxos to cover fee for bigger tx", func(t *testing.T) {
		// given:
		given, then, cleanup := testabilities.New(t)
		defer cleanup()

		// and:
		funder := given.NewFunderService()

		// and:
		targetSatoshis := int64(100)

		// Because the transaction size makes the fee = 2, one satoshi above the target satoshis is not enough.
		given.UTXO().OwnedBy(testusers.Alice).WithSatoshis(targetSatoshis + 1).P2PKH().Stored()

		// when:
		result, err := funder.Fund2(ctx, targetSatoshis, 1500, desiredUTXONumberToPreferSingleChange, testDesiredUTXOValue, testusers.Alice.ID)

		// then:
		then.Result(result).WithError(err)
	})
}
