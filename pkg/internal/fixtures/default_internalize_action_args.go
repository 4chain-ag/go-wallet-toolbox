package fixtures

import (
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk/primitives"
	"github.com/bsv-blockchain/universal-test-vectors/pkg/testabilities"
	"github.com/stretchr/testify/require"
	"testing"
)

func DefaultInternalizeActionArgs(t *testing.T) wdk.InternalizeActionArgs {
	t.Helper()

	spec := testabilities.GivenTX().WithInput(100).WithP2PKHOutput(999)

	atomicBeef, err := spec.TX().AtomicBEEF(false)
	require.NoError(t, err)

	return wdk.InternalizeActionArgs{
		Tx: atomicBeef,
		Outputs: []*wdk.InternalizeOutput{
			{
				OutputIndex: 0,
				Protocol:    wdk.WalletPaymentProtocol,
				PaymentRemittance: &wdk.WalletPayment{
					DerivationPrefix:  DerivationPrefix,
					DerivationSuffix:  DerivationSuffix,
					SenderIdentityKey: UserIdentityKey,
				},
			},
		},
		Labels: []primitives.StringUnder300{
			"label1", "label2",
		},
		Description:    "description",
		SeekPermission: nil,
	}
}
