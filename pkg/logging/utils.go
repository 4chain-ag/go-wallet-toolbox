package logging

import (
	"context"
	"fmt"
	"log/slog"
)

const (
	// ServiceKey is the key used to store the service name in the logger attrs.
	ServiceKey = "service"
)

// Child returns a new logger with the given service name added to the logger attrs.
func Child(logger *slog.Logger, serviceName string) *slog.Logger {
	return logger.With(
		slog.String(ServiceKey, serviceName),
	)
}

// Sprintf logs a message with the given level and arguments.
func Sprintf(logger *slog.Logger, level slog.Level, format string, args ...any) {
	logger.Log(context.Background(), level, fmt.Sprintf(format, args...))
}

// NopIfNil returns a new NOP logger if the given logger is nil, otherwise returns the given logger.
// NOTE: NOP logger discards all logs.
func NopIfNil(logger *slog.Logger) *slog.Logger {
	if logger == nil {
		return New().Nop().Logger()
	}
	return logger
}
