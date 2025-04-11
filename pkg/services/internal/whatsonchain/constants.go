package whatsonchain

import "time"

const (
	// DefaultBSVExchangeUpdateInterval is a duration after which the BSV Exchange Rate should be updated
	DefaultBSVExchangeUpdateInterval = 15 * time.Minute
	// DefaultFiatExchangeUpdateInterval is a duration after which the Fiat Exchange Rate should be updated
	DefaultFiatExchangeUpdateInterval = 24 * time.Hour
	// Retries is the number of retries the client should make when trying to query for the chain's resource
	Retries = 2
	// RetriesWaitTime is the duration to wait to retry a failed call to outside service
	RetriesWaitTime = 2 * time.Second
)
