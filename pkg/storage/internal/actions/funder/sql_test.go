package funder_test

import (
	"context"
	"testing"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/actions/funder/testabilities"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/testabilities/testusers"
	"github.com/go-softwarelab/common/pkg/must"
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

	t.Run("user has funded by himself the transaction and part of the fee", func(t *testing.T) {
		// given:
		given, then, cleanup := testabilities.New(t)
		defer cleanup()

		// and:
		funder := given.NewFunderService()

		// and:
		given.UTXO().OwnedBy(testusers.Alice).WithSatoshis(1).P2PKH().Stored()

		// when:
		result, err := funder.Fund(ctx, -1, 1500, desiredUTXONumberToPreferSingleChange, testDesiredUTXOValue, testusers.Alice.ID)

		// then:
		then.Result(result).WithoutError(err).
			HasAllocatedUTXOs().RowIndexes(0).
			HasNoChange().
			HasFee(2)
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

	t.Run("adding utxo can increase the fee", func(t *testing.T) {
		// given:
		given, then, cleanup := testabilities.New(t)
		defer cleanup()

		// and:
		funder := given.NewFunderService()

		// and:
		given.UTXO().OwnedBy(testusers.Alice).WithSatoshis(102).P2PKH().Stored()

		// when:
		result, err := funder.Fund(ctx, 100, 999, desiredUTXONumberToPreferSingleChange, testDesiredUTXOValue, testusers.Alice.ID)

		// then:
		then.Result(result).WithoutError(err).
			HasAllocatedUTXOs().RowIndexes(0).
			HasFee(2).
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

	t.Run("allocate biggest utxos first", func(t *testing.T) {
		// given:
		given, then, cleanup := testabilities.New(t)
		defer cleanup()

		// and:
		funder := given.NewFunderService()

		// and:
		given.UTXO().OwnedBy(testusers.Alice).WithSatoshis(200).P2PKH().Stored()
		given.UTXO().OwnedBy(testusers.Alice).WithSatoshis(100).P2PKH().Stored()
		given.UTXO().OwnedBy(testusers.Alice).WithSatoshis(10101).P2PKH().Stored()
		given.UTXO().OwnedBy(testusers.Alice).WithSatoshis(1).P2PKH().Stored()
		given.UTXO().OwnedBy(testusers.Alice).WithSatoshis(300).P2PKH().Stored()

		// when:
		result, err := funder.Fund(ctx, 100, smallTransactionSize, desiredUTXONumberToPreferSingleChange, testDesiredUTXOValue, testusers.Alice.ID)

		// then:
		then.Result(result).WithoutError(err).
			HasAllocatedUTXOs().RowIndexes(2).
			HasFee(1).
			HasChangeCount(1).ForAmount(10000)

	})

	t.Run("allocate several utxos and calculate the change from them", func(t *testing.T) {
		// given:
		given, then, cleanup := testabilities.New(t)
		defer cleanup()

		// and:
		funder := given.NewFunderService()

		// and:
		given.UTXO().OwnedBy(testusers.Alice).WithSatoshis(200).P2PKH().Stored()
		given.UTXO().OwnedBy(testusers.Alice).WithSatoshis(100).P2PKH().Stored()
		given.UTXO().OwnedBy(testusers.Alice).WithSatoshis(1).P2PKH().Stored()
		given.UTXO().OwnedBy(testusers.Alice).WithSatoshis(300).P2PKH().Stored()

		// when:
		result, err := funder.Fund(ctx, 549, smallTransactionSize, desiredUTXONumberToPreferSingleChange, testDesiredUTXOValue, testusers.Alice.ID)

		// then:
		then.Result(result).WithoutError(err).
			HasAllocatedUTXOs().RowIndexes(0, 1, 3).
			HasFee(1).
			HasChangeCount(1).ForAmount(50)

	})

	t.Run("adding change increases the fee", func(t *testing.T) {
		// given:
		given, then, cleanup := testabilities.New(t)
		defer cleanup()

		// and:
		funder := given.NewFunderService()

		// when:
		result, err := funder.Fund(ctx, -102, 990, desiredUTXONumberToPreferSingleChange, testDesiredUTXOValue, testusers.Alice.ID)

		// then:
		then.Result(result).WithoutError(err).
			HasChangeCount(1).ForAmount(100).
			HasFee(2)
	})

	t.Run("adding change will increase the fee so that there won't be any change, so we're giving extra fee to miner", func(t *testing.T) {
		// given:
		given, then, cleanup := testabilities.New(t)
		defer cleanup()

		// and:
		funder := given.NewFunderService()

		// when:
		result, err := funder.Fund(ctx, -2, 999, desiredUTXONumberToPreferSingleChange, testDesiredUTXOValue, testusers.Alice.ID)

		// then:
		then.Result(result).WithoutError(err).
			HasFee(2).
			HasNoChange()
	})

	testCasesSplitUserProvidedInputIntoChanges := map[string]struct {
		expectedChangeValue           int
		expectedNumberOfChangeOutputs int
	}{
		"change (value: 249) below minimum desired utxo creates single output": {
			expectedChangeValue:           249,
			expectedNumberOfChangeOutputs: 1,
		},
		"change (value: 250) below minimum desired utxo creates single output": {
			expectedChangeValue:           250,
			expectedNumberOfChangeOutputs: 1,
		},
		"change equal to minimum desired utxo creates single output": {
			expectedChangeValue:           1000,
			expectedNumberOfChangeOutputs: 1,
		},
		"change (value 1001) below 125% of minimum desired utxo creates single output": {
			expectedChangeValue:           1001,
			expectedNumberOfChangeOutputs: 1,
		},
		"change (value 1249) below 125% of minimum desired utxo creates single output": {
			expectedChangeValue:           1249,
			expectedNumberOfChangeOutputs: 1,
		},
		"change equal to 125% of minimum desired utxo creates two outputs": {
			expectedChangeValue:           1250,
			expectedNumberOfChangeOutputs: 2,
		},
		"change equal to 200% of minimum desired utxo creates two outputs": {
			expectedChangeValue:           2000,
			expectedNumberOfChangeOutputs: 2,
		},
		"change above 200% but below 225% of minimum desired utxo creates two outputs": {
			expectedChangeValue:           2249,
			expectedNumberOfChangeOutputs: 2,
		},
		"change above 225% of minimum desired utxo creates three outputs": {
			expectedChangeValue:           2250,
			expectedNumberOfChangeOutputs: 3,
		},
		"change equal to (minimum desired utxo) times (number of desired utxo) creates desired utxo number of changes": {
			expectedChangeValue:           3000,
			expectedNumberOfChangeOutputs: 3,
		},
		"change above the (minimum desired utxo) times (number of desired utxo) creates desired utxo number of changes": {
			expectedChangeValue:           10000,
			expectedNumberOfChangeOutputs: 3,
		},
	}
	for name, test := range testCasesSplitUserProvidedInputIntoChanges {
		t.Run(name, func(t *testing.T) {
			// given:
			fee := 1

			// and: targetSatoshis should cover the fee and the expected change value
			// and it must be negative to simulate that user provides by himself the inputs to cover those values.
			targetSatoshis := must.ConvertToInt64(-(test.expectedChangeValue + fee))

			// and: this is the limit for number of changes we don't want to exceed (in those test cases)
			const numberOfDesiredUTXOs = 3

			// and:
			given, then, cleanup := testabilities.New(t)
			defer cleanup()

			// and:
			funder := given.NewFunderService()

			// when:
			result, err := funder.Fund(ctx, targetSatoshis, smallTransactionSize, numberOfDesiredUTXOs, testDesiredUTXOValue, testusers.Alice.ID)

			// then:
			then.Result(result).WithoutError(err).
				HasChangeCount(test.expectedNumberOfChangeOutputs).ForAmount(test.expectedChangeValue)
		})
	}
}
