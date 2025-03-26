package wdk

import (
	"context"
	"fmt"

	"github.com/filecoin-project/go-jsonrpc"
)

// NewClient returns WalletStorageWriterClient and allows connection to rpc server
func NewClient(addr string, overrideOptions ...StorageClientOverrides) (*WalletStorageWriterClient, func(), error) {
	opts := defaultClientOptions()
	client := &WalletStorageWriterClient{
		client: &rpcWalletStorageWriter{},
	}

	for _, opt := range overrideOptions {
		opt(&opts)
	}

	cleanup, err := jsonrpc.NewMergeClient(
		context.Background(),
		addr,
		"remote_storage",
		[]any{client.client},
		nil,
		opts.options...,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize new RPC client: %w", err)
	}

	return client, cleanup, nil
}
