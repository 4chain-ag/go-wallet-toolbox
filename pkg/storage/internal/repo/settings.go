package repo

import (
	"fmt"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/defs"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/database/models"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Settings struct {
	db *gorm.DB
}

func NewSettings(db *gorm.DB) *Settings {
	return &Settings{db: db}
}

func (s *Settings) ReadSettings() (*wdk.TableSettings, error) {
	var settings models.Setting
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
	}, nil
}

func (s *Settings) SaveSettings(settings *wdk.TableSettings) error {
	err := s.db.
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&models.Setting{
			StorageIdentityKey: settings.StorageIdentityKey,
			StorageName:        settings.StorageName,
			Chain:              string(settings.Chain),
			MaxOutputScript:    settings.MaxOutputScript,
		}).Error
	if err != nil {
		return fmt.Errorf("failed to save settings: %w", err)
	}

	return nil
}
