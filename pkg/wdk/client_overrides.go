package wdk

import (
	"net/http"

	"github.com/filecoin-project/go-jsonrpc"
)

// InternalOverrides is a function that can be used to override internal dependencies.
// This is meant to be used for testing purposes.
type InternalOverrides = func(*overrides)

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
func WithHttpClient(httpClient *http.Client) InternalOverrides {
	return func(o *overrides) {
		o.options = append(o.options, jsonrpc.WithHTTPClient(httpClient))
	}
}
