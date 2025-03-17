package server

import (
	"fmt"
	"github.com/filecoin-project/go-jsonrpc"
	"net/http"
)

type RPCServer struct {
	Handler *jsonrpc.RPCServer
}

func NewRPCHandler() *RPCServer {
	rpcServer := jsonrpc.NewServer(jsonrpc.WithServerMethodNamer(jsonrpc.NoNamespaceDecapitalizedMethodNamer))

	// create a handler instance and register it
	serverHandler := &Handler{}
	rpcServer.Register("Handler", serverHandler)

	return &RPCServer{
		Handler: rpcServer,
	}
}

func (s *RPCServer) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /{$}", s.Handler.ServeHTTP)
	mux.HandleFunc("POST /.well-known/auth", s.handleAuth)
}

func (s *RPCServer) handleAuth(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handleAuth")
	_ = r.Body.Close()

	w.WriteHeader(http.StatusInternalServerError)

	//_, _ = w.Write([]byte(`{}`))
	//w.WriteHeader(http.StatusOK)
}
