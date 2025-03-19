package logging

import (
	"context"
	"fmt"
	"log/slog"
)

func Sprintf(logger *slog.Logger, level slog.Level, format string, args ...any) {
	logger.Log(context.Background(), level, fmt.Sprintf(format, args...))
}
