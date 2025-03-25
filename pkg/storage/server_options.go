package storage

// ServerOptions represents configurable options for the storage server
type ServerOptions struct {
	Port uint
}

func defaultServerOptions() ServerOptions {
	return ServerOptions{
		Port: 8080,
	}
}

// ServerOption is a function that modifies the server options
type ServerOption func(*ServerOptions)

// WithPort sets the port for the storage server
func WithPort(port uint) ServerOption {
	return func(o *ServerOptions) {
		o.Port = port
	}
}
