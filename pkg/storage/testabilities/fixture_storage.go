package testabilities

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/defs"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/logging"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/mocks"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/server"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	StorageServerPrivKey = "8143f5ed6c5b41c3d084d39d49e161d8dde4b50b0685a4e4ac23959d3b8a319b"
	StorageIdentityKey   = "028f2daab7808b79368d99eef1ebc2d35cdafe3932cafe3d83cf17837af034ec29" // that matches StorageServerPrivKey
	StorageName          = "test-storage"
	StorageHandlerName   = "storage_server"
)

type StorageFixture interface {
	GormProvider() *storage.Provider
	StartedRPCServerFor(provider wdk.WalletStorageWriter) (cleanup func())
	RPCClient() (*wdk.WalletStorageWriterClient, func())
	MockProvider() *mocks.MockWalletStorageWriter
}

type storageFixture struct {
	t          testing.TB
	require    *require.Assertions
	logger     *slog.Logger
	testServer *httptest.Server
}

func (s *storageFixture) GormProvider() *storage.Provider {
	s.t.Helper()

	storageIdentityKey, err := wdk.IdentityKey(StorageServerPrivKey)
	s.require.NoError(err)

	dbConfig := defs.DefaultDBConfig()
	dbConfig.SQLite.ConnectionString = "file:storage.test.sqlite?mode=memory"
	dbConfig.MaxIdleConnections = 1
	dbConfig.MaxOpenConnections = 1

	activeStorage, err := storage.NewGORMProvider(s.logger, dbConfig, defs.NetworkTestnet)
	s.require.NoError(err)

	_, err = activeStorage.Migrate(StorageName, storageIdentityKey)
	s.require.NoError(err)

	return activeStorage
}

func (s *storageFixture) StartedRPCServerFor(provider wdk.WalletStorageWriter) (cleanup func()) {
	s.t.Helper()
	rpcServer := server.NewRPCHandler(s.logger, StorageHandlerName, provider)

	mux := http.NewServeMux()
	rpcServer.Register(mux)

	s.testServer = httptest.NewServer(mux)
	return s.testServer.Close
}

func (s *storageFixture) RPCClient() (client *wdk.WalletStorageWriterClient, cleanup func()) {
	s.t.Helper()
	client, cleanup, err := wdk.NewClient(s.testServer.URL, wdk.WithHttpClient(s.testServer.Client()))
	s.require.NoError(err)
	return client, cleanup
}

func (s *storageFixture) MockProvider() *mocks.MockWalletStorageWriter {
	s.t.Helper()
	ctrl := gomock.NewController(s.t)

	return mocks.NewMockWalletStorageWriter(ctrl)
}

func Given(t testing.TB) StorageFixture {
	return &storageFixture{
		t:       t,
		require: require.New(t),
		logger:  logging.NewTestLogger(t),
	}
}
