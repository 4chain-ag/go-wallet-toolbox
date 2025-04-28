package fixtures

import (
	"testing"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk/primitives"
	"github.com/bsv-blockchain/universal-test-vectors/pkg/testabilities"
	"github.com/go-softwarelab/common/pkg/to"
	"github.com/stretchr/testify/require"
)

func DefaultInternalizeActionArgs(t *testing.T, protocol wdk.InternalizeProtocol) wdk.InternalizeActionArgs {
	t.Helper()

	spec := testabilities.GivenTX().WithInput(1000).WithP2PKHOutput(ExpectedValueToInternalize)

	atomicBeef, err := spec.TX().AtomicBEEF(false)
	require.NoError(t, err)

	outputSpec := &wdk.InternalizeOutput{
		OutputIndex: 0,
		Protocol:    protocol,
	}
	if protocol == wdk.WalletPaymentProtocol {
		outputSpec.PaymentRemittance = &wdk.WalletPayment{
			DerivationPrefix:  DerivationPrefix,
			DerivationSuffix:  DerivationSuffix,
			SenderIdentityKey: UserIdentityKey,
		}
	} else {
		outputSpec.InsertionRemittance = &wdk.BasketInsertion{
			Basket:             CustomBasket,
			CustomInstructions: to.Ptr("custom instructions"),
			Tags:               []primitives.StringUnder300{"tag1", "tag2"},
		}
	}

	return wdk.InternalizeActionArgs{
		Tx: atomicBeef,
		Outputs: []*wdk.InternalizeOutput{
			outputSpec,
		},
		Labels: []primitives.StringUnder300{
			"label1", "label2",
		},
		Description:    "description",
		SeekPermission: nil,
	}
}
