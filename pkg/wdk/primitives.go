package wdk

import (
	"fmt"
	"regexp"
)

// String5to2000Bytes represents a string used for descriptions,
// with a length between 5 and 2000 characters.
type String5to2000Bytes string

func (d String5to2000Bytes) Validate() error {
	if len(d) < 5 {
		return fmt.Errorf("at least 5 length")
	}
	if len(d) > 2000 {
		return fmt.Errorf("no more than 2000 length")
	}
	return nil
}

// Base64String is a string in base64 format
type Base64String string

// HexString is a string in hexadecimal format
type HexString string

var hexRegex = regexp.MustCompile("^[0-9a-fA-F]+$")

func (h HexString) Validate() error {
	if len(h)%2 != 0 {
		return fmt.Errorf("even length, not %d", len(h))
	}

	if !hexRegex.MatchString(string(h)) {
		return fmt.Errorf("hexadecimal string")
	}
	return nil
}

// BooleanDefaultTrue is a boolean with a default value of true
type BooleanDefaultTrue bool

// BooleanDefaultFalse is a boolean with a default value of false
type BooleanDefaultFalse bool

// PositiveInteger represents a positive integer value
type PositiveInteger uint

// SatoshiValue Represents a value in Satoshis, constrained by the max supply of Bitcoin (2.1 * 10^15 Satoshis).
// @minimum 1
// @maximum 2100000000000000
type SatoshiValue uint

const MaxSatoshis = 2100000000000000

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

func (b IdentifierStringUnder300) Validate() error {
	if len(b) > 300 {
		return fmt.Errorf("no more than 300 length")
	}
	if len(b) == 0 {
		return fmt.Errorf("at least 1 length")
	}
	return nil
}

// TXIDHexString is a hexadecimal transaction ID
type TXIDHexString string

// BEEF An array of integers, each ranging from 0 to 255, indicating transaction data in BEEF (BRC-62) format.
type BEEF []byte

// OutPoint identifies a unique transaction output by its txid and index vout
type OutPoint struct {
	// TxID Transaction double sha256 hash as big endian hex string
	TxID string
	// Vout Zero based output index within the transaction
	Vout int
}

// AuthID represents the identity of the user making the request
type AuthID struct {
	IdentityKey string
	UserID      *int
	IsActive    *bool
}
