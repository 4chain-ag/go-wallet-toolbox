package storage_test

import (
	"context"
	"testing"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/defs"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/fixtures"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/testabilities"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
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
			Migrate(gomock.Any(), fixtures.StorageName, fixtures.StorageIdentityKey).
			Return("current-migration-version", nil)

		// when:
		migrationVersion, err := client.Migrate(context.Background(), fixtures.StorageName, fixtures.StorageIdentityKey)

		// then:
		require.NoError(t, err)
		assert.Equal(t, "current-migration-version", migrationVersion)
	})

	t.Run("MakeAvailable", func(t *testing.T) {
		// given:
		storageResult := &wdk.TableSettings{
			StorageName:        fixtures.StorageName,
			StorageIdentityKey: fixtures.StorageIdentityKey,
			Chain:              defs.NetworkTestnet,
			MaxOutputScript:    1024,
		}

		mockStorage.EXPECT().
			MakeAvailable(gomock.Any()).
			Return(storageResult, nil)

		// when:
		response, err := client.MakeAvailable(context.Background())

		// then:
		require.NoError(t, err)
		assert.EqualValues(t, storageResult, response)
	})

	t.Run("FindOrInsertUser", func(t *testing.T) {
		// given:
		userIdentityKey := "03f17660f611ce531402a2ce1e070380b6fde57aca211d707bfab27bce42d86beb"

		storageResult := &wdk.FindOrInsertUserResponse{
			User: wdk.TableUser{
				IdentityKey: userIdentityKey,
			},
			IsNew: false,
		}

		// and:
		mockStorage.EXPECT().
			FindOrInsertUser(gomock.Any(), fixtures.StorageIdentityKey).
			Return(storageResult, nil)

		// when:
		response, err := client.FindOrInsertUser(context.Background(), fixtures.StorageIdentityKey)

		// then:
		require.NoError(t, err)
		require.NotNil(t, response)
		assert.EqualValues(t, storageResult, response)
	})

	t.Run("CreateAction", func(t *testing.T) {
		t.Skip("Not implemented yet")
	})
}
