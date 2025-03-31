package repo

import "gorm.io/gorm"

type Repositories struct {
	*Migrator
	*Settings
	*Users
	*OutputBaskets
	*UTXOs
}

func NewSQLRepositories(db *gorm.DB) *Repositories {
	repositories := &Repositories{
		Migrator:      NewMigrator(db),
		Settings:      NewSettings(db),
		OutputBaskets: NewOutputBaskets(db),
		UTXOs:         NewUTXOs(db),
	}
	repositories.Users = NewUsers(db, repositories.Settings, repositories.OutputBaskets)

	return repositories
}
