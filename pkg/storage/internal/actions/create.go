package actions

import (
	"log/slog"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/logging"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
)

type create struct {
	logger *slog.Logger
}

func newCreateAction(logger *slog.Logger) *create {
	logger = logging.Child(logger, "createAction")
	return &create{
		logger: logger,
	}
}

func (c *create) Create(auth wdk.AuthID, args wdk.ValidCreateActionArgs) (*wdk.StorageCreateActionResult, error) {
	return nil, nil
}
