package infra

import (
	"fmt"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/server"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/config"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/logging"
)

// Server is a struct that holds the "infra" server configuration
type Server struct {
	Config Config

	logger  *slog.Logger
	storage *storage.Provider
}

// NewServer creates a new server instance with given options, like config file path or a prefix for environment variables
func NewServer(opts ...InitOption) (*Server, error) {
	options := defaultOptions()
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

	storageIdentityKey, err := wdk.IdentityKey(cfg.ServerPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage identity key: %w", err)
	}

	activeStorage, err := storage.NewGORMProvider(logger, cfg.DBConfig, cfg.BSVNetwork)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage provider: %w", err)
	}

	// TODO: Get database name from somewhere
	_, err = activeStorage.Migrate("local", storageIdentityKey)
	if err != nil {
		return nil, fmt.Errorf("failed to migrate storage: %w", err)
	}

	return &Server{
		Config: cfg,

		logger: logger,
	}, nil
}

// ListenAndServe starts the JSON-RPC server
func (s *Server) ListenAndServe() error {
	rpcServer := server.NewRPCHandler(s.logger, "remote_storage", s.storage)

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
	if options.Logger != nil {
		return options.Logger
	}

	if !cfg.Logging.Enabled {
		return logging.New().Nop().Logger()
	}

	return logging.New().
		WithLevel(cfg.Logging.Level).
		WithHandler(cfg.Logging.Handler, os.Stdout).
		Logger()
}
