package testabilities

import (
	"testing"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/defs"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/logging"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/actions/funder"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/database"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/database/models"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/testabilities/dbfixtures"
	"github.com/stretchr/testify/require"
)

type FunderFixture interface {
	NewFunderService() *funder.SQL
	UTXO() UserUTXOFixture
}

var feeModel = defs.FeeModel{
	Type:  defs.SatPerKB,
	Value: 1,
}

type funderFixture struct {
	t            testing.TB
	db           *database.Database
	createdUTXOs []*models.UserUTXO
}

func newFixture(t testing.TB) (given FunderFixture, cleanup func()) {
	db, dbCleanup := dbfixtures.TestDatabase(t)
	return &funderFixture{
		t:            t,
		db:           db,
		createdUTXOs: make([]*models.UserUTXO, 0),
	}, dbCleanup
}

func (f *funderFixture) NewFunderService() *funder.SQL {
	repo := f.db.CreateRepositories().UTXOs
	return funder.NewSQL(logging.NewTestLogger(f.t), repo, feeModel)
}

func (f *funderFixture) UTXO() UserUTXOFixture {
	index := uint(len(f.createdUTXOs))
	return newUtxoFixture(f.t, f, index)
}

func (f *funderFixture) Save(utxo *models.UserUTXO) {
	err := f.db.DB.Create(&utxo).Error
	require.NoError(f.t, err)
	f.createdUTXOs = append(f.createdUTXOs, utxo)
}
