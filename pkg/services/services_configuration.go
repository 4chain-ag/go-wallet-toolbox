package services

import (
	"time"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/defs"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/services/internal/whatsonchain"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
)

// WalletServicesConfiguration is a struct that has options for wallet services
type WalletServicesConfiguration struct {
	Chain                           defs.BSVNetwork
	TaalApiKey                      string
	BitailsApiKey                   *string
	FiatExchangeRates               wdk.FiatExchangeRates
	FiatUpdateInterval              *time.Duration
	DisableMapiCallback             bool
	ExchangeratesApiKey             string
	ChaintracksFiatExchangeRatesUrl string
	Chaintracks                     any // TODO: create *ChaintracksServiceClient
	ArcUrl                          string
	ArcConfig                       any // TODO: create *ArcConfig

	WhatsOnChainConfiguration whatsonchain.WhatsOnChainConfiguration
}
