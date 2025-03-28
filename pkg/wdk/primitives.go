package wdk

// DescriptionString5to2000Bytes represents a string used for descriptions,
// with a length between 5 and 2000 characters.
type DescriptionString5to2000Bytes string

// Base64String is a string in base64 format
type Base64String string

// DescriptionString5to50Bytes is a string used for descriptions, with a length between 5 and 50 characters.
type DescriptionString5to50Bytes string

// PositiveIntegerDefault10Max10000  is a positive integer that defaults to 10, and has an upper bound of 10000.
type PositiveIntegerDefault10Max10000 uint

// CertificateFieldNameUnder50Bytes Represents a certificate field name with a maximum length of 50 characters
type CertificateFieldNameUnder50Bytes string

// HexString is a string in hexadecimal format
type HexString string

// PubKeyHex is a compressed DER secp256k1 public key, exactly 66 hex characters (33 bytes) in length.
type PubKeyHex HexString

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

// PositiveIntegerOrZero represents a positive integer or zero value
type PositiveIntegerOrZero uint

// BasketStringUnder300Bytes is a string used for basket names, with a length under 300 bytes
type BasketStringUnder300Bytes string

// TXIDHexString is a hexadecimal transaction ID
type TXIDHexString string

// BEEF An array of integers, each ranging from 0 to 255, indicating transaction data in BEEF (BRC-62) format.
type BEEF []byte

// OutpointString represents a transaction ID and output index pair.
// The TXID is given as a hex string followed by a period "." and then the output index is given as a decimal integer.
type OutpointString string

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
