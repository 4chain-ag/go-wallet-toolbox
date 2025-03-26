package actions

import (
	"log/slog"
)

type Actions struct {
	*create
}

func New(logger *slog.Logger, funder Funder) *Actions {
	return &Actions{
		create: newCreateAction(logger, funder),
	}
}
