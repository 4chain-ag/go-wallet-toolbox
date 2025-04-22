package testabilities

import (
	"testing"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/database"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/testabilities/dbfixtures"
)

func New(t testing.TB) (given FunderFixture, then FunderAssertion, cleanup func()) {
	db, cleanup := dbfixtures.TestDatabase(t)
	given, then = NewWithDatabase(t, db)
	return
}

func NewWithDatabase(t testing.TB, db *database.Database) (given FunderFixture, then FunderAssertion) {
	given = newFixture(t, db)
	fixture := given.(*funderFixture)
	then = newFunderAssertion(t, fixture)
	return
}
