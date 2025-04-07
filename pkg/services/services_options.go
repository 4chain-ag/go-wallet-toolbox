package services

import (
	"fmt"
	"time"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/defs"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
)

// ServicesOptions is a function that can be used to override services options
type Options = func(*WalletServicesOptions)

type WalletServicesOptions struct {
	TaalApiKey                      string
	BitailsApiKey                   *string
	WhatsOnChainApiKey              *string
	BsvExchangeRate                 wdk.BSVExchangeRate
	BsvUpdateMsecs                  int
	FiatExchangeRates               wdk.FiatExchangeRates
	FiatUpdateMsecs                 int
	DisableMapiCallback             bool
	ExchangeratesApiKey             string
	ChaintracksFiatExchangeRatesUrl string
	Chaintracks                     any // TODO: create *ChaintracksServiceClient
	ArcUrl                          string
	ArcConfig                       any // TODO: create *ArcConfig
}

func defaultServicesOptions(chain defs.BSVNetwork) WalletServicesOptions {
	var taalApiKey string
	var port int
	var arcUrl string

	if chain == defs.NetworkMainnet {
		taalApiKey = "mainnet_9596de07e92300c6287e4393594ae39c"
		port = 8084
		arcUrl = "https://api.taal.com/arc"
	} else {
		taalApiKey = "testnet_0e6cf72133b43ea2d7861da2a38684e3"
		port = 8083
		arcUrl = "https://arc-test.taal.com/arc"
	}

	return WalletServicesOptions{
		TaalApiKey: taalApiKey,
		BsvExchangeRate: wdk.BSVExchangeRate{
			Timestamp: time.Date(2023, time.December, 13, 0, 0, 0, 0, time.UTC),
			Base:      "USD",
			Rate:      47.52,
		},
		BsvUpdateMsecs: 1000 * 60 * 15, // 15 minutes
		FiatExchangeRates: wdk.FiatExchangeRates{
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
