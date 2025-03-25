package actions

import (
	"log/slog"
)

type Actions struct {
	*create
}

func New(logger *slog.Logger) *Actions {
	return &Actions{
		create: newCreateAction(logger),
	}
}
