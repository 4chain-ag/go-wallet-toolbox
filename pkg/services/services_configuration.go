package services

import (
	"time"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/defs"
)

// WalletServicesConfiguration is a struct that has options for wallet services
type WalletServicesConfiguration struct {
	Chain                           defs.BSVNetwork        `mapstructure:"chain"`
	TaalApiKey                      string                 `mapstructure:"taal_api_key"`
	BitailsApiKey                   *string                `mapstructure:"bitails_api_key"`
	WhatsOnChainApiKey              string                 `mapstructure:"whats_on_chain_api_key"`
	BsvExchangeRate                 *defs.BSVExchangeRate  `mapstructure:"bsv_exchange_rate"`
	BsvUpdateInterval               time.Duration          `mapstructure:"bsv_update_interval"`
	FiatExchangeRates               defs.FiatExchangeRates `mapstructure:"fiat_exchange_rates"`
	FiatUpdateInterval              time.Duration          `mapstructure:"fiat_update_interval"`
	DisableMapiCallback             bool                   `mapstructure:"disable_mapi_callback"`
	ExchangeratesApiKey             string                 `mapstructure:"exchangerates_api_key"`
	ChaintracksFiatExchangeRatesUrl string                 `mapstructure:"chaintracks_fiat_exchange_rates_url"`
	Chaintracks                     any                    `mapstructure:"chaintracks"` // TODO: create *ChaintracksServiceClient
	ArcUrl                          string                 `mapstructure:"arc_url"`
	ArcConfig                       any                    `mapstructure:"arc"` // TODO: create *ArcConfig
}
