package actions

import (
	"log/slog"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/defs"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/repo"
)

type Actions struct {
	*create
	*internalize
}

func New(logger *slog.Logger, funder Funder, commission defs.Commission, repos *repo.Repositories) *Actions {
	return &Actions{
		create:      newCreateAction(logger, funder, commission, repos.OutputBaskets, repos.Transactions, repos.Outputs),
		internalize: newInternalizeAction(logger),
	}
}
