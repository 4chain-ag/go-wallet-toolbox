package funder_test

import (
	"context"
	"testing"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/actions/funder/testabilities"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/testabilities/testusers"
	"github.com/go-softwarelab/common/pkg/must"
)

func TestFunderSQLFund(t *testing.T) {
	const smallTransactionSize = 44
	const transactionSizeForHigherFee = 1001
	var ctx = context.Background()

	testCasesErrors := map[string]struct {
		possessedUTXOs int64
		targetSatoshis int64
		txSize         uint64
	}{
		"return error when user has no utxo": {
			possessedUTXOs: 0,
			targetSatoshis: 100,
			txSize:         smallTransactionSize,
		},
		"return error when user fund the transaction by himself but has not enough utxo to cover the fee": {
			possessedUTXOs: 0,
			targetSatoshis: 0,
			txSize:         smallTransactionSize,
		},
		"return error when user has not enough utxo to cover the transaction": {
			possessedUTXOs: 50,
			targetSatoshis: 100,
			txSize:         smallTransactionSize,
		},
		"return error when user has not enough utxos to cover fee": {
			possessedUTXOs: 100,
			targetSatoshis: 100,
			txSize:         smallTransactionSize,
		},
		"return error when user has not enough utxos to cover fee for bigger tx": {
			// Because the transaction size makes the fee = 2, one satoshi above the target satoshis is not enough.
			possessedUTXOs: 101,
			targetSatoshis: 100,
			txSize:         transactionSizeForHigherFee,
		},
	}
	for name, test := range testCasesErrors {
		t.Run(name, func(t *testing.T) {
			// given:
			given, then, cleanup := testabilities.New(t)
			defer cleanup()

			// and:
			funder := given.NewFunderService()

			// and:
			given.UTXO().OwnedBy(testusers.Alice).WithSatoshis(test.possessedUTXOs).P2PKH().Stored()

			// and:
			basket := given.BasketFor(testusers.Alice).ThatPrefersSingleChange()

			// when:
			result, err := funder.Fund(ctx, test.targetSatoshis, test.txSize, basket, testusers.Alice.ID)

			// then:
			then.Result(result).WithError(err)
		})
	}

	// CreateAction can receive args with inputs that aren't tracked by this wallet
	// those are the test cases for handling such transactions with inputs.
	testCasesForFundingWithoutAllocatingUTXO := map[string]struct {
		possessedUTXOs int64
		targetSatoshis int64
		txSize         uint64
		expectations   func(testabilities.SuccessFundingResultAssertion)
	}{
		"user has funded exactly the transaction and fee by himself": {
			targetSatoshis: -1,
			txSize:         smallTransactionSize,

			expectations: func(thenResult testabilities.SuccessFundingResultAssertion) {
				thenResult.DoesNotAllocateUTXOs().
					HasNoChange().
					HasFee(1)
			},
		},
		"user has funded exactly the transaction and fee for bigger size of tx by himself": {
			targetSatoshis: -2,
			txSize:         transactionSizeForHigherFee,

			expectations: func(thenResult testabilities.SuccessFundingResultAssertion) {
				thenResult.DoesNotAllocateUTXOs().
					HasNoChange().
					HasFee(2)
			},
		},
		"user has funded by himself more then the transaction and fee": {
			targetSatoshis: -1001,
			txSize:         smallTransactionSize,

			expectations: func(thenResult testabilities.SuccessFundingResultAssertion) {
				thenResult.DoesNotAllocateUTXOs().
					HasChangeCount(1).ForAmount(1000).
					HasFee(1)
			},
		},
		"user has funded by himself the transaction but not the fee": {
			possessedUTXOs: 1,

			targetSatoshis: 0,
			txSize:         smallTransactionSize,

			expectations: func(thenResult testabilities.SuccessFundingResultAssertion) {
				thenResult.HasAllocatedUTXOs().ForTotalAmount(1).
					HasNoChange().
					HasFee(1)
			},
		},
		"user has funded by himself the transaction and part of the fee": {
			possessedUTXOs: 1,

			targetSatoshis: -1,
			txSize:         transactionSizeForHigherFee,

			expectations: func(thenResult testabilities.SuccessFundingResultAssertion) {
				thenResult.HasAllocatedUTXOs().ForTotalAmount(1).
					HasNoChange().
					HasFee(2)
			},
		},
	}
	for name, test := range testCasesForFundingWithoutAllocatingUTXO {
		t.Run(name, func(t *testing.T) {
			// given:
			given, then, cleanup := testabilities.New(t)
			defer cleanup()

			// and:
			funder := given.NewFunderService()

			// and:
			given.UTXO().OwnedBy(testusers.Alice).WithSatoshis(test.possessedUTXOs).P2PKH().Stored()

			// and:
			basket := given.BasketFor(testusers.Alice).ThatPrefersSingleChange()

			// when:
			result, err := funder.Fund(ctx, test.targetSatoshis, test.txSize, basket, testusers.Alice.ID)

			// then:
			test.expectations(then.Result(result).WithoutError(err))
		})
	}

	testCasesFundWholeTransaction := map[string]struct {
		havingUTXOsInDB func(testabilities.FunderFixture)
		targetSatoshis  int64
		txSize          uint64
		expectations    func(testabilities.SuccessFundingResultAssertion)
	}{
		"target satoshis and fee are equal to the only one utxo satoshis": {
			havingUTXOsInDB: func(given testabilities.FunderFixture) {
				given.UTXO().OwnedBy(testusers.Alice).WithSatoshis(101).P2PKH().Stored()
			},

			targetSatoshis: 100,
			txSize:         smallTransactionSize,

			expectations: func(thenResult testabilities.SuccessFundingResultAssertion) {
				thenResult.HasAllocatedUTXOs().RowIndexes(0).
					HasNoChange().
					HasFee(1)
			},
		},
		"adding utxo can increase the fee": {
			havingUTXOsInDB: func(given testabilities.FunderFixture) {
				given.UTXO().OwnedBy(testusers.Alice).WithSatoshis(102).P2PKH().Stored()
			},

			targetSatoshis: 100,
			txSize:         999,

			expectations: func(thenResult testabilities.SuccessFundingResultAssertion) {
				thenResult.HasAllocatedUTXOs().RowIndexes(0).
					HasFee(2).
					HasNoChange()
			},
		},
		"user has a lot of small utxo to they will cover the target sats and fee": {
			havingUTXOsInDB: func(given testabilities.FunderFixture) {
				// Funder is collecting utxos by 1000 rows, so we need to have more than 1000 utxos to test this case.
				for range 1600 {
					given.UTXO().OwnedBy(testusers.Alice).WithSatoshis(1).P2PKH().Stored()
				}
			},

			targetSatoshis: 1363,
			txSize:         smallTransactionSize,

			expectations: func(thenResult testabilities.SuccessFundingResultAssertion) {
				thenResult.HasAllocatedUTXOs().ForTotalAmount(1600).
					HasNoChange()
			},
		},
		"user has single big utxo and basket is aiming for smallest number of changes": {
			havingUTXOsInDB: func(given testabilities.FunderFixture) {
				given.UTXO().OwnedBy(testusers.Alice).WithSatoshis(10101).P2PKH().Stored()
			},

			targetSatoshis: 100,
			txSize:         smallTransactionSize,

			expectations: func(thenResult testabilities.SuccessFundingResultAssertion) {
				thenResult.HasAllocatedUTXOs().RowIndexes(0).
					HasFee(1).
					HasChangeCount(1).ForAmount(10000)
			},
		},
		"allocate biggest utxos first": {
			havingUTXOsInDB: func(given testabilities.FunderFixture) {
				given.UTXO().OwnedBy(testusers.Alice).WithSatoshis(200).P2PKH().Stored()
				given.UTXO().OwnedBy(testusers.Alice).WithSatoshis(100).P2PKH().Stored()
				given.UTXO().OwnedBy(testusers.Alice).WithSatoshis(10101).P2PKH().Stored()
				given.UTXO().OwnedBy(testusers.Alice).WithSatoshis(1).P2PKH().Stored()
				given.UTXO().OwnedBy(testusers.Alice).WithSatoshis(300).P2PKH().Stored()
			},

			targetSatoshis: 100,
			txSize:         smallTransactionSize,

			expectations: func(thenResult testabilities.SuccessFundingResultAssertion) {
				thenResult.HasAllocatedUTXOs().RowIndexes(2).
					HasFee(1).
					HasChangeCount(1).ForAmount(10000)
			},
		},
		"allocate several utxos and calculate the change from them": {
			havingUTXOsInDB: func(given testabilities.FunderFixture) {
				given.UTXO().OwnedBy(testusers.Alice).WithSatoshis(200).P2PKH().Stored()
				given.UTXO().OwnedBy(testusers.Alice).WithSatoshis(100).P2PKH().Stored()
				given.UTXO().OwnedBy(testusers.Alice).WithSatoshis(1).P2PKH().Stored()
				given.UTXO().OwnedBy(testusers.Alice).WithSatoshis(300).P2PKH().Stored()
			},

			targetSatoshis: 549,
			txSize:         smallTransactionSize,

			expectations: func(thenResult testabilities.SuccessFundingResultAssertion) {
				thenResult.HasAllocatedUTXOs().RowIndexes(0, 1, 3).
					HasFee(1).
					HasChangeCount(1).ForAmount(50)
			},
		},
	}
	for name, test := range testCasesFundWholeTransaction {
		t.Run(name, func(t *testing.T) {
			// given:
			given, then, cleanup := testabilities.New(t)
			defer cleanup()

			// and:
			funder := given.NewFunderService()

			// and:
			test.havingUTXOsInDB(given)

			// and:
			basket := given.BasketFor(testusers.Alice).ThatPrefersSingleChange()

			// when:
			result, err := funder.Fund(ctx, test.targetSatoshis, test.txSize, basket, testusers.Alice.ID)

			// then:
			test.expectations(then.Result(result).WithoutError(err))

		})
	}

	t.Run("adding change increases the fee", func(t *testing.T) {
		// given:
		given, then, cleanup := testabilities.New(t)
		defer cleanup()

		// and:
		funder := given.NewFunderService()

		// and:
		basket := given.BasketFor(testusers.Alice).ThatPrefersSingleChange()

		// when:
		result, err := funder.Fund(ctx, -102, 990, basket, testusers.Alice.ID)

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

		// and:
		basket := given.BasketFor(testusers.Alice).ThatPrefersSingleChange()

		// when:
		result, err := funder.Fund(ctx, -2, 999, basket, testusers.Alice.ID)

		// then:
		then.Result(result).WithoutError(err).
			HasFee(2).
			HasNoChange()
	})

	t.Run("produce single change when basket NumberOfDesiredUTXOs is 0", func(t *testing.T) {
		// given:
		given, then, cleanup := testabilities.New(t)
		defer cleanup()

		// and:
		funder := given.NewFunderService()

		// and:
		basket := given.BasketFor(testusers.Alice).WithNumberOfDesiredUTXOs(0)

		// when:
		result, err := funder.Fund(ctx, -5001, smallTransactionSize, basket, testusers.Alice.ID)

		// then:
		then.Result(result).WithoutError(err).
			HasChangeCount(1).ForAmount(5000)
	})

	t.Run("produce single change when basket NumberOfDesiredUTXOs is negative (value: -5)", func(t *testing.T) {
		// given:
		given, then, cleanup := testabilities.New(t)
		defer cleanup()

		// and:
		funder := given.NewFunderService()

		// and:
		basket := given.BasketFor(testusers.Alice).WithNumberOfDesiredUTXOs(-5)

		// when:
		result, err := funder.Fund(ctx, -5001, smallTransactionSize, basket, testusers.Alice.ID)

		// then:
		then.Result(result).WithoutError(err).
			HasChangeCount(1).ForAmount(5000)
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

			// and:
			given, then, cleanup := testabilities.New(t)
			defer cleanup()

			// and:
			funder := given.NewFunderService()

			// and: basket with limit of 3 outputs
			basket := given.BasketFor(testusers.Alice).WithNumberOfDesiredUTXOs(3)

			// when:
			result, err := funder.Fund(ctx, targetSatoshis, smallTransactionSize, basket, testusers.Alice.ID)

			// then:
			then.Result(result).WithoutError(err).
				HasChangeCount(test.expectedNumberOfChangeOutputs).ForAmount(test.expectedChangeValue)
		})
	}

}
