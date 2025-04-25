package testabilities

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/defs"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/fixtures"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/logging"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/mocks"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/satoshi"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/database"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/database/models"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/server"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/testabilities/dbfixtures"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/testabilities/testusers"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	txtestabilities "github.com/bsv-blockchain/universal-test-vectors/pkg/testabilities"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type StorageFixture interface {
	Provider() ProviderFixture

	StartedRPCServerFor(provider wdk.WalletStorageWriter) (cleanup func())
	RPCClient() (*storage.WalletStorageWriterClient, func())

	MockProvider() *mocks.MockWalletStorageWriter

	Faucet(activeStorage *storage.Provider, user testusers.User) FaucetFixture
}

type FaucetFixture interface {
	TopUp(satoshis satoshi.Value) (txtestabilities.TransactionSpec, *models.UserUTXO)
}

type storageFixture struct {
	t          testing.TB
	require    *require.Assertions
	logger     *slog.Logger
	testServer *httptest.Server
	db         *database.Database
}

func (s *storageFixture) StartedRPCServerFor(provider wdk.WalletStorageWriter) (cleanup func()) {
	s.t.Helper()
	rpcServer := server.NewRPCHandler(s.logger, fixtures.StorageHandlerName, provider)

	mux := http.NewServeMux()
	rpcServer.Register(mux)

	s.testServer = httptest.NewServer(mux)
	return s.testServer.Close
}

func (s *storageFixture) RPCClient() (client *storage.WalletStorageWriterClient, cleanup func()) {
	s.t.Helper()
	client, cleanup, err := storage.NewClient(s.testServer.URL, storage.WithHttpClient(s.testServer.Client()))
	s.require.NoError(err)
	return client, cleanup
}

func (s *storageFixture) MockProvider() *mocks.MockWalletStorageWriter {
	s.t.Helper()
	ctrl := gomock.NewController(s.t)

	return mocks.NewMockWalletStorageWriter(ctrl)
}

func (s *storageFixture) Provider() ProviderFixture {
	s.t.Helper()
	return &providerFixture{
		t:       s.t,
		require: s.require,
		logger:  s.logger,
		db:      s.db,

		network:    defs.NetworkTestnet,
		commission: defs.Commission{},
		feeModel:   defs.DefaultFeeModel(),
	}
}

func (s *storageFixture) Faucet(activeStorage *storage.Provider, user testusers.User) FaucetFixture {
	s.t.Helper()
	ctx := context.Background()

	_, err := activeStorage.FindOrInsertUser(ctx, user.PrivKey)
	s.require.NoError(err)

	basket, err := s.db.CreateRepositories().
		FindBasketByName(context.Background(), user.ID, wdk.BasketNameForChange)
	require.NoError(s.t, err)

	return &faucetFixture{
		t:        s.t,
		user:     user,
		db:       s.db,
		basketID: basket.BasketID,
	}
}

func Given(t testing.TB) StorageFixture {
	db, _ := dbfixtures.TestDatabase(t)
	return &storageFixture{
		t:       t,
		require: require.New(t),
		logger:  logging.NewTestLogger(t),
		db:      db,
	}
}
