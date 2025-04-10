package defs

import "time"

const (
	// FifteenMinutes is a time.Duration of 15 minutes in milliseconds
	FifteenMinutes = 15 * time.Minute
	// TwentyFourHours is a time.Duration of 24 hours in milliseconds
	TwentyFourHours = 24 * time.Hour
	// Retries is the number of retries the client should make when trying to query for the chain's resource
	Retries = 2
	// RetriesWaitTime is the duration to wait to retry a failed call to outside service
	RetriesWaitTime = 2 * time.Second
)
