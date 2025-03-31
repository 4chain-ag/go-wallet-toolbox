package actions

import (
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/repo"
	"log/slog"
)

type Actions struct {
	*create
}

func New(logger *slog.Logger, funder Funder, repos *repo.Repositories) *Actions {
	return &Actions{
		create: newCreateAction(logger, funder, repos.OutputBaskets),
	}
}
