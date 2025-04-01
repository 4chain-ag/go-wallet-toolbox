package storage

import (
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/actions"
	"gorm.io/gorm"
)

// ProviderOption is function for additional setup of Provider itself.
type ProviderOption func(*providerOptions)

type providerOptions struct {
	gormDB *gorm.DB
	funder actions.Funder
}

// WithGORM sets the GORM database for the provider.
func WithGORM(gormDB *gorm.DB) ProviderOption {
	return func(o *providerOptions) {
		o.gormDB = gormDB
	}
}

// WithFunder sets the funder for the provider.
func WithFunder(funder actions.Funder) ProviderOption {
	return func(o *providerOptions) {
		o.funder = funder
	}
}

func toOptions(opts []ProviderOption) *providerOptions {
	options := &providerOptions{}
	for _, opt := range opts {
		opt(options)
	}
	return options
}
