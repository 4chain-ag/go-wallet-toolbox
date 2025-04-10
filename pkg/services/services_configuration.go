package services

import (
	"time"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/defs"
)

// WalletServicesConfiguration is a struct that has options for wallet services
type WalletServicesConfiguration struct {
	Chain                           defs.BSVNetwork
	TaalApiKey                      string
	BitailsApiKey                   *string
	WhatsOnChainApiKey              string
	BsvExchangeRate                 *defs.BSVExchangeRate
	BsvUpdateInterval               time.Duration
	FiatExchangeRates               defs.FiatExchangeRates
	FiatUpdateInterval              time.Duration
	DisableMapiCallback             bool
	ExchangeratesApiKey             string
	ChaintracksFiatExchangeRatesUrl string
	Chaintracks                     any // TODO: create *ChaintracksServiceClient
	ArcUrl                          string
	ArcConfig                       any // TODO: create *ArcConfig
}
