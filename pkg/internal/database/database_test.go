package database_test

import (
	"testing"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/database"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/database/models"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	t.Run("Test database", func(t *testing.T) {
		var settingsFromDB models.Settings
		storageID := "12344666"
		storageDB, err := database.NewDatabase(&database.Config{
			SQLiteConfig: database.SQLiteConfig{
				ConnectionString: "./spv-wallet.db",
			},
			Engine: database.SQLite,
		}, nil)
		require.NoError(t, err)

		settings := models.Settings{
			StorageIdentityKey: storageID,
			StorageName:        "test-name",
			Chain:              "test",
			DBType:             "SQLite",
			MaxOutputs:         100,
		}
		_ = storageDB.DB.AutoMigrate(&models.Settings{})
		storageDB.DB.Create(&settings)

		res := storageDB.DB.Where("storage_identity_key = ?", storageID).First(&settingsFromDB)

		require.NoError(t, res.Error)
		require.Equal(t, storageID, settingsFromDB.StorageIdentityKey)
	})
}
