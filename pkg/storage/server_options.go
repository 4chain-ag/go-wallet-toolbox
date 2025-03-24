package storage

type ServerOptions struct {
	Port uint
}

func defaultServerOptions() ServerOptions {
	return ServerOptions{
		Port: 8080,
	}
}

type ServerOption func(*ServerOptions)

func WithPort(port uint) ServerOption {
	return func(o *ServerOptions) {
		o.Port = port
	}
}
