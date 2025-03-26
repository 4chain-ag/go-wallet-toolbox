package wdk

import (
	"context"
	"fmt"

	"github.com/filecoin-project/go-jsonrpc"
)

// NewClient returns WalletStorageWriterClient that allows connection to rpc server
func NewClient(addr string, overrideOptions ...StorageClientOption) (*WalletStorageWriterClient, func(), error) {
	opts := defaultClientOptions()
	for _, opt := range overrideOptions {
		opt(&opts)
	}

	client := &WalletStorageWriterClient{
		client: &rpcWalletStorageWriter{},
	}

	rpcClientOptions := []jsonrpc.Option{
		jsonrpc.WithMethodNamer(jsonrpc.NoNamespaceDecapitalizedMethodNamer),
	}
	if opts.httpClient != nil {
		rpcClientOptions = append(rpcClientOptions, jsonrpc.WithHTTPClient(opts.httpClient))
	}

	cleanup, err := jsonrpc.NewMergeClient(
		context.Background(),
		addr,
		"remote_storage",
		[]any{client.client},
		nil,
		rpcClientOptions...,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize new RPC client: %w", err)
	}

	return client, cleanup, nil
}
