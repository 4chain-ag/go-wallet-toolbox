package repo

import (
	"fmt"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/defs"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/database/models"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	"gorm.io/gorm"
)

type Settings struct {
	db *gorm.DB
}

func NewSettings(db *gorm.DB) *Settings {
	return &Settings{db: db}
}

func (s *Settings) ReadSettings() (*wdk.TableSettings, error) {
	var settings models.Settings
	err := s.db.First(&settings).Error
	if err != nil {
		return nil, fmt.Errorf("failed to read settings: %w", err)
	}

	chain, err := defs.ParseBSVNetworkStr(settings.Chain)
	if err != nil {
		return nil, fmt.Errorf("failed to parse chain from settings: %w", err)
	}

	return &wdk.TableSettings{
		StorageIdentityKey: settings.StorageIdentityKey,
		StorageName:        settings.StorageName,
		CreatedAt:          settings.CreatedAt,
		UpdatedAt:          settings.UpdatedAt,
		Chain:              chain,
		MaxOutputScript:    settings.MaxOutputScript,

		//DbType:             settings.DbType, //from-kt: returning DB type what is used on the server side is a security risk
	}, nil
}

func (s *Settings) SaveSettings(settings *wdk.TableSettings) error {
	err := s.db.Create(&models.Settings{
		StorageIdentityKey: settings.StorageIdentityKey,
		StorageName:        settings.StorageName,
		Chain:              string(settings.Chain),
		MaxOutputScript:    settings.MaxOutputScript,
		//DbType:             settings.DbType, //from-kt: DB type should be determined by the server side
	}).Error
	if err != nil {
		return fmt.Errorf("failed to save settings: %w", err)
	}

	return nil
}
