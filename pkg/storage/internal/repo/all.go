package repo

import "gorm.io/gorm"

type Repositories struct {
	*Migrator
	*Settings
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		Migrator: NewMigrator(db),
		Settings: NewSettings(db),
	}
}
