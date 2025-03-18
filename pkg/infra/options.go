package infra

type InitParams struct {
	EnvPrefix  string
	ConfigFile string
}

func DefaultParams() InitParams {
	return InitParams{
		EnvPrefix:  "INFRA",
		ConfigFile: "",
	}
}

type InitOption func(*InitParams)

func WithEnvPrefix(prefix string) InitOption {
	return func(o *InitParams) {
		o.EnvPrefix = prefix
	}
}

func WithConfigFile(file string) InitOption {
	return func(o *InitParams) {
		o.ConfigFile = file
	}
}
