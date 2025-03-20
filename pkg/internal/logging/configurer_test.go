package logging_test

import (
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/logging"
	"testing"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/defs"
	"github.com/stretchr/testify/require"
)

func TestTextLogger(t *testing.T) {
	// given:
	stringWriter := &logging.TestWriter{}
	logger := logging.New().
		WithLevel(defs.LogLevelDebug).
		WithHandler(defs.TextHandler, stringWriter).
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
		WithLevel(defs.LogLevelDebug).
		WithHandler(defs.JSONHandler, stringWriter).
		Logger()

	// when:
	logger.Debug("debug message")

	// then:
	msg := stringWriter.String()
	require.Contains(t, msg, `"time"`)
	require.Contains(t, msg, `"level":"DEBUG"`)
	require.Contains(t, msg, `"msg":"debug message"`)
}
