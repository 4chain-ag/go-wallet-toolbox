package whatsonchain

import (
	"time"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
)

// WhatsOnChainConfiguration is a struct that configures WhatsOnChain service
type WhatsOnChainConfiguration struct {
	ApiKey            string
	BsvExchangeRate   *wdk.BSVExchangeRate
	BsvUpdateInterval *time.Duration
}
