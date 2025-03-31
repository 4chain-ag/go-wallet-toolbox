package methodtests

import (
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/utils"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/testabilities/testusers"
	"testing"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/fixtures"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/testabilities"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	"github.com/stretchr/testify/require"
)

func TestNilAuth(t *testing.T) {
	given := testabilities.Given(t)

	// given:
	activeStorage := given.GormProvider()

	// when:
	_, err := activeStorage.CreateAction(wdk.AuthID{UserID: nil}, fixtures.DefaultValidCreateActionArgs())

	// then:
	require.Error(t, err)
}

func TestCreateActionHappyPath(t *testing.T) {
	given := testabilities.Given(t)

	// given:
	activeStorage := given.GormProvider()

	// and:
	// TODO: remove this after dzolt-4chain merge his PR
	userResp, err := activeStorage.FindOrInsertUser(testusers.Alice.PrivKey)
	require.NoError(t, err)

	// when:
	_, err = activeStorage.CreateAction(
		wdk.AuthID{UserID: utils.Ptr(userResp.User.UserID)},
		fixtures.DefaultValidCreateActionArgs(),
	)

	// then:
	require.NoError(t, err)
}
