package actions

import (
	"log/slog"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/repo"
)

type Actions struct {
	*create
}

func New(logger *slog.Logger, funder Funder, repos *repo.Repositories) *Actions {
	return &Actions{
		create: newCreateAction(logger, funder, repos.OutputBaskets),
	}
}
