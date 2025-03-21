package integrationtests_test

import (
	"context"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/defs"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/logging"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/integrationtests/testabilities"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/server"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMakeAvailable(t *testing.T) {
	// given:
	logger := logging.NewTestLogger(t)

	storageIdentityKey, err := wdk.IdentityKey(testabilities.StorageIdentityKey)
	require.NoError(t, err)

	dbConfig := defs.DefaultDBConfig()
	dbConfig.SQLite.ConnectionString = "file:storage-test.db?mode=memory"
	dbConfig.MaxIdleConnections = 1
	dbConfig.MaxOpenConnections = 1

	activeStorage, err := storage.NewGORMProvider(logger, dbConfig, defs.NetworkTestnet)
	require.NoError(t, err)

	_, err = activeStorage.Migrate("test", storageIdentityKey)
	require.NoError(t, err)

	// given server:
	rpcServer := server.NewRPCHandler(logger, "storage_server", activeStorage)

	mux := http.NewServeMux()
	rpcServer.Register(mux)

	testSrv := httptest.NewServer(mux)
	defer testSrv.Close()

	// and client:
	var client struct {
		MakeAvailable func() (*wdk.TableSettings, error)
	}
	closer, err := jsonrpc.NewMergeClient(
		context.Background(),
		testSrv.URL,
		"storage_server",
		[]any{&client},
		nil,
		jsonrpc.WithMethodNamer(jsonrpc.NoNamespaceDecapitalizedMethodNamer),
	)
	require.NoError(t, err)
	defer closer()

	// when:
	tableSettings, err := client.MakeAvailable()

	// then:
	require.NoError(t, err)
	require.Equal(t, "test", tableSettings.StorageName)
	require.Equal(t, storageIdentityKey, tableSettings.StorageIdentityKey)
	require.Equal(t, defs.NetworkTestnet, tableSettings.Chain)
	require.Equal(t, 1024, tableSettings.MaxOutputScript)
}
