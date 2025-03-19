package logging_test

import (
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/logging"
	"github.com/stretchr/testify/require"
	"log/slog"
	"strings"
	"testing"
)

func TestTextLogger(t *testing.T) {
	// given:
	stringWriter := &StringWriter{}
	logger := logging.New().
		WithLevel(slog.LevelDebug).
		WithHandler(logging.TextHandler, stringWriter).
		Logger()

	// when:
	logger.Debug("debug message")

	// then:
	msg := stringWriter.builder.String()
	require.Contains(t, msg, "time=")
	require.Contains(t, msg, "level=DEBUG")
	require.Contains(t, msg, `msg="debug message"`)
}

func TestJSONLogger(t *testing.T) {
	// given:
	stringWriter := &StringWriter{}
	logger := logging.New().
		WithLevel(slog.LevelDebug).
		WithHandler(logging.JSONHandler, stringWriter).
		Logger()

	// when:
	logger.Debug("debug message")

	// then:
	msg := stringWriter.builder.String()
	require.Contains(t, msg, `"time"`)
	require.Contains(t, msg, `"level":"DEBUG"`)
	require.Contains(t, msg, `"msg":"debug message"`)
}

type StringWriter struct {
	builder strings.Builder
}

func (w *StringWriter) Write(p []byte) (n int, err error) {
	return w.builder.Write(p)
}
