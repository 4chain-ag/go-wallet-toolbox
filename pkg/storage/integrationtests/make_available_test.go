package integrationtests_test

import (
	"testing"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/defs"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/integrationtests/testabilities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMakeAvailable(t *testing.T) {
	given := testabilities.Given(t)

	// given:
	activeStorage := given.GormProvider()

	// and:
	cleanupSrv := given.StartedRPCServerFor(activeStorage)
	defer cleanupSrv()

	// and:
	client, cleanupCli := given.RPCClient()
	defer cleanupCli()

	// when:
	tableSettings, err := client.MakeAvailable()

	// then:
	require.NoError(t, err)

	assert.Equal(t, testabilities.StorageName, tableSettings.StorageName)
	assert.Equal(t, testabilities.StorageIdentityKey, tableSettings.StorageIdentityKey)
	assert.Equal(t, defs.NetworkTestnet, tableSettings.Chain)
	assert.Equal(t, 1024, tableSettings.MaxOutputScript)
}
