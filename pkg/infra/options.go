package infra

// InitParams is the parameters for initializing the "infra" server
type InitParams struct {
	EnvPrefix  string
	ConfigFile string
}

// DefaultParams returns the default parameters to initialize the "infra" server
func DefaultParams() InitParams {
	return InitParams{
		EnvPrefix:  "INFRA",
		ConfigFile: "",
	}
}

// InitOption is a function that sets a parameter for initializing the "infra" server
type InitOption func(*InitParams)

// WithEnvPrefix sets the environment variable prefix for the "infra" server, all environment variables will be prefixed with this:
// e.g. "INFRA_HTTP_PORT=8100"
func WithEnvPrefix(prefix string) InitOption {
	return func(o *InitParams) {
		o.EnvPrefix = prefix
	}
}

// WithConfigFile sets the configuration file for the "infra" server, the configuration file is in YAML format
func WithConfigFile(file string) InitOption {
	return func(o *InitParams) {
		o.ConfigFile = file
	}
}
