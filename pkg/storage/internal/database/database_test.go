package database_test

import (
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/repo"
	"testing"
	"time"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/defs"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/database"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/database/models"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	t.Run("Find or create user", func(t *testing.T) {
		// given:
		storageID := "12344666"
		storageDB, err := database.NewDatabase(defs.Database{
			SQLite: defs.SQLite{
				ConnectionString: "file::memory:",
			},
			MaxConnectionTime:     360 * time.Second,
			MaxIdleConnections:    1,
			MaxConnectionIdleTime: 360 * time.Second,
			MaxOpenConnections:    1,
			Engine:                defs.DBTypeSQLite,
		}, nil)
		require.NoError(t, err)

		// and:
		err = storageDB.DB.AutoMigrate(&models.Settings{}, &models.User{})
		require.NoError(t, err)

		// and:
		settings := models.Settings{
			StorageIdentityKey: storageID,
			StorageName:        "test-name",
			Chain:              "test",
			MaxOutputScript:    100,
		}
		storageDB.DB.Create(&settings)

		// and:
		repository := repo.NewRepositories(storageDB.DB)

		//when:
		usr, err := repository.FindOrCreateUser("028f2daab7808b79368d99eef1ebc2d35cdafe3932cafe3d83cf17837af034ec29")
		require.NoError(t, err)
		require.Equal(t, usr.User.ActiveStorage, settings.StorageIdentityKey)

		usr1, err1 := repository.FindOrCreateUser("028f2daab7808b79368d99eef1ebc2d35cdafe3932cafe3d83cf17837af034ec29")
		require.NoError(t, err1)
		require.Equal(t, usr1.User.ActiveStorage, settings.StorageIdentityKey)
	})

	t.Run("Test database", func(t *testing.T) {
		// given:
		var settingsFromDB models.Settings
		storageID := "12344666"
		settings := models.Settings{
			StorageIdentityKey: storageID,
			StorageName:        "test-name",
			Chain:              "test",
			MaxOutputScript:    100,
		}

		// and:
		storageDB, err := database.NewDatabase(defs.Database{
			SQLite: defs.SQLite{
				ConnectionString: "file::memory:",
			},
			MaxConnectionTime:     5 * time.Second,
			MaxIdleConnections:    1,
			MaxConnectionIdleTime: 5 * time.Second,
			MaxOpenConnections:    1,
			Engine:                defs.DBTypeSQLite,
		}, nil)
		require.NoError(t, err)

		// when:
		err = storageDB.DB.AutoMigrate(&models.Settings{})
		require.NoError(t, err)
		storageDB.DB.Create(&settings)
		res := storageDB.DB.Where("storage_identity_key = ?", storageID).First(&settingsFromDB)

		// then:
		require.NoError(t, res.Error)
		require.Equal(t, storageID, settingsFromDB.StorageIdentityKey)
	})
}
