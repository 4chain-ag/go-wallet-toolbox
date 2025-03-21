package storage

import (
	"fmt"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/defs"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/database"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/repo"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	"log/slog"
)

type Repository interface {
	Migrate() error
	ReadSettings() (*wdk.TableSettings, error)
	SaveSettings(settings *wdk.TableSettings) error
}

type Provider struct {
	Chain defs.BSVNetwork

	settings *wdk.TableSettings
	repo     Repository
}

func NewSQLiteProvider(logger *slog.Logger, chain defs.BSVNetwork, connectionString string) (*Provider, error) {
	db, err := database.NewDatabase(&database.Config{
		Engine: defs.DBTypeSQLite,
		SQLiteConfig: database.SQLiteConfig{
			ConnectionString: connectionString,
		},
		// TODO: Do it differently after Damian's PR merged
		MaxIdleConnections: 1,
		MaxOpenConnections: 1,
	}, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create database: %w", err)
	}

	return &Provider{
		Chain: chain,
		repo:  repo.NewRepositories(db.DB),
	}, nil
}

func (p *Provider) Migrate(storageName, storageIdentityKey string) (string, error) {
	err := p.repo.Migrate()
	if err != nil {
		return "", fmt.Errorf("failed to migrate: %w", err)
	}

	// TODO: what if p.Chain != Chain from DB?

	err = p.repo.SaveSettings(&wdk.TableSettings{
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

func (p *Provider) MakeAvailable() (*wdk.TableSettings, error) {
	settings, err := p.repo.ReadSettings()
	if err != nil {
		return nil, fmt.Errorf("failed to read settings: %w", err)
	}

	p.settings = settings
	return settings, nil
}
