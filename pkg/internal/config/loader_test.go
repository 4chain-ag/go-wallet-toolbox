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
	C SubConfig `mapstructure:"c_sub_config"`
}

type SubConfig struct {
	D string `mapstructure:"d_nested_field"`
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

const yamlConfig = `
b_with_long_name: 3
c_sub_config:
  d_nested_field: file_world
`

const dotEnvConfig = `
TEST_B_WITH_LONG_NAME=4
TEST_C_SUB_CONFIG_D_NESTED_FIELD="dotenv_world"
`

const dotEnvConfigEmptyPrefix = `
B_WITH_LONG_NAME=4
C_SUB_CONFIG_D_NESTED_FIELD="dotenv_world"
`

const jsonConfig = `
{
	"b_with_long_name": 5,
	"c_sub_config": {
		"d_nested_field": "json_world"
	}
}
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
	t.Setenv("TEST_C_SUB_CONFIG_D_NESTED_FIELD", "env_world")

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
	configFilePath := tempConfig(t, yamlConfig, "yaml")

	// when:
	err := loader.SetConfigFilePath(configFilePath)

	// then:
	require.NoError(t, err)

	// and:
	cfg, err := loader.Load()

	// then:
	require.NoError(t, err)
	require.Equal(t, "default_hello", cfg.A)
	require.Equal(t, 3, cfg.B)
	require.Equal(t, "file_world", cfg.C.D)
}

func TestDotEnvConfig(t *testing.T) {
	// given:
	loader := config.NewLoader(Defaults, "TEST")

	// and:
	t.Setenv("TEST_A", "env_hello")

	// and:
	configFilePath := tempConfig(t, dotEnvConfig, "env")

	// when:
	err := loader.SetConfigFilePath(configFilePath)

	// then:
	require.NoError(t, err)

	// and:
	cfg, err := loader.Load()

	// then:
	require.NoError(t, err)
	require.Equal(t, "env_hello", cfg.A)
	require.Equal(t, 4, cfg.B)
	require.Equal(t, "dotenv_world", cfg.C.D)
}

func TestJSONConfig(t *testing.T) {
	// given:
	loader := config.NewLoader(Defaults, "TEST")

	// and:
	t.Setenv("TEST_A", "env_hello")

	// and:
	configFilePath := tempConfig(t, jsonConfig, "json")

	// when:
	err := loader.SetConfigFilePath(configFilePath)

	// then:
	require.NoError(t, err)

	// and:
	cfg, err := loader.Load()

	// then:
	require.NoError(t, err)
	require.Equal(t, "env_hello", cfg.A)
	require.Equal(t, 5, cfg.B)
	require.Equal(t, "json_world", cfg.C.D)
}

func TestMixedConfig(t *testing.T) {
	// given:
	loader := config.NewLoader(Defaults, "TEST")

	// and:
	t.Setenv("TEST_B_WITH_LONG_NAME", "2")

	// and:
	configFilePath := tempConfig(t, yamlConfig, "yaml")

	// when:
	err := loader.SetConfigFilePath(configFilePath)

	// then:
	require.NoError(t, err)

	// and:
	cfg, err := loader.Load()

	// then:
	require.NoError(t, err)
	require.Equal(t, "default_hello", cfg.A)
	require.Equal(t, 2, cfg.B)
	require.Equal(t, "file_world", cfg.C.D)
}

func TestWithEmptyPrefix(t *testing.T) {
	// given:
	loader := config.NewLoader(Defaults, "")

	// and:
	t.Setenv("A", "env_hello")

	// and:
	configFilePath := tempConfig(t, dotEnvConfigEmptyPrefix, "env")

	// when:
	err := loader.SetConfigFilePath(configFilePath)

	// then:
	require.NoError(t, err)

	// and:
	cfg, err := loader.Load()

	// then:
	require.NoError(t, err)
	require.Equal(t, "env_hello", cfg.A)
	require.Equal(t, 4, cfg.B)
	require.Equal(t, "dotenv_world", cfg.C.D)
}

func TestEnvOverridesDotEnv(t *testing.T) {
	// given:
	loader := config.NewLoader(Defaults, "TEST")

	// and:
	t.Setenv("TEST_B_WITH_LONG_NAME", "2")
	t.Setenv("TEST_C_SUB_CONFIG_D_NESTED_FIELD", "env_world")

	// and:
	configFilePath := tempConfig(t, dotEnvConfig, "env")

	// when:
	err := loader.SetConfigFilePath(configFilePath)

	// then:
	require.NoError(t, err)

	// and:
	cfg, err := loader.Load()

	// then:
	require.NoError(t, err)
	require.Equal(t, "default_hello", cfg.A)
	require.Equal(t, 2, cfg.B)
	require.Equal(t, "env_world", cfg.C.D)
}

func tempConfig(t *testing.T, content, extension string) string {
	tmpDir := t.TempDir()
	configFilePath := fmt.Sprintf("%s/config.%s", tmpDir, extension)
	err := os.WriteFile(configFilePath, []byte(content), 0644)
	require.NoError(t, err)

	return configFilePath
}
