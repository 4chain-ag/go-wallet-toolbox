package whatsonchain

import "time"

// BSVExchangeRate is the rate struct for BSV exchange
type BSVExchangeRate struct {
	Timestamp time.Time
	Rate      float64
	Base      string
}

// BSVExchangeRateResponse is the response from WhatsOnChain for bsv exchange range
type BSVExchangeRateResponse struct {
	Time     int     `json:"time"`
	Rate     float64 `json:"rate"`
	Currency string  `json:"currency"`
}
