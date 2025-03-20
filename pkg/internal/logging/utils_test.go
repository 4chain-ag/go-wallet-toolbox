package logging_test

import (
	"fmt"
	"testing"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/defs"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/logging"
	"github.com/stretchr/testify/require"
)

func TestChildLogger(t *testing.T) {
	// given:
	stringWriter := &logging.TestWriter{}
	logger := logging.New().
		WithLevel(defs.LogLevelDebug).
		WithHandler(defs.TextHandler, stringWriter).
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

func TestErrorf(t *testing.T) {
	// given:
	stringWriter := &logging.TestWriter{}
	logger := logging.New().
		WithLevel(defs.LogLevelError).
		WithHandler(defs.TextHandler, stringWriter).
		Logger()

	// and:
	err := fmt.Errorf("error message")

	// when:
	logging.Errorf(logger, err, "additional context %d", 123)

	// then:
	msg := stringWriter.String()
	require.Contains(t, msg, `msg="additional context 123"`)
	require.Contains(t, msg, `error="error message"`)
}

func TestNopIfNil(t *testing.T) {
	// when:
	logger := logging.DefaultIfNil(nil)

	// then:
	require.NotNil(t, logger)
}
