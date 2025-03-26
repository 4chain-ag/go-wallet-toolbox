package wdk

import (
	"net/http"

	"github.com/filecoin-project/go-jsonrpc"
)

// StorageClientOverrides is a function that can be used to override internal dependencies.
// This is meant to be used for testing purposes.
type StorageClientOverrides = func(*overrides)

type overrides struct {
	options []jsonrpc.Option
}

func defaultClientOptions() overrides {
	return overrides{
		options: []jsonrpc.Option{
			jsonrpc.WithMethodNamer(jsonrpc.NoNamespaceDecapitalizedMethodNamer),
		},
	}
}

// WithHttpClient is a function that can be used to override the http.Client used by the client.
// This is meant to be used for testing purposes.
func WithHttpClient(httpClient *http.Client) StorageClientOverrides {
	return func(o *overrides) {
		o.options = append(o.options, jsonrpc.WithHTTPClient(httpClient))
	}
}
