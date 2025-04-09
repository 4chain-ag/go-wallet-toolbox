package services

import (
	"fmt"
	"time"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/defs"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/services/internal/whatsonchain"
)

// Options is a function that can be used to override services options
type Options = func(*WalletServicesOptions)

// WalletServicesOptions is a struct that has options for wallet services
type WalletServicesOptions struct {
	TaalApiKey                      string
	BitailsApiKey                   *string
	WhatsOnChainApiKey              string
	BsvExchangeRate                 whatsonchain.BSVExchangeRate
	BsvUpdateMsecs                  int
	FiatExchangeRates               FiatExchangeRates
	FiatUpdateMsecs                 int
	DisableMapiCallback             bool
	ExchangeratesApiKey             string
	ChaintracksFiatExchangeRatesUrl string
	Chaintracks                     any // TODO: create *ChaintracksServiceClient
	ArcUrl                          string
	ArcConfig                       any // TODO: create *ArcConfig
}

func defaultServicesOptions(chain defs.BSVNetwork) *WalletServicesOptions {
	var taalApiKey string
	var port int
	var arcUrl string

	if chain == defs.NetworkMainnet {
		//nolint:gosec
		taalApiKey = "mainnet_9596de07e92300c6287e4393594ae39c"
		port = 8084
		arcUrl = "https://api.taal.com/arc"
	} else {
		//nolint:gosec
		taalApiKey = "testnet_0e6cf72133b43ea2d7861da2a38684e3"
		port = 8083
		arcUrl = "https://arc-test.taal.com/arc"
	}

	return &WalletServicesOptions{
		TaalApiKey: taalApiKey,
		BsvExchangeRate: whatsonchain.BSVExchangeRate{
			Timestamp: time.Date(2023, time.December, 13, 0, 0, 0, 0, time.UTC),
			Base:      "USD",
			Rate:      47.52,
		},
		BsvUpdateMsecs: 1000 * 60 * 15, // 15 minutes
		FiatExchangeRates: FiatExchangeRates{
			Timestamp: time.Date(2023, time.December, 13, 0, 0, 0, 0, time.UTC),
			Base:      "USD",
			Rates: map[string]float64{
				"USD": 1,
				"GBP": 0.8,
				"EUR": 0.93,
			},
		},
		FiatUpdateMsecs:                 1000 * 60 * 60 * 24, // 24 hours
		DisableMapiCallback:             true,                // rely on WalletMonitor by default
		ExchangeratesApiKey:             "bd539d2ff492bcb5619d5f27726a766f",
		ChaintracksFiatExchangeRatesUrl: fmt.Sprintf("https://npm-registry.babbage.systems:%d/getFiatExchangeRates", port),
		Chaintracks:                     nil, // TODO: implement me
		ArcUrl:                          arcUrl,
		ArcConfig:                       nil, // TODO: implement me
	}
}

// WithTaalApiKey sets the Taal API key.
func WithTaalApiKey(apiKey string) Options {
	return func(o *WalletServicesOptions) {
		o.TaalApiKey = apiKey
	}
}

// WithBitailsApiKey sets the Bitails API key.
func WithBitailsApiKey(apiKey *string) Options {
	return func(o *WalletServicesOptions) {
		o.BitailsApiKey = apiKey
	}
}

// WithWhatsOnChainApiKey sets the WhatsOnChain API key.
func WithWhatsOnChainApiKey(apiKey string) Options {
	return func(o *WalletServicesOptions) {
		o.WhatsOnChainApiKey = apiKey
	}
}

// WithBsvExchangeRate sets the BSV exchange rate.
func WithBsvExchangeRate(exchangeRate whatsonchain.BSVExchangeRate) Options {
	return func(o *WalletServicesOptions) {
		o.BsvExchangeRate = exchangeRate
	}
}

// WithBsvUpdateMsecs sets the update interval for BSV exchange rates in milliseconds.
func WithBsvUpdateMsecs(updateMsecs int) Options {
	return func(o *WalletServicesOptions) {
		o.BsvUpdateMsecs = updateMsecs
	}
}

// WithFiatExchangeRates sets the fiat exchange rates.
func WithFiatExchangeRates(fiatRates FiatExchangeRates) Options {
	return func(o *WalletServicesOptions) {
		o.FiatExchangeRates = fiatRates
	}
}

// WithFiatUpdateMsecs sets the update interval for fiat exchange rates in milliseconds.
func WithFiatUpdateMsecs(updateMsecs int) Options {
	return func(o *WalletServicesOptions) {
		o.FiatUpdateMsecs = updateMsecs
	}
}

// WithDisableMapiCallback disables or enables MAPI callbacks.
func WithDisableMapiCallback(disable bool) Options {
	return func(o *WalletServicesOptions) {
		o.DisableMapiCallback = disable
	}
}

// WithExchangeratesApiKey sets the ExchangeRates API key.
func WithExchangeratesApiKey(apiKey string) Options {
	return func(o *WalletServicesOptions) {
		o.ExchangeratesApiKey = apiKey
	}
}

// WithChaintracksFiatExchangeRatesUrl sets the Chaintracks fiat exchange rates URL.
func WithChaintracksFiatExchangeRatesUrl(url string) Options {
	return func(o *WalletServicesOptions) {
		o.ChaintracksFiatExchangeRatesUrl = url
	}
}

// WithChaintracks sets the Chaintracks service client.
func WithChaintracks(chaintracks any) Options {
	return func(o *WalletServicesOptions) {
		o.Chaintracks = chaintracks
	}
}

// WithArcUrl sets the ARC URL.
func WithArcUrl(url string) Options {
	return func(o *WalletServicesOptions) {
		o.ArcUrl = url
	}
}

// WithArcConfig sets the ARC configuration.
func WithArcConfig(config any) Options {
	return func(o *WalletServicesOptions) {
		o.ArcConfig = config
	}
}
