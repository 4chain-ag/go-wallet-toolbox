package infra

import (
	"fmt"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/defs"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/config"
)

// Config is the configuration for the "remote storage server" service (aka "infra")
type Config struct {
	ServerPrivateKey string          `mapstructure:"server_private_key"`
	BSVNetwork       defs.BSVNetwork `mapstructure:"bsv_network"`
	DBConfig         defs.Database        `mapstructure:"db"`
	HTTPConfig       HTTPConfig      `mapstructure:"http"`
	Logging          LogConfig       `mapstructure:"logging"`
}

// DBConfig is the configuration for the database
type DBConfig struct {
	Engine defs.DBType `mapstructure:"engine"`
}

// HTTPConfig is the configuration for the HTTP server related settings
type HTTPConfig struct {
	Port uint `mapstructure:"port"`
}

// LogConfig is the configuration for the logging
type LogConfig struct {
	Enabled bool            `mapstructure:"enabled"`
	Level   defs.LogLevel   `mapstructure:"level"`
	Handler defs.LogHandler `mapstructure:"handler"`
}

// Defaults returns the default configuration
func Defaults() Config {
	return Config{
		ServerPrivateKey: "", // it is not optional, user must provide it

		BSVNetwork: defs.NetworkMainnet,
		DBConfig:   defs.DefaultDBConfig(),
		HTTPConfig: HTTPConfig{
			Port: 8100,
		},
		Logging: LogConfig{
			Enabled: true,
			Level:   defs.LogLevelInfo,
			Handler: defs.JSONHandler,
		},
	}
}

// Validate validates the whole configuration
func (c *Config) Validate() (err error) {
	if c.ServerPrivateKey == "" {
		return fmt.Errorf("server private key is required")
	}

	if c.BSVNetwork, err = defs.ParseBSVNetworkStr(string(c.BSVNetwork)); err != nil {
		return fmt.Errorf("invalid BSV network: %w", err)
	}

	if err = c.DBConfig.Validate(); err != nil {
		return fmt.Errorf("invalid DB config: %w", err)
	}

	if err = c.Logging.Validate(); err != nil {
		return fmt.Errorf("invalid HTTP config: %w", err)
	}

	return nil
}

// Validate validates the DB configuration
func (c *DBConfig) Validate() (err error) {
	if c.Engine, err = defs.ParseDBTypeStr(string(c.Engine)); err != nil {
		return fmt.Errorf("invalid DB engine: %w", err)
	}

	return nil
}

// Validate validates the HTTP configuration
func (c *LogConfig) Validate() (err error) {
	if c.Level, err = defs.ParseLogLevelStr(string(c.Level)); err != nil {
		return fmt.Errorf("invalid log level: %w", err)
	}

	if c.Handler, err = defs.ParseHandlerTypeStr(string(c.Handler)); err != nil {
		return fmt.Errorf("invalid log handler: %w", err)
	}

	return nil
}

// ToYAMLFile writes the configuration to a YAML file
func (c *Config) ToYAMLFile(filename string) error {
	err := config.ToYAMLFile(c, filename)
	if err != nil {
		return fmt.Errorf("failed to write config to file: %w", err)
	}
	return nil
}
