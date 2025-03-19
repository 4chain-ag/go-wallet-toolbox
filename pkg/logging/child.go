package logging

import "log/slog"

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
