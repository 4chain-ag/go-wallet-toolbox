package funder_test

import (
	"context"
	"testing"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/actions/funder/testabilities"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/testabilities/testusers"
)

const smallTransactionSize = 44
const desiredUTXONumberToPreferSingleChange = 1
const testDesiredUTXOValue = 1000

var ctx = context.Background()

func TestFunderSQLFundSuccessResult(t *testing.T) {
	t.Run("return error when user has no utxo", func(t *testing.T) {
		// given:
		given, then, cleanup := testabilities.New(t)
		defer cleanup()

		// and:
		funder := given.NewFunderService()

		// when:
		result, err := funder.Fund(ctx, 100, smallTransactionSize, desiredUTXONumberToPreferSingleChange, testDesiredUTXOValue, testusers.Alice.ID)

		// then:
		then.Result(result).WithError(err)
	})

	t.Run("user has funded exactly the transaction and fee by himself", func(t *testing.T) {
		// given:
		given, then, cleanup := testabilities.New(t)
		defer cleanup()

		// and:
		funder := given.NewFunderService()

		// when:
		result, err := funder.Fund(ctx, -1, smallTransactionSize, desiredUTXONumberToPreferSingleChange, testDesiredUTXOValue, testusers.Alice.ID)

		// then:
		then.Result(result).WithoutError(err).
			DoesNotAllocateUTXOs().
			HasNoChange().
			HasFee(1)

	})

	t.Run("user has funded exactly the transaction and fee for bigger size of tx by himself", func(t *testing.T) {
		// given:
		given, then, cleanup := testabilities.New(t)
		defer cleanup()

		// and:
		funder := given.NewFunderService()

		// when:
		result, err := funder.Fund(ctx, -2, 1001, desiredUTXONumberToPreferSingleChange, testDesiredUTXOValue, testusers.Alice.ID)

		// then:
		then.Result(result).WithoutError(err).
			DoesNotAllocateUTXOs().
			HasNoChange().
			HasFee(2)

	})

	t.Run("return error when user fund the transaction by himself but has not enough utxo to cover the fee", func(t *testing.T) {
		// given:
		given, then, cleanup := testabilities.New(t)
		defer cleanup()

		// and:
		funder := given.NewFunderService()

		// when:
		result, err := funder.Fund(ctx, 0, smallTransactionSize, desiredUTXONumberToPreferSingleChange, testDesiredUTXOValue, testusers.Alice.ID)

		// then:
		then.Result(result).WithError(err)
	})

	t.Run("user has funded by himself the transaction but not the fee", func(t *testing.T) {
		// given:
		given, then, cleanup := testabilities.New(t)
		defer cleanup()

		// and:
		funder := given.NewFunderService()

		// and:
		given.UTXO().OwnedBy(testusers.Alice).WithSatoshis(1).P2PKH().Stored()

		// when:
		result, err := funder.Fund(ctx, 0, smallTransactionSize, desiredUTXONumberToPreferSingleChange, testDesiredUTXOValue, testusers.Alice.ID)

		// then:
		then.Result(result).WithoutError(err).
			HasAllocatedUTXOs().RowIndexes(0).
			HasNoChange().
			HasFee(1)
	})

	t.Run("test when target sat are < 0, but after adding the fee it will be more than 0", func(t *testing.T) {
		// TODO
		t.Skip("not implemented")
	})

	t.Run("test that assigning the change will increase fee and decrease change", func(t *testing.T) {
		// TODO
		t.Skip("not implemented")
	})

	t.Run("test that if adding change will increase the fee so that there is no need for change, then return higher fee but no change", func(t *testing.T) {
		// TODO
		t.Skip("not implemented")
	})

	t.Run("user has funded by himself more then the transaction and fee", func(t *testing.T) {
		// given:
		given, then, cleanup := testabilities.New(t)
		defer cleanup()

		// and:
		funder := given.NewFunderService()

		// when:
		result, err := funder.Fund(ctx, -1001, smallTransactionSize, desiredUTXONumberToPreferSingleChange, testDesiredUTXOValue, testusers.Alice.ID)

		// then:
		then.Result(result).WithoutError(err).
			DoesNotAllocateUTXOs().
			HasChangeCount(1).ForAmount(1000).
			HasFee(1)

	})

	t.Run("return error when user has not enough utxo to cover the transaction", func(t *testing.T) {
		// given:
		given, then, cleanup := testabilities.New(t)
		defer cleanup()

		// and:
		funder := given.NewFunderService()

		given.UTXO().OwnedBy(testusers.Alice).WithSatoshis(10).P2PKH().Stored()

		// when:
		result, err := funder.Fund(ctx, 100, smallTransactionSize, desiredUTXONumberToPreferSingleChange, testDesiredUTXOValue, testusers.Alice.ID)

		// then:
		then.Result(result).WithError(err)
	})

	t.Run("target satoshis and fee are equal to the only one utxo satoshis", func(t *testing.T) {
		// given:
		given, then, cleanup := testabilities.New(t)
		defer cleanup()

		// and:
		funder := given.NewFunderService()

		// and:
		given.UTXO().OwnedBy(testusers.Alice).WithSatoshis(101).P2PKH().Stored()

		// when:
		result, err := funder.Fund(ctx, 100, smallTransactionSize, desiredUTXONumberToPreferSingleChange, testDesiredUTXOValue, testusers.Alice.ID)

		// then:
		then.Result(result).WithoutError(err).
			HasAllocatedUTXOs().RowIndexes(0).
			HasNoChange().
			HasFee(1)

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
		result, err := funder.Fund(ctx, targetSatoshis, smallTransactionSize, desiredUTXONumberToPreferSingleChange, testDesiredUTXOValue, testusers.Alice.ID)

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
		result, err := funder.Fund(ctx, targetSatoshis, 1500, desiredUTXONumberToPreferSingleChange, testDesiredUTXOValue, testusers.Alice.ID)

		// then:
		then.Result(result).WithError(err)
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
		result, err := funder.Fund(ctx, 1363, smallTransactionSize, desiredUTXONumberToPreferSingleChange, testDesiredUTXOValue, testusers.Alice.ID)

		// then:
		then.Result(result).WithoutError(err).
			HasAllocatedUTXOs().ForTotalAmount(1600).
			HasNoChange()

	})

	t.Run("user has single big utxo and aiming for smallest number of changes", func(t *testing.T) {
		// given:
		given, then, cleanup := testabilities.New(t)
		defer cleanup()

		// and:
		funder := given.NewFunderService()

		// and:
		given.UTXO().OwnedBy(testusers.Alice).WithSatoshis(10101).P2PKH().Stored()

		// when:
		result, err := funder.Fund(ctx, 100, smallTransactionSize, desiredUTXONumberToPreferSingleChange, testDesiredUTXOValue, testusers.Alice.ID)

		// then:
		then.Result(result).WithoutError(err).
			HasAllocatedUTXOs().RowIndexes(0).
			HasFee(1).
			HasChangeCount(1).ForAmount(10000)

	})
}
