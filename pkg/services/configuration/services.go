package configuration

import (
	"time"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/defs"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
)

// WalletServices is a struct that has options for wallet services
type WalletServices struct {
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

	WhatsOnChain WhatsOnChain
}
