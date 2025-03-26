package storage_test

import (
	"testing"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/defs"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/testabilities"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRPCCommunication(t *testing.T) {
	given := testabilities.Given(t)

	// given:
	mockStorage := given.MockProvider()

	// and server:
	cleanupSrv := given.StartedRPCServerFor(mockStorage)
	defer cleanupSrv()

	// and client:
	client, cleanupCli := given.RPCClient()
	defer cleanupCli()

	t.Run("Migrate", func(t *testing.T) {
		// given:
		mockStorage.EXPECT().
			Migrate(testabilities.StorageName, testabilities.StorageIdentityKey).
			Return("current-migration-version", nil)

		// when:
		migrationVersion, err := client.Migrate(testabilities.StorageName, testabilities.StorageIdentityKey)

		// then:
		require.NoError(t, err)
		assert.Equal(t, "current-migration-version", migrationVersion)
	})

	t.Run("MakeAvailable", func(t *testing.T) {
		// given:
		mockStorage.EXPECT().
			MakeAvailable().
			Return(&wdk.TableSettings{
				StorageName:        testabilities.StorageName,
				StorageIdentityKey: testabilities.StorageIdentityKey,
				Chain:              defs.NetworkTestnet,
				MaxOutputScript:    1024,
			}, nil)

		// when:
		tableSettings, err := client.MakeAvailable()

		// then:
		require.NoError(t, err)
		assert.Equal(t, testabilities.StorageName, tableSettings.StorageName)
		assert.Equal(t, testabilities.StorageIdentityKey, tableSettings.StorageIdentityKey)
		assert.Equal(t, defs.NetworkTestnet, tableSettings.Chain)
		assert.Equal(t, 1024, tableSettings.MaxOutputScript)
	})

	t.Run("FindOrInsertUser", func(t *testing.T) {
		// given:
		userIdentityKey := "03f17660f611ce531402a2ce1e070380b6fde57aca211d707bfab27bce42d86beb"

		// and:
		mockStorage.EXPECT().
			FindOrInsertUser(testabilities.StorageIdentityKey).
			Return(&wdk.FindOrInsertUserResponse{
				User: wdk.TableUser{
					IdentityKey: userIdentityKey,
				},
				IsNew: false,
			}, nil)

		// when:
		tableUser, err := client.FindOrInsertUser(testabilities.StorageIdentityKey)

		// then:
		require.NoError(t, err)
		require.NotNil(t, tableUser)
		assert.Equal(t, false, tableUser.IsNew)
		assert.Equal(t, userIdentityKey, tableUser.User.IdentityKey)
	})

	t.Run("CreateAction", func(t *testing.T) {
		t.Skip("Not implemented yet")
	})
}
