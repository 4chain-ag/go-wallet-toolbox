package server

import (
	"context"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRPCServer(t *testing.T) {
	// given server:
	rpcServer := NewRPCHandler()

	mux := http.NewServeMux()
	rpcServer.Register(mux)

	testSrv := httptest.NewServer(mux)
	defer testSrv.Close()

	// and client:
	var client struct {
		MakeAvailable func() TableSettings
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
