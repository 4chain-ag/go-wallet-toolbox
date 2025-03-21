package server

import (
	"fmt"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/logging"
	"github.com/filecoin-project/go-jsonrpc"
	"log/slog"
	"net/http"
)

// RPCServer is a JSON-RPC server
type RPCServer struct {
	Handler *jsonrpc.RPCServer
}

// NewRPCHandler creates a new RPCServer instance
func NewRPCHandler(parentLogger *slog.Logger, name string, handler any) *RPCServer {
	logger := logging.Child(parentLogger, "rpc_server")

	rpcServer := jsonrpc.NewServer(
		jsonrpc.WithServerMethodNamer(jsonrpc.NoNamespaceDecapitalizedMethodNamer),
		jsonrpc.WithTracer(tracer(logger)),
	)

	rpcServer.Register(name, handler)

	return &RPCServer{
		Handler: rpcServer,
	}
}

// Register registers the RPCServer with the provided ServeMux
func (s *RPCServer) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /{$}", s.Handler.ServeHTTP)
	mux.HandleFunc("POST /.well-known/auth", s.handleAuth) //fixme: this is a workaround to pass the client to the next step, it will be handled by the auth middleware
}

func (s *RPCServer) handleAuth(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handleAuth")
	_ = r.Body.Close()

	// from-kt: this is a workaround to pass the client to the next step
	w.WriteHeader(http.StatusInternalServerError)
}
