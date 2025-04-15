package methodtests

import (
	"context"
	"testing"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/defs"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/fixtures"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/testabilities/testusers"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/testabilities"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	"github.com/go-softwarelab/common/pkg/to"
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

	// when:
	result, err := activeStorage.CreateAction(
		context.Background(),
		wdk.AuthID{UserID: to.Ptr(testusers.Bob.ID)},
		args,
	)

	// TODO: Test DB state: but after we make actual getter methods, like ListActions

	// then:
	require.NoError(t, err)
	require.Equal(t, 24, len(result.DerivationPrefix))
	require.Equal(t, 16, len(result.Reference))
	require.Equal(t, args.Version, result.Version)
	require.Equal(t, args.LockTime, result.LockTime)
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
	require.Equal(t, 24, len(result.DerivationPrefix))
	require.Equal(t, 16, len(result.Reference))
	require.Equal(t, args.Version, result.Version)
	require.Equal(t, args.LockTime, result.LockTime)

	// TODO: More assertions to check the commission was applied correctly
}
