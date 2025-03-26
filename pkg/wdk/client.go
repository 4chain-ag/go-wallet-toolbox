package wdk

import (
	"context"

	"github.com/filecoin-project/go-jsonrpc"
)

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

	return client, cleanup, err
}
