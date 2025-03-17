package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/server"
)

func main() {
	rpcServer := server.NewRPCHandler()

	mux := http.NewServeMux()
	rpcServer.Register(mux)

	fmt.Println("Listening on :8101")

	s := &http.Server{
		Addr:              ":8101",
		Handler:           mux,
		ReadHeaderTimeout: 3 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	err := s.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
