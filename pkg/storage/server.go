package storage

import (
	"fmt"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/logging"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/server"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	"log/slog"
	"net/http"
	"time"
)

type Server struct {
	provider wdk.StorageProvider
	logger   *slog.Logger
	options  ServerOptions
}

func NewServer(logger *slog.Logger, storage wdk.StorageProvider, opts ...ServerOption) *Server {
	options := defaultServerOptions()
	for _, opt := range opts {
		opt(&options)
	}

	return &Server{
		provider: storage,
		logger:   logging.Child(logger, "storage_server"),
	}
}

// Start starts the server
// NOTE: This method is blocking
func (s *Server) Start() error {
	rpcServer := server.NewRPCHandler(s.logger, "remote_storage", s.provider)

	mux := http.NewServeMux()
	rpcServer.Register(mux)

	port := s.options.Port
	httpServer := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           mux,
		ReadHeaderTimeout: 3 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	s.logger.Info("Listening...", slog.Any("port", port))
	err := httpServer.ListenAndServe()
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	return nil
}
