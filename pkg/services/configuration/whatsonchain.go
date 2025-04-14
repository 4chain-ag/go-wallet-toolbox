package configuration

import (
	"time"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
)

// WhatsOnChain is a struct that configures WhatsOnChain service
type WhatsOnChain struct {
	ApiKey            string
	BsvExchangeRate   wdk.BSVExchangeRate
	BsvUpdateInterval *time.Duration
}
