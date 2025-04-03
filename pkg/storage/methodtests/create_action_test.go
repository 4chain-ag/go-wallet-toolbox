package methodtests

import (
	"testing"

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

	// when:
	_, err := activeStorage.CreateAction(
		wdk.AuthID{UserID: to.Ptr(testusers.Alice.ID)},
		fixtures.DefaultValidCreateActionArgs(),
	)

	// then:
	require.NoError(t, err)
}
