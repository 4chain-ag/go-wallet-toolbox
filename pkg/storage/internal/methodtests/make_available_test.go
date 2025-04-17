package methodtests_test

import (
	"context"
	"testing"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/defs"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/fixtures"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/testabilities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMakeAvailable(t *testing.T) {
	given := testabilities.Given(t)

	// given:
	activeStorage := given.Provider().GORMWithCleanDatabase()

	// when:
	tableSettings, err := activeStorage.MakeAvailable(context.Background())

	// then:
	require.NoError(t, err)

	assert.Equal(t, fixtures.StorageName, tableSettings.StorageName)
	assert.Equal(t, fixtures.StorageIdentityKey, tableSettings.StorageIdentityKey)
	assert.Equal(t, defs.NetworkTestnet, tableSettings.Chain)
	assert.Equal(t, 1024, tableSettings.MaxOutputScript)
}
