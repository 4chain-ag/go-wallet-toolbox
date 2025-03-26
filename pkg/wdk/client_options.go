package wdk

import (
	"net/http"
)

// StorageClientOption is a function that can be used to override internal dependencies.
type StorageClientOption = func(*ClientOptions)

// ClientOptions represents configurable options for the client
type ClientOptions struct {
	httpClient *http.Client
}

func defaultClientOptions() ClientOptions {
	return ClientOptions{
		httpClient: nil,
	}
}

// WithHttpClient is a function that can be used to override the http.Client used by the client.
func WithHttpClient(httpClient *http.Client) StorageClientOption {
	return func(o *ClientOptions) {
		o.httpClient = httpClient
	}
}
