package infra

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/config"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/logging"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/server"
)

// Server is a struct that holds the "infra" server configuration
type Server struct {
	Config Config

	logger *slog.Logger
}

// NewServer creates a new server instance with given options, like config file path or a prefix for environment variables
func NewServer(opts ...InitOption) (*Server, error) {
	options := DefaultOptions()
	for _, option := range opts {
		option(&options)
	}

	loader := config.NewLoader(Defaults, options.EnvPrefix)
	if options.ConfigFile != "" {
		err := loader.SetConfigFilePath(options.ConfigFile)
		if err != nil {
			return nil, fmt.Errorf("failed to set config file path: %w", err)
		}
	}
	cfg, err := loader.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	err = cfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	logger := logging.Child(makeLogger(&cfg, &options), "infra")

	return &Server{
		Config: cfg,

		logger: logger,
	}, nil
}

// ListenAndServe starts the JSON-RPC server
func (s *Server) ListenAndServe() error {
	rpcServer := server.NewRPCHandler(s.logger)

	mux := http.NewServeMux()
	rpcServer.Register(mux)

	port := s.Config.HTTPConfig.Port
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

func makeLogger(cfg *Config, options *Options) *slog.Logger {
	if !cfg.Logging.Enabled {
		return logging.DefaultIfNil(options.Logger)
	}

	return logging.New().
		WithLevel(cfg.Logging.Level).
		WithHandler(cfg.Logging.Handler, os.Stdout).
		Logger()
}
