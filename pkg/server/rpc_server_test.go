package server_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/logging"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/server"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/stretchr/testify/require"
)

func TestRPCServer(t *testing.T) {
	// given server:
	rpcServer := server.NewRPCHandler(logging.New().Nop().Logger())

	mux := http.NewServeMux()
	rpcServer.Register(mux)

	testSrv := httptest.NewServer(mux)
	defer testSrv.Close()

	// and client:
	var client struct {
		MakeAvailable func() server.TableSettings
	}
	closer, err := jsonrpc.NewMergeClient(
		context.Background(),
		testSrv.URL,
		"SimpleServerHandler",
		[]any{&client},
		nil,
		jsonrpc.WithMethodNamer(jsonrpc.NoNamespaceDecapitalizedMethodNamer),
	)
	require.NoError(t, err)
	defer closer()

	// when:
	tableSettings := client.MakeAvailable()

	// then:
	require.NotEmpty(t, tableSettings.StorageIdentityKey)
}
