package logging

import (
	"context"
	"fmt"
	"log/slog"
	"os"
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

// Fatalf logs the error and exits the program.
func Fatalf(logger *slog.Logger, err error, msg string, args ...any) {
	logger.Error("Fatal error: "+fmt.Sprintf(msg, args...), slog.String("error", err.Error()))
	os.Exit(1)
}

// NopIfNil returns a new NOP logger if the given logger is nil, otherwise returns the given logger.
// NOTE: NOP logger discards all logs.
func NopIfNil(logger *slog.Logger) *slog.Logger {
	if logger == nil {
		return New().Nop().Logger()
	}
	return logger
}
