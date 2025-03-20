package defs

// LogLevel represents different log levels which can be configured.
type LogLevel string

// Supported log levels (based on slog).
const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

// ParseLogLevelStr parses a string into a LogLevel (case-insensitive).
func ParseLogLevelStr(level string) (LogLevel, error) {
	return parseEnumCaseInsensitive(level, LogLevelDebug, LogLevelInfo, LogLevelWarn, LogLevelError)
}

// LogHandler represents different log handler types which can be configured.
type LogHandler string

// Supported handler types (based on slog).
const (
	JSONHandler LogHandler = "json"
	TextHandler LogHandler = "text"
)

// ParseHandlerTypeStr parses a string into a LogHandler (case-insensitive).
func ParseHandlerTypeStr(handlerType string) (LogHandler, error) {
	return parseEnumCaseInsensitive(handlerType, JSONHandler, TextHandler)
}
