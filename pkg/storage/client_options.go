package storage

import (
	"net/http"

	"github.com/filecoin-project/go-jsonrpc"
)

// ClientOptions is a function that can be used to override internal dependencies.
// This is meant to be used for testing purposes.
type ClientOptions = func(*clientOptions)

type clientOptions struct {
	options []jsonrpc.Option
}

func defaultClientOptions() clientOptions {
	return clientOptions{
		options: []jsonrpc.Option{
			jsonrpc.WithMethodNameFormatter(jsonrpc.NewMethodNameFormatter(false, jsonrpc.LowerFirstCharCase)),
		},
	}
}

// WithHttpClient is a function that can be used to override the http.Client used by the client.
// This is meant to be used for testing purposes.
func WithHttpClient(httpClient *http.Client) ClientOptions {
	return func(o *clientOptions) {
		o.options = append(o.options, jsonrpc.WithHTTPClient(httpClient))
	}
}
