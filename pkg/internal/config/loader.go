package config

import (
	"fmt"
	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
	"os"
	"strings"
)

const (
	DefaultConfigFilePath = "config.yaml"
)

type Loader[T any] struct {
	cfg            T
	envPrefix      string
	configFilePath string
	viper          *viper.Viper
}

func NewLoader[T any](defaults func() T, envPrefix string) *Loader[T] {
	return &Loader[T]{
		cfg:            defaults(),
		envPrefix:      envPrefix,
		configFilePath: DefaultConfigFilePath,
		viper:          viper.New(),
	}
}

func (l *Loader[T]) SetConfigFilePath(path string) {
	l.configFilePath = path
}

// Load loads the configuration from the environment and the config file.
// NOTE: The priority of the values is as follows:
// 1. Environment variables
// 2. Config file
// 3. Default values
//
// NOTE: The config file is optional.
// NOTE: For multilevel nested structs, the keys in the ENV variables should be separated by underscores.
// e.g. for the nesting:
// a:
//
//	b_with_long_name:
//	  c: value
//
// the ENV variable should be named as: <ENVPREFIX>_A_B_WITH_LONG_NAME_C
// the ENVPREFIX is the prefix that is passed to the NewLoader function.
func (l *Loader[T]) Load() (T, error) {
	if err := l.setViperDefaults(); err != nil {
		return l.cfg, err
	}

	l.loadFromEnv()

	if err := l.loadFromFile(); err != nil {
		return l.cfg, err
	}

	if err := l.viperToCfg(); err != nil {
		return l.cfg, err
	}

	return l.cfg, nil
}

func (l *Loader[T]) setViperDefaults() error {
	defaultsMap := make(map[string]any)
	if err := mapstructure.Decode(l.cfg, &defaultsMap); err != nil {
		err = fmt.Errorf("error occurred while setting defaults: %w", err)
		return err
	}

	for k, v := range defaultsMap {
		l.viper.SetDefault(k, v)
	}

	return nil
}

func (l *Loader[T]) loadFromEnv() {
	l.viper.SetEnvPrefix(l.envPrefix)
	l.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	l.viper.AutomaticEnv()
}

func (l *Loader[T]) loadFromFile() error {
	if l.configFilePath == DefaultConfigFilePath {
		_, err := os.Stat(l.configFilePath)
		if os.IsNotExist(err) {
			// Config file not specified. Using defaults
			return nil
		}
	}

	l.viper.SetConfigFile(l.configFilePath)
	if err := l.viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error while reading config fiel: %w", err)
	}

	return nil
}

func (l *Loader[T]) viperToCfg() error {
	if err := l.viper.Unmarshal(&l.cfg); err != nil {
		return fmt.Errorf("error while unmarshalling config from viper: %w", err)
	}
	return nil
}
