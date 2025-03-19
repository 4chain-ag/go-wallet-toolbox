package server

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/filecoin-project/go-jsonrpc"
)

// RPCServer is a JSON-RPC server
type RPCServer struct {
	Handler *jsonrpc.RPCServer

	logger *slog.Logger
}

// NewRPCHandler creates a new RPCServer instance
func NewRPCHandler(logger *slog.Logger) *RPCServer {
	rpcServer := jsonrpc.NewServer(
		jsonrpc.WithServerMethodNamer(jsonrpc.NoNamespaceDecapitalizedMethodNamer),
		jsonrpc.WithTracer(tracer(logger)),
	)

	// create a handler instance and register it
	serverHandler := &Handler{}
	rpcServer.Register("Handler", serverHandler)

	return &RPCServer{
		Handler: rpcServer,

		logger: logger,
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
