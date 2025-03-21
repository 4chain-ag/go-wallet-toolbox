package repo

import (
	"fmt"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/database/models"
	"gorm.io/gorm"
)

type Migrator struct {
	db *gorm.DB
}

func NewMigrator(db *gorm.DB) *Migrator {
	return &Migrator{db: db}
}

func (m *Migrator) Migrate() error {
	err := m.db.AutoMigrate(&models.Settings{})
	if err != nil {
		return fmt.Errorf("failed to migrate settings: %w", err)
	}
	
	return nil
}
