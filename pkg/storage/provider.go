package storage

import (
	"fmt"
	"log/slog"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/defs"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/database"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/repo"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
)

// Repository is an interface for the actual storage repository.
type Repository interface {
	Migrate() error
	ReadSettings() (*wdk.SettingsDTO, error)
	SaveSettings(settings *wdk.SettingsDTO) error
}

// Provider is a storage provider.
type Provider struct {
	Chain defs.BSVNetwork

	settings *wdk.SettingsDTO
	repo     Repository
}

// NewGORMProvider creates a new storage provider with GORM repository.
func NewGORMProvider(logger *slog.Logger, dbConfig defs.Database, chain defs.BSVNetwork) (*Provider, error) {
	db, err := database.NewDatabase(dbConfig, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create database: %w", err)
	}

	return &Provider{
		Chain: chain,
		repo:  repo.NewRepositories(db.DB),
	}, nil
}

// Migrate migrates the storage and saves the settings.
func (p *Provider) Migrate(storageName, storageIdentityKey string) (string, error) {
	err := p.repo.Migrate()
	if err != nil {
		return "", fmt.Errorf("failed to migrate: %w", err)
	}

	// TODO: what if p.Chain != Chain from DB?

	err = p.repo.SaveSettings(&wdk.SettingsDTO{
		StorageIdentityKey: storageIdentityKey,
		StorageName:        storageName,
		Chain:              p.Chain,
		MaxOutputScript:    DefaultMaxScriptLength,
	})
	if err != nil {
		return "", fmt.Errorf("failed to save settings: %w", err)
	}

	// NOTE: GORM automigrate does not support db versioning
	// from-kt: In TS version I can't find any usage of returned version
	version := "auto-migrated"

	return version, nil
}

// MakeAvailable reads the settings and makes them available.
func (p *Provider) MakeAvailable() (*wdk.SettingsDTO, error) {
	settings, err := p.repo.ReadSettings()
	if err != nil {
		return nil, fmt.Errorf("failed to read settings: %w", err)
	}

	p.settings = settings
	return settings, nil
}
