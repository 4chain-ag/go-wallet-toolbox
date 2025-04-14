package testabilities

import (
	"context"
	"log/slog"
	"testing"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/defs"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/testabilities/dbfixtures"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/testabilities/testusers"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	"github.com/stretchr/testify/require"
)

type ProviderFixture interface {
	WithNetwork(network defs.BSVNetwork) ProviderFixture
	WithCommission(commission defs.Commission) ProviderFixture
	WithFeeModel(feeModel defs.FeeModel) ProviderFixture

	GORM() *storage.Provider
	GORMWithCleanDatabase() *storage.Provider
}

type providerFixture struct {
	network    defs.BSVNetwork
	commission defs.Commission
	feeModel   defs.FeeModel

	t       testing.TB
	require *require.Assertions
	logger  *slog.Logger
}

func (p *providerFixture) WithNetwork(network defs.BSVNetwork) ProviderFixture {
	p.network = network
	return p
}

func (p *providerFixture) WithCommission(commission defs.Commission) ProviderFixture {
	p.commission = commission
	return p
}

func (p *providerFixture) WithFeeModel(feeModel defs.FeeModel) ProviderFixture {
	p.feeModel = feeModel
	return p
}

func (p *providerFixture) GORM() *storage.Provider {
	p.t.Helper()
	provider := p.GORMWithCleanDatabase()

	p.seedUsers(provider)

	return provider
}

func (p *providerFixture) GORMWithCleanDatabase() *storage.Provider {
	p.t.Helper()

	storageIdentityKey, err := wdk.IdentityKey(StorageServerPrivKey)
	p.require.NoError(err)

	dbConfig := dbfixtures.DBConfigForTests()

	activeStorage, err := storage.NewGORMProvider(p.logger, storage.GORMProviderConfig{
		DB:         dbConfig,
		Chain:      p.network,
		FeeModel:   p.feeModel,
		Commission: p.commission,
	}, storage.WithFunder(&MockFunder{}))
	p.require.NoError(err)

	_, err = activeStorage.Migrate(context.Background(), StorageName, storageIdentityKey)
	p.require.NoError(err)

	return activeStorage
}

func (p *providerFixture) seedUsers(provider *storage.Provider) {
	for _, user := range testusers.All() {
		res, err := provider.FindOrInsertUser(context.Background(), user.PrivKey)
		p.require.NoError(err)

		user.ID = res.User.UserID
	}
}
