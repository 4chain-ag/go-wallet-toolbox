package defs

import (
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
