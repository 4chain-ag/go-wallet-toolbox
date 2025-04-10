package defs

import (
	"fmt"
	"time"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk/primitives"
)

// FiatExchangeRates is the rate struct for fiat currency
type FiatExchangeRates struct {
	Timestamp time.Time                       `mapstructure:"timestamp"`
	Rates     map[primitives.Currency]float64 `mapstructure:"rates"`
	Base      primitives.Currency             `mapstructure:"base"`
}

// BSVExchangeRate is the rate struct for BSV exchange
type BSVExchangeRate struct {
	Timestamp time.Time           `mapstructure:"timestamp"`
	Rate      float64             `mapstructure:"rates"`
	Base      primitives.Currency `mapstructure:"base"`
}

// BSVExchangeRateResponse is the response from WhatsOnChain for bsv exchange range
type BSVExchangeRateResponse struct {
	Time     int     `json:"time" mapstructure:"time"`
	Rate     float64 `json:"rate" mapstructure:"rate"`
	Currency string  `json:"currency" mapstructure:"currency"`
}

type ServicesConfig struct {
	Chain                           BSVNetwork        `mapstructure:"chain"`
	TaalApiKey                      string            `mapstructure:"taal_api_key"`
	BitailsApiKey                   *string           `mapstructure:"bitails_api_key"`
	WhatsOnChainApiKey              string            `mapstructure:"whats_on_chain_api_key"`
	BsvExchangeRate                 *BSVExchangeRate  `mapstructure:"bsv_exchange_rate"`
	BsvUpdateInterval               time.Duration     `mapstructure:"bsv_update_interval"`
	FiatExchangeRates               FiatExchangeRates `mapstructure:"fiat_exchange_rates"`
	FiatUpdateInterval              time.Duration     `mapstructure:"fiat_update_interval"`
	DisableMapiCallback             bool              `mapstructure:"disable_mapi_callback"`
	ExchangeratesApiKey             string            `mapstructure:"exchangerates_api_key"`
	ChaintracksFiatExchangeRatesUrl string            `mapstructure:"chaintracks_fiat_exchange_rates_url"`
	Chaintracks                     any               `mapstructure:"chaintracks"` // TODO: create *ChaintracksServiceClient
	ArcUrl                          string            `mapstructure:"arc_url"`
	ArcConfig                       any               `mapstructure:"arc"` // TODO: create *ArcConfig
}

// DefaultServicesConfig returns default services configuration
func DefaultServicesConfig(chain BSVNetwork) ServicesConfig {
	var taalApiKey string
	var port int
	var arcUrl string

	if chain == NetworkMainnet {
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

	return ServicesConfig{
		TaalApiKey:        taalApiKey,
		BsvExchangeRate:   nil,
		BsvUpdateInterval: FifteenMinutes,
		FiatExchangeRates: FiatExchangeRates{
			Timestamp: time.Date(2023, time.December, 13, 0, 0, 0, 0, time.UTC),
			Base:      primitives.USD,
			Rates: map[primitives.Currency]float64{
				primitives.USD: 1,
				primitives.GBP: 0.8,
				primitives.EUR: 0.93,
			},
		},
		FiatUpdateInterval:              TwentyFourHours,
		DisableMapiCallback:             true, // rely on WalletMonitor by default
		ExchangeratesApiKey:             "bd539d2ff492bcb5619d5f27726a766f",
		ChaintracksFiatExchangeRatesUrl: fmt.Sprintf("https://npm-registry.babbage.systems:%d/getFiatExchangeRates", port),
		Chaintracks:                     nil, // TODO: implement me
		ArcUrl:                          arcUrl,
		ArcConfig:                       nil, // TODO: implement me
	}
}
