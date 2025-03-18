package infra

import (
	"fmt"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/defs"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/config"
)

type Config struct {
	BSVNetwork defs.BSVNetwork `mapstructure:"bsv_network"`
	DBConfig   DBConfig        `mapstructure:"db"`
	HTTPConfig HTTPConfig      `mapstructure:"http"`
}

type DBConfig struct {
	Engine defs.DBType `mapstructure:"engine"`
}

type HTTPConfig struct {
	Port uint `mapstructure:"port"`
}

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

func (c Config) Validate() error {
	if _, err := defs.ParseBSVNetworkStr(string(c.BSVNetwork)); err != nil {
		return err
	}

	if err := c.DBConfig.Validate(); err != nil {
		return err
	}

	return nil
}

func (c DBConfig) Validate() error {
	if _, err := defs.ParseDBTypeStr(string(c.Engine)); err != nil {
		return err
	}

	return nil
}

func (c Config) ToYAMLFile(filename string) error {
	err := config.ToYAMLFile(c, filename)
	if err != nil {
		return fmt.Errorf("failed to write config to file: %w", err)
	}
	return nil
}
