package logging_test

import (
	"log/slog"
	"testing"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/logging"
	"github.com/stretchr/testify/require"
)

func TestSprintf(t *testing.T) {
	// given:
	stringWriter := &logging.TestWriter{}
	logger := logging.New().
		WithLevel(slog.LevelDebug).
		WithHandler(logging.TextHandler, stringWriter).
		Logger()

	// when:
	logging.Sprintf(logger, slog.LevelInfo, "info message: %d", 123)

	// then:
	msg := stringWriter.String()
	require.Contains(t, msg, "info message: 123")
}
