package infra

import (
	"fmt"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/defs"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/config"
)

// Config is the configuration for the "remote storage server" service (aka "infra")
type Config struct {
	BSVNetwork defs.BSVNetwork `mapstructure:"bsv_network"`
	DBConfig   DBConfig        `mapstructure:"db"`
	HTTPConfig HTTPConfig      `mapstructure:"http"`
}

// DBConfig is the configuration for the database
type DBConfig struct {
	Engine defs.DBType `mapstructure:"engine"`
}

// HTTPConfig is the configuration for the HTTP server related settings
type HTTPConfig struct {
	Port uint `mapstructure:"port"`
}

// Defaults returns the default configuration
func Defaults() Config {
	return Config{
		BSVNetwork: defs.NetworkMainnet,
		DBConfig: DBConfig{
			Engine: defs.DBTypeSQLite,
		},
		HTTPConfig: HTTPConfig{
			Port: 8100,
		},
	}
}

// Validate validates the whole configuration
func (c Config) Validate() error {
	if _, err := defs.ParseBSVNetworkStr(string(c.BSVNetwork)); err != nil {
		return fmt.Errorf("invalid BSV network: %w", err)
	}

	if err := c.DBConfig.Validate(); err != nil {
		return fmt.Errorf("invalid DB config: %w", err)
	}

	return nil
}

// Validate validates the DB configuration
func (c DBConfig) Validate() error {
	if _, err := defs.ParseDBTypeStr(string(c.Engine)); err != nil {
		return fmt.Errorf("invalid DB engine: %w", err)
	}

	return nil
}

// ToYAMLFile writes the configuration to a YAML file
func (c Config) ToYAMLFile(filename string) error {
	err := config.ToYAMLFile(c, filename)
	if err != nil {
		return fmt.Errorf("failed to write config to file: %w", err)
	}
	return nil
}
