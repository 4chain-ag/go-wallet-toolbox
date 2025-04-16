package methodtests

import (
	"context"
	"slices"
	"testing"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/defs"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/fixtures"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/testabilities/testusers"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/testabilities"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk/primitives"
	"github.com/go-softwarelab/common/pkg/to"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNilAuth(t *testing.T) {
	given := testabilities.Given(t)

	// given:
	activeStorage := given.Provider().GORM()

	// when:
	_, err := activeStorage.CreateAction(context.Background(), wdk.AuthID{UserID: nil}, fixtures.DefaultValidCreateActionArgs())

	// then:
	require.Error(t, err)
}

func TestCreateActionHappyPath(t *testing.T) {
	given := testabilities.Given(t)

	// given:
	activeStorage := given.Provider().GORM()

	// and:
	args := fixtures.DefaultValidCreateActionArgs()
	providedOutput := args.Outputs[0]

	// when:
	result, err := activeStorage.CreateAction(
		context.Background(),
		wdk.AuthID{UserID: to.Ptr(testusers.Bob.ID)},
		args,
	)

	// then:
	require.NoError(t, err)
	assert.Equal(t, 24, len(result.DerivationPrefix))
	assert.Equal(t, 16, len(result.Reference))
	assert.Equal(t, args.Version, result.Version)
	assert.Equal(t, args.LockTime, result.LockTime)
	assert.Equal(t, 1, len(result.Outputs))

	resultOutput := result.Outputs[0]
	assert.Equal(t, wdk.ProvidedByYou, resultOutput.ProvidedBy)
	assert.Empty(t, resultOutput.Purpose)
	assert.Equal(t, providedOutput.Satoshis, resultOutput.Satoshis)
	assert.Equal(t, providedOutput.Basket, resultOutput.Basket)
	assert.Equal(t, providedOutput.LockingScript, resultOutput.LockingScript)
	assert.Equal(t, providedOutput.CustomInstructions, resultOutput.CustomInstructions)
	assert.Equal(t, providedOutput.Tags, resultOutput.Tags)

	// TODO: Test DB state: but after we make actual getter methods, like ListActions
}

func TestCreateActionWithCommission(t *testing.T) {
	given := testabilities.Given(t)

	// given:
	activeStorage := given.Provider().
		WithCommission(defs.Commission{
			PubKeyHex: "03398d26f180996f8a2cb175a99620630d76257ccfef4ac7d303c8aa6f90c3190c",
			Satoshis:  10,
		}).
		GORM()

	// and:
	args := fixtures.DefaultValidCreateActionArgs()

	// when:
	result, err := activeStorage.CreateAction(context.Background(), wdk.AuthID{UserID: to.Ptr(testusers.Alice.ID)}, args)

	// then:
	require.NoError(t, err)
	assert.Equal(t, 24, len(result.DerivationPrefix))
	assert.Equal(t, 16, len(result.Reference))
	assert.Equal(t, args.Version, result.Version)
	assert.Equal(t, args.LockTime, result.LockTime)
	assert.Equal(t, 2, len(result.Outputs))

	commissionOutput, _ := findOutput(t, result.Outputs, func(p wdk.StorageCreateTransactionSdkOutput) bool {
		return p.Purpose == "storage-commission"
	})
	assert.Equal(t, primitives.SatoshiValue(10), commissionOutput.Satoshis)
	assert.Nil(t, commissionOutput.Basket)
	assert.Equal(t, wdk.ProvidedByStorage, commissionOutput.ProvidedBy)
	assert.Nil(t, commissionOutput.DerivationSuffix)
	assert.NotEmpty(t, commissionOutput.LockingScript)
	assert.NoError(t, commissionOutput.LockingScript.Validate())
	assert.Empty(t, commissionOutput.OutputDescription)
	assert.Nil(t, commissionOutput.CustomInstructions)
	assert.Empty(t, commissionOutput.Tags)
}

func TestCreateActionShuffleOutputs(t *testing.T) {
	given := testabilities.Given(t)

	// given:
	activeStorage := given.Provider().
		WithCommission(defs.Commission{
			PubKeyHex: "03398d26f180996f8a2cb175a99620630d76257ccfef4ac7d303c8aa6f90c3190c",
			Satoshis:  10,
		}).
		GORM()

	// and:
	args := fixtures.DefaultValidCreateActionArgs()

	commissionOutputVouts := map[uint32]struct{}{}
	for range 100 {
		// when:
		result, _ := activeStorage.CreateAction(
			context.Background(),
			wdk.AuthID{UserID: to.Ptr(testusers.Bob.ID)},
			args,
		)

		// then:
		found := slices.IndexFunc(result.Outputs, func(p wdk.StorageCreateTransactionSdkOutput) bool {
			return p.Purpose == "storage-commission"
		})
		commissionOutputVouts[result.Outputs[found].Vout] = struct{}{}

		if len(commissionOutputVouts) > 1 {
			t.Log(t, "Random shuffle works! Found commission outputs at different vouts")
			return
		}
	}

	t.Error(t, "Expected commission output to be shuffled, but it was not")
}

func findOutput(
	t *testing.T,
	outputs []wdk.StorageCreateTransactionSdkOutput,
	finder func(p wdk.StorageCreateTransactionSdkOutput) bool,
) (*wdk.StorageCreateTransactionSdkOutput, uint32) {
	t.Helper()
	index := slices.IndexFunc(outputs, finder)
	require.GreaterOrEqual(t, index, 0)

	return &outputs[index], uint32(index)
}
