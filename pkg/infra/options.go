package infra

// Options is the parameters for initializing the "infra" server
type Options struct {
	EnvPrefix  string
	ConfigFile string
}

// DefaultOptions returns the default parameters to initialize the "infra" server
func DefaultOptions() Options {
	return Options{
		EnvPrefix:  "INFRA",
		ConfigFile: "",
	}
}

// InitOption is a function that sets a parameter for initializing the "infra" server
type InitOption func(*Options)

// WithEnvPrefix sets the environment variable prefix for the "infra" server, all environment variables will be prefixed with this:
// e.g. "INFRA_HTTP_PORT=8100"
func WithEnvPrefix(prefix string) InitOption {
	return func(o *Options) {
		o.EnvPrefix = prefix
	}
}

// WithConfigFile sets the configuration file for the "infra" server, the configuration file is in YAML format
func WithConfigFile(file string) InitOption {
	return func(o *Options) {
		o.ConfigFile = file
	}
}
