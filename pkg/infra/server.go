package infra

import (
	"fmt"
	"net/http"
	"time"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/config"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/server"
)

type Server struct {
	Config Config
}

func NewServer(opts ...InitOption) (*Server, error) {
	params := DefaultParams()
	for _, option := range opts {
		option(&params)
	}

	loader := config.NewLoader(Defaults, params.EnvPrefix)
	if params.ConfigFile != "" {
		loader.SetConfigFilePath(params.ConfigFile)
	}
	cfg, err := loader.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return &Server{
		Config: cfg,
	}, nil
}

func (s *Server) ListenAndServe() error {
	rpcServer := server.NewRPCHandler()

	mux := http.NewServeMux()
	rpcServer.Register(mux)

	httpServer := &http.Server{
		Addr:              fmt.Sprintf(":%d", s.Config.HTTPConfig.Port),
		Handler:           mux,
		ReadHeaderTimeout: 3 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	err := httpServer.ListenAndServe()
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}
