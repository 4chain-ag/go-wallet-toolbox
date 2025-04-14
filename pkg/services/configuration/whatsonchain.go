package configuration

import (
	"time"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
)

// WhatsOnChain is a struct that configures WhatsOnChain service
type WhatsOnChain struct {
	APIKey            string              `mapstructure:"api_key"`
	BSVExchangeRate   wdk.BSVExchangeRate `mapstructure:"bsv_exchange_rate"`
	BSVUpdateInterval *time.Duration      `mapstructure:"bsv_update_interval"`
}
