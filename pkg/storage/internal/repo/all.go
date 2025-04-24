package repo

import "gorm.io/gorm"

type Repositories struct {
	*Migrator
	*Settings
	*Users
	*OutputBaskets
	*Certificates
	*UTXOs
	*Transactions
	*Outputs
}

func NewSQLRepositories(db *gorm.DB) *Repositories {
	repositories := &Repositories{
		Migrator:      NewMigrator(db),
		Settings:      NewSettings(db),
		OutputBaskets: NewOutputBaskets(db),
		Certificates:  NewCertificates(db),
		UTXOs:         NewUTXOs(db),
		Transactions:  NewTransactions(db),
		Outputs:       NewOutputs(db),
	}
	repositories.Users = NewUsers(db, repositories.Settings, repositories.OutputBaskets)

	return repositories
}
