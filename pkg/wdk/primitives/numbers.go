package primitives

import "fmt"

// PositiveIntegerDefault10Max10000  is a positive integer that defaults to 10, and has an upper bound of 10000.
type PositiveIntegerDefault10Max10000 uint

// Validate checks if the integer is maximum 10000
func (i PositiveIntegerDefault10Max10000) Validate() error {
	if i > 10000 {
		return fmt.Errorf("is larger than 10000")
	}

	return nil
}

// PositiveInteger represents a positive integer value
type PositiveInteger uint64

// SatoshiValue Represents a value in Satoshis, constrained by the max supply of Bitcoin (2.1 * 10^15 Satoshis).
// @maximum 2100000000000000
type SatoshiValue uint64

// MaxSatoshis is the maximum number of Satoshis in the Bitcoin supply
const MaxSatoshis = 2100000000000000

// Validate checks if the value is less than the maximum number of Satoshis
func (s SatoshiValue) Validate() error {
	if s > MaxSatoshis {
		return fmt.Errorf("less than %d", MaxSatoshis)
	}
	return nil
}

// PositiveIntegerOrZero represents a positive integer or zero value
type PositiveIntegerOrZero uint

// IdentifierStringUnder300 is a string used for basket names, with a length under 300 bytes
type IdentifierStringUnder300 string

// Validate checks if the string is under 300 bytes long and not empty
func (b IdentifierStringUnder300) Validate() error {
	if len(b) > 300 {
		return fmt.Errorf("no more than 300 length")
	}
	if len(b) == 0 {
		return fmt.Errorf("at least 1 length")
	}
	return nil
}
