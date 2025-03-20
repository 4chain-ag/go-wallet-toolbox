package logging

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/defs"
)

const (
	// ServiceKey is the key used to store the service name in the logger attrs.
	ServiceKey = "service"
)

var strLevelToSlog = map[defs.LogLevel]slog.Level{
	defs.LogLevelDebug: slog.LevelDebug,
	defs.LogLevelInfo:  slog.LevelInfo,
	defs.LogLevelWarn:  slog.LevelWarn,
	defs.LogLevelError: slog.LevelError,
}

// Child returns a new logger with the given service name added to the logger attrs.
func Child(logger *slog.Logger, serviceName string) *slog.Logger {
	return DefaultIfNil(logger).With(
		slog.String(ServiceKey, serviceName),
	)
}

func Errorf(logger *slog.Logger, err error, format string, args ...any) {
	logger.Error(fmt.Sprintf(format, args...), slog.String("error", err.Error()))
}

// Fatalf logs the error and exits the program.
func Fatalf(logger *slog.Logger, err error, msg string, args ...any) {
	Errorf(logger, err, "Fatal error: "+msg, args...)
	os.Exit(1)
}

// DefaultIfNil returns the default logger if the given logger is nil.
func DefaultIfNil(logger *slog.Logger) *slog.Logger {
	if logger == nil {
		return slog.Default()
	}
	return logger
}
