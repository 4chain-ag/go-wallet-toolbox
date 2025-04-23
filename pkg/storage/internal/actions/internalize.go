package actions

import (
	"context"
	"log/slog"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/logging"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
)

type internalize struct {
	logger *slog.Logger
}

func newInternalizeAction(logger *slog.Logger) *internalize {
	logger = logging.Child(logger, "createAction")
	return &internalize{
		logger: logger,
	}
}

func (in *internalize) Internalize(ctx context.Context, userID int, args *wdk.InternalizeActionArgs) (*wdk.InternalizeActionResult, error) {
	panic("not implemented")
}
