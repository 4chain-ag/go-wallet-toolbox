package main

import (
	"github.com/4chain-ag/go-wallet-toolbox/pkg/logging"
	"log/slog"
	"os"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/infra"
)

func main() {
	rootLogger := logging.New().
		WithLevel(slog.LevelDebug).
		WithHandler(logging.TextHandler, os.Stdout).
		Logger()

	server, err := infra.NewServer(
		infra.WithConfigFile("infra-config.yaml"),
		infra.WithLogger(logging.Child(rootLogger, "infra")),
	)
	if err != nil {
		logging.Fatalf(rootLogger, err, "failed to create server")
	}

	err = server.ListenAndServe()
	if err != nil {
		logging.Fatalf(rootLogger, err, "Server failed")
	}
}
