package wdk

import "time"

// Currency represents supported currency types
type Currency string

// Supported currency types
const (
	USD Currency = "USD"
	EUR Currency = "EUR"
	GBP Currency = "GBP"
)

// FiatExchangeRates is the rate struct for fiat currency
type FiatExchangeRates struct {
	Timestamp time.Time            `mapstructure:"timestamp"`
	Rates     map[Currency]float64 `mapstructure:"rates"`
	Base      Currency             `mapstructure:"base"`
}

// BSVExchangeRate is the rate struct for BSV exchange
type BSVExchangeRate struct {
	Timestamp time.Time `mapstructure:"timestamp"`
	Rate      float64   `mapstructure:"rates"`
	Base      Currency  `mapstructure:"base"`
}
