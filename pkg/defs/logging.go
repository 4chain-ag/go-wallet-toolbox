package defs

type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

func ParseLogLevelStr(level string) (LogLevel, error) {
	return parseEnumCaseInsensitive(level, LogLevelDebug, LogLevelInfo, LogLevelWarn, LogLevelError)
}

type LogHandler string

// Supported handler types.
const (
	JSONHandler LogHandler = "json"
	TextHandler LogHandler = "text"
)

func ParseHandlerTypeStr(handlerType string) (LogHandler, error) {
	return parseEnumCaseInsensitive(handlerType, JSONHandler, TextHandler)
}
