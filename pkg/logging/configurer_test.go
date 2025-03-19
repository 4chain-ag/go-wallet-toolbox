package logging_test

import (
	"log/slog"
	"testing"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/logging"
	"github.com/stretchr/testify/require"
)

func TestTextLogger(t *testing.T) {
	// given:
	stringWriter := &logging.TestWriter{}
	logger := logging.New().
		WithLevel(slog.LevelDebug).
		WithHandler(logging.TextHandler, stringWriter).
		Logger()

	// when:
	logger.Debug("debug message")

	// then:
	msg := stringWriter.String()
	require.Contains(t, msg, "time=")
	require.Contains(t, msg, "level=DEBUG")
	require.Contains(t, msg, `msg="debug message"`)
}

func TestJSONLogger(t *testing.T) {
	// given:
	stringWriter := &logging.TestWriter{}
	logger := logging.New().
		WithLevel(slog.LevelDebug).
		WithHandler(logging.JSONHandler, stringWriter).
		Logger()

	// when:
	logger.Debug("debug message")

	// then:
	msg := stringWriter.String()
	require.Contains(t, msg, `"time"`)
	require.Contains(t, msg, `"level":"DEBUG"`)
	require.Contains(t, msg, `"msg":"debug message"`)
}
