package storage_test

import (
	"fmt"
	"testing"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/models"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	t.Run("Test storage", func(t *testing.T) {
		var settingsFromDB []models.Settings
		storageID := "12344666"
		store := storage.NewStorage(&storage.StorageConfig{
			SQLiteConfig: &storage.SQLiteConfig{
				DatabasePath: "./spv-wallet.db",
			},
			LogLevel: "debug",
			Engine:   storage.SQLite,
		}, nil)

		settings := models.Settings{
			StorageIdentityKey: storageID,
			StorageName:        "test-name",
			Chain:              "test",
			DBType:             "SQLite",
			MaxOutputs:         100,
		}
		_ = store.DB.AutoMigrate(&models.Settings{})
		result := store.DB.Create(&settings)

		fmt.Printf("Error: %v, RowsAffected: %d", result.Error, result.RowsAffected)
		fmt.Printf("%+v,", settings)
		res := store.DB.Where("storage_identity_key = ?", storageID).Find(&settingsFromDB)

		require.NoError(t, res.Error)
		require.Equal(t, 1, len(settingsFromDB))
		require.Equal(t, storageID, settingsFromDB[0].StorageIdentityKey)
	})
}
