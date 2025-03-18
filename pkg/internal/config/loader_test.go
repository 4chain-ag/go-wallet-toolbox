package config_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/config"
	"github.com/stretchr/testify/require"
)

type MockConfig struct {
	A string    `mapstructure:"a"`
	B int       `mapstructure:"b_with_long_name"`
	C SubConfig `mapstructure:"c"`
}

func Defaults() MockConfig {
	return MockConfig{
		A: "default_hello",
		B: 1,
		C: SubConfig{
			D: "default_world",
		},
	}
}

type SubConfig struct {
	D string `mapstructure:"d"`
}

const yamlConfig = `
# field "a" is skipped

b_with_long_name: 3
c:
  d: file_world
`

func TestDefaults(t *testing.T) {
	// given:
	loader := config.NewLoader(Defaults, "TEST")

	// when:
	cfg, err := loader.Load()

	// then:
	require.NoError(t, err)
	require.Equal(t, "default_hello", cfg.A)
	require.Equal(t, 1, cfg.B)
	require.Equal(t, "default_world", cfg.C.D)
}

func TestEnvVariables(t *testing.T) {
	// given:
	loader := config.NewLoader(Defaults, "TEST")

	// and:
	t.Setenv("TEST_B_WITH_LONG_NAME", "2")
	t.Setenv("TEST_C_D", "env_world")

	// when:
	cfg, err := loader.Load()

	// then:
	require.NoError(t, err)
	require.Equal(t, "default_hello", cfg.A)
	require.Equal(t, 2, cfg.B)
	require.Equal(t, "env_world", cfg.C.D)
}

func TestFileConfig(t *testing.T) {
	// given:
	loader := config.NewLoader(Defaults, "TEST")

	// and:
	configFilePath := tempConfig(t, yamlConfig)

	// when:
	loader.SetConfigFilePath(configFilePath)

	// and:
	cfg, err := loader.Load()

	// then:
	require.NoError(t, err)
	require.Equal(t, "default_hello", cfg.A)
	require.Equal(t, 3, cfg.B)
	require.Equal(t, "file_world", cfg.C.D)
}

func TestMixedConfig(t *testing.T) {
	// given:
	loader := config.NewLoader(Defaults, "TEST")

	// and:
	t.Setenv("TEST_B_WITH_LONG_NAME", "2")

	// and:
	configFilePath := tempConfig(t, yamlConfig)

	// when:
	loader.SetConfigFilePath(configFilePath)

	// and:
	cfg, err := loader.Load()

	// then:
	require.NoError(t, err)
	require.Equal(t, "default_hello", cfg.A)
	require.Equal(t, 2, cfg.B)
	require.Equal(t, "file_world", cfg.C.D)
}

func tempConfig(t *testing.T, content string) string {
	tmpDir := t.TempDir()
	configFilePath := fmt.Sprintf("%s/config.yaml", tmpDir)
	err := os.WriteFile(configFilePath, []byte(content), 0644)
	require.NoError(t, err)

	return configFilePath
}
