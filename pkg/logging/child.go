package logging

import "log/slog"

const (
	// ServiceKey is the key used to store the service name in the logger context.
	ServiceKey = "service"
)

func Child(logger *slog.Logger, serviceName string) *slog.Logger {
	return logger.With(
		slog.String(ServiceKey, serviceName),
	)
}
