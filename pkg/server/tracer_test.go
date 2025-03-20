package server_test

import (
	"context"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/defs"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/logging"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/server"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTracer(t *testing.T) {
	// TODO: Reorganize the tests when testabilities are introduced
	// given:
	testWriter := logging.TestWriter{}
	logger := logging.New().WithLevel(defs.LogLevelDebug).WithHandler(defs.TextHandler, &testWriter).Logger()

	// given server:
	rpcServer := server.NewRPCHandler(logger)

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
	_ = client.MakeAvailable()

	// then:
	msg := testWriter.String()
	assert.Contains(t, msg, "time=")
	assert.Contains(t, msg, "level=INFO")
	assert.Contains(t, msg, `msg="Handling RPC call"`)
	assert.Contains(t, msg, `method=makeAvailable`)
}
