package server

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/logging"
	"github.com/filecoin-project/go-jsonrpc"
)

type RPCServer struct {
	Handler *jsonrpc.RPCServer
}

func NewRPCHandler(parentLogger *slog.Logger, name string, handler any) *RPCServer {
	logger := logging.Child(parentLogger, "rpc_server")

	rpcServer := jsonrpc.NewServer(
		jsonrpc.WithServerMethodNameFormatter(jsonrpc.NewMethodNameFormatter(false, jsonrpc.LowerFirstCharCase)),
		jsonrpc.WithTracer(tracer(logger)),
	)

	rpcServer.Register(name, handler)

	return &RPCServer{
		Handler: rpcServer,
	}
}

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
