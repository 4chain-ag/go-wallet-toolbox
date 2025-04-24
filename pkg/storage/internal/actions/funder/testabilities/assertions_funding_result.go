package testabilities

import (
	"testing"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/satoshi"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/actions"
	"github.com/go-softwarelab/common/pkg/slices"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type FunderAssertion interface {
	Result(result *actions.FundingResult) FundingResultAssertion
}

type FundingResultAssertion interface {
	WithError(err error)
	WithoutError(err error) SuccessFundingResultAssertion
}

type AllocatedUTXOsAssertion interface {
	RowIndexes(indexes ...int) SuccessFundingResultAssertion
	ForTotalAmount(satoshis uint64) SuccessFundingResultAssertion
}

type SuccessFundingResultAssertion interface {
	DoesNotAllocateUTXOs() SuccessFundingResultAssertion
	HasAllocatedUTXOs() AllocatedUTXOsAssertion
	HasNoChange() SuccessFundingResultAssertion
	HasFee(fee int) SuccessFundingResultAssertion
	HasChangeCount(i int) ChangeAssertion
}

type ChangeAssertion interface {
	ForAmount(satoshis int) SuccessFundingResultAssertion
}

type funderAssertion struct {
	testing.TB
	result  *actions.FundingResult
	fixture *funderFixture
}

func newFunderAssertion(t testing.TB, fixture *funderFixture) FunderAssertion {
	return &funderAssertion{
		TB:      t,
		fixture: fixture,
	}
}

func (a *funderAssertion) Result(result *actions.FundingResult) FundingResultAssertion {
	a.Helper()
	a.result = result
	return a
}

func (a *funderAssertion) WithError(err error) {
	a.Helper()
	assert.Nil(a, a.result, "Expected error result")
	require.Error(a, err, "Expected error result")
}

func (a *funderAssertion) WithoutError(err error) SuccessFundingResultAssertion {
	a.Helper()
	assert.NoError(a, err, "Expected success result")
	require.NotNil(a, a.result, "Expected success result")
	return a
}

func (a *funderAssertion) DoesNotAllocateUTXOs() SuccessFundingResultAssertion {
	a.Helper()
	assert.Len(a, a.result.AllocatedUTXOs, 0, "Expected no allocated UTXOs")
	return a
}

func (a *funderAssertion) HasAllocatedUTXOs() AllocatedUTXOsAssertion {
	a.Helper()
	assert.NotEmptyf(a, a.result.AllocatedUTXOs, "Expected allocated UTXOs")
	return a
}

func (a *funderAssertion) ForTotalAmount(satoshis uint64) SuccessFundingResultAssertion {
	a.Helper()
	total := satoshi.Zero()
	for _, utxo := range a.result.AllocatedUTXOs {
		total += utxo.Satoshis
	}
	assert.EqualValuesf(a, satoshis, total, "Expected allocated UTXO to be for total %d but was %d", satoshis, total)
	return a
}

func (a *funderAssertion) RowIndexes(indexes ...int) SuccessFundingResultAssertion {
	a.Helper()
	expected := slices.Map(indexes, func(index int) *actions.UTXO {
		return &actions.UTXO{
			OutputID: a.fixture.createdUTXOs[index].OutputID,
			Satoshis: satoshi.MustFrom(a.fixture.createdUTXOs[index].Satoshis),
		}
	})

	assert.ElementsMatchf(a, expected, a.result.AllocatedUTXOs, "The allocated elements to match the expected ones")

	return a
}

func (a *funderAssertion) HasFee(fee int) SuccessFundingResultAssertion {
	a.Helper()
	assert.EqualValuesf(a, fee, a.result.Fee, "Expected fee to be %d but was %d", fee, a.result.Fee)
	return a
}

func (a *funderAssertion) HasNoChange() SuccessFundingResultAssertion {
	a.Helper()

	assert.Zerof(a, a.result.ChangeCount, "Unexpected number of changes")
	assert.Zerof(a, a.result.ChangeAmount, "Unexpected amount for changes")
	return a
}

func (a *funderAssertion) HasChangeCount(count int) ChangeAssertion {
	a.Helper()
	assert.EqualValuesf(a, count, a.result.ChangeCount, "Expected change count to be %d but was %d", count, a.result.ChangeCount)
	return a
}

func (a *funderAssertion) ForAmount(satoshis int) SuccessFundingResultAssertion {
	a.Helper()
	assert.EqualValuesf(a, satoshis, a.result.ChangeAmount, "Expected change amount to be %d but was %d", satoshis, a.result.ChangeAmount)
	return a
}
