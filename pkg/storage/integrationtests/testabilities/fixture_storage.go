package testabilities

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/defs"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/logging"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/server"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/stretchr/testify/require"
)

const (
	StorageServerPrivKey = "8143f5ed6c5b41c3d084d39d49e161d8dde4b50b0685a4e4ac23959d3b8a319b"
	StorageIdentityKey   = "028f2daab7808b79368d99eef1ebc2d35cdafe3932cafe3d83cf17837af034ec29" // that matches StorageServerPrivKey
	StorageName          = "test-storage"
	StorageHandlerName   = "storage_server"
)

type StorageFixture interface {
	GormProvider() *storage.Provider
	StartedRPCServerFor(provider *storage.Provider) (cleanup func())
	RPCClient() (*TestClient, func())

	ServerURL() string
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
	dbConfig.SQLite.ConnectionString = "file:storage-test.db?mode=memory"
	dbConfig.MaxIdleConnections = 1
	dbConfig.MaxOpenConnections = 1

	activeStorage, err := storage.NewGORMProvider(s.logger, dbConfig, defs.NetworkTestnet)
	s.require.NoError(err)

	_, err = activeStorage.Migrate(StorageName, storageIdentityKey)
	s.require.NoError(err)

	return activeStorage
}

func (s *storageFixture) StartedRPCServerFor(provider *storage.Provider) (cleanup func()) {
	s.t.Helper()
	rpcServer := server.NewRPCHandler(s.logger, StorageHandlerName, provider)

	mux := http.NewServeMux()
	rpcServer.Register(mux)

	s.testServer = httptest.NewServer(mux)
	return s.testServer.Close
}

func (s *storageFixture) RPCClient() (client *TestClient, cleanup func()) {
	s.t.Helper()
	client = &TestClient{}
	closer, err := jsonrpc.NewMergeClient(
		context.Background(),
		s.testServer.URL,
		"storage_server",
		[]any{client},
		nil,
		jsonrpc.WithHTTPClient(s.testServer.Client()),
		jsonrpc.WithMethodNamer(jsonrpc.NoNamespaceDecapitalizedMethodNamer),
	)
	s.require.NoError(err)
	return client, closer
}

func (s *storageFixture) ServerURL() string {
	return s.testServer.URL
}

func Given(t testing.TB) StorageFixture {
	return &storageFixture{
		t:       t,
		require: require.New(t),
		logger:  logging.NewTestLogger(t),
	}
}
