package whatsonchain

import (
	"time"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/services/internal"
)

// BSVExchangeRate is the rate struct for BSV exchange
type BSVExchangeRate struct {
	Timestamp time.Time
	Rate      float64
	Base      internal.Currency
}

// BSVExchangeRateResponse is the response from WhatsOnChain for bsv exchange range
type BSVExchangeRateResponse struct {
	Time     int     `json:"time"`
	Rate     float64 `json:"rate"`
	Currency string  `json:"currency"`
}
