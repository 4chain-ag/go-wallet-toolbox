package wdk

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"

	"github.com/go-softwarelab/common/pkg/to"
)

const (
	// pubKeyBytesLenCompressed is the length of the compressed pub key
	pubKeyBytesLenCompressed = 33
	// pubKeyBytesLenUncompressed is the length of the uncompressed pub key
	pubKeyBytesLenUncompressed = 65
)

// String5to2000Bytes represents a string used for descriptions,
// with a length between 5 and 2000 characters.
type String5to2000Bytes string

// Validate checks if the string is between 5 and 2000 characters long
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

// Validate will check if string is proper based64 encoded string
func (s Base64String) Validate() error {
	// Step 1: Check if the string's length is divisible by 4 (Base64 requirement)
	if len(s)%4 != 0 {
		return fmt.Errorf("base64 string length must be divisible by 4")
	}

	// Step 2: Validate padding
	if strings.HasSuffix(string(s), "===") {
		return fmt.Errorf("invalid base64 padding")
	}

	// Step 3: Check if the string is valid Base64
	_, err := base64.StdEncoding.DecodeString(string(s))
	if err != nil {
		return fmt.Errorf("invalid base64 string")
	}

	return nil
}

// DescriptionString5to50Bytes is a string used for descriptions, with a length between 5 and 50 characters.
type DescriptionString5to50Bytes string

// Validate checks if the string is between 5 and 50 characters long
func (s DescriptionString5to50Bytes) Validate() error {
	if len(s) < 5 {
		return fmt.Errorf("at least 5 length")
	}
	if len(s) > 50 {
		return fmt.Errorf("no more than 50 length")
	}
	return nil
}

// PositiveIntegerDefault10Max10000  is a positive integer that defaults to 10, and has an upper bound of 10000.
type PositiveIntegerDefault10Max10000 uint

// Validate checks if the integer is maximum 10000
func (i PositiveIntegerDefault10Max10000) Validate() error {
	if i > 10000 {
		return fmt.Errorf("is larger than 10000")
	}

	return nil
}

// CertificateFieldNameUnder50Bytes Represents a certificate field name with a maximum length of 50 characters
type CertificateFieldNameUnder50Bytes string

// Validate checks if the string is under 50 length
func (s CertificateFieldNameUnder50Bytes) Validate() error {
	if len(s) < 1 {
		return fmt.Errorf("at least 1 length")
	}

	if len(s) > 50 {
		return fmt.Errorf("no more than 50 length")
	}
	return nil
}

// PubKeyHex is a compressed DER secp256k1 public key, exactly 66 hex characters (33 bytes) in length.
type PubKeyHex HexString

// Validate checks if the string is valid pubkey hexadecimal string
func (pkh PubKeyHex) Validate() error {
	// The public key is stored as a hex string, which means each byte is represented by 2 characters.
	// To get the actual byte length of the public key, we divide the hex string length by 2.
	pkhHalfLen := len(pkh) / 2
	if pkhHalfLen != pubKeyBytesLenCompressed && pkhHalfLen != pubKeyBytesLenUncompressed {
		return fmt.Errorf("invalid pubKey hex length: %d", len(pkh))
	}

	// Validate as HexString
	hs := HexString(pkh)
	if err := hs.Validate(); err != nil {
		return fmt.Errorf("invalid pubKey hex string: %w", err)
	}

	return nil
}

// HexString is a string in hexadecimal format
type HexString string

var hexRegex = regexp.MustCompile("^[0-9a-fA-F]+$")

// Validate checks if the string is a valid hexadecimal string
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

// Value returns the boolean value with a default when nil
func (b *BooleanDefaultTrue) Value() bool {
	if b == nil {
		return true
	}
	return bool(*b)
}

// BooleanDefaultFalse is a boolean with a default value of false
type BooleanDefaultFalse bool

// Value returns the boolean value with a default when nil
func (b *BooleanDefaultFalse) Value() bool {
	if b == nil {
		return false
	}
	return bool(*b)
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

// TXIDHexString is a hexadecimal transaction ID
type TXIDHexString string

// BEEF An array of integers, each ranging from 0 to 255, indicating transaction data in BEEF (BRC-62) format.
type BEEF []byte

// OutpointString represents a transaction ID and output index pair.
// The TXID is given as a hex string followed by a period "." and then the output index is given as a decimal integer.
type OutpointString string

// Validate checks if the string is proper outpoint string and contains outpoint index after "."
func (s OutpointString) Validate() error {
	split := strings.Split(string(s), ".")

	if len(split) != 2 {
		return fmt.Errorf("txid as hexstring and numeric output index joined with '.'")
	}

	// check if after decimal point there is an outpoint index
	_, err := to.UInt64FromString(split[1])
	if err != nil {
		return fmt.Errorf("txid as hexstring and numeric output index joined with '.'")
	}

	return nil
}

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

// ReqHistoryNote is the history representation of the request
type ReqHistoryNote struct {
	When        *string
	What        string
	ExtraFields map[string]interface{}
}
