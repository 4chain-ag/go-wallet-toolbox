package infra_test

import (
	"testing"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/defs"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/infra"
	"github.com/stretchr/testify/require"
)

func TestCaseInsensitiveEnums(t *testing.T) {
	// given:
	t.Setenv("TEST_DB_ENGINE", "SQLite")
	t.Setenv("TEST_BSV_NETWORK", "MAIN")
	t.Setenv("TEST_LOGGING_LEVEL", "DeBug")
	t.Setenv("TEST_LOGGING_HANDLER", "Text")

	// when:
	infraSrv, err := infra.NewServer(infra.WithEnvPrefix("TEST"))

	// then:
	require.NoError(t, err)
	require.Equal(t, defs.DBTypeSQLite, infraSrv.Config.DBConfig.Engine)
	require.Equal(t, defs.NetworkMainnet, infraSrv.Config.BSVNetwork)
	require.Equal(t, defs.LogLevelDebug, infraSrv.Config.Logging.Level)
	require.Equal(t, defs.TextHandler, infraSrv.Config.Logging.Handler)
}

func TestEnums(t *testing.T) {
	tests := map[string]struct {
		envKey string
	}{
		"DB engine": {
			envKey: "TEST_DB_ENGINE",
		},
		"BSV network": {
			envKey: "TEST_BSV_NETWORK",
		},
		"HTTP port": {
			envKey: "TEST_HTTP_PORT",
		},
		"Logging level": {
			envKey: "TEST_LOGGING_LEVEL",
		},
		"Logging handler": {
			envKey: "TEST_LOGGING_HANDLER",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			t.Setenv(test.envKey, "wrong")

			// when:
			_, err := infra.NewServer(infra.WithEnvPrefix("TEST"))

			// then:
			require.Error(t, err)
		})
	}
}
