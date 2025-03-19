package logging_test

import (
	"log/slog"
	"testing"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/logging"
	"github.com/stretchr/testify/require"
)

func TestChildLogger(t *testing.T) {
	// given:
	stringWriter := &logging.TestWriter{}
	logger := logging.New().
		WithLevel(slog.LevelDebug).
		WithHandler(logging.TextHandler, stringWriter).
		Logger()

	// when:
	childLogger := logging.Child(logger, "child")

	// and:
	childLogger.Debug("debug message")

	// then:
	msg := stringWriter.String()
	require.Contains(t, msg, "service=child")
	require.Contains(t, msg, `msg="debug message"`)
}
