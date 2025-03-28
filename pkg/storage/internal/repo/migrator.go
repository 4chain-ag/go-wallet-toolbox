package repo

import (
	"fmt"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/database/models"
	"gorm.io/gorm"
)

type Migrator struct {
	db *gorm.DB
}

func NewMigrator(db *gorm.DB) *Migrator {
	return &Migrator{db: db}
}

func (m *Migrator) Migrate() error {
	err := m.db.AutoMigrate(
		&models.Settings{},
		&models.User{},
		&models.OutputBaskets{},
		&models.CertificateField{},
		&models.Certificate{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate settings: %w", err)
	}

	return nil
}
