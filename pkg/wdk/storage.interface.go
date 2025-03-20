package wdk

// DescriptionString5to2000Bytes represents a string used for descriptions,
// with a length between 5 and 2000 characters.
type DescriptionString5to2000Bytes string

// Base64String is a string in base64 format
type Base64String string

// HexString is a string in hexadecimal format
type HexString string

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

// OutPoint identifies a unique transaction output by its txid and index vout
type OutPoint struct {
	// TxID Transaction double sha256 hash as big endian hex string
	TxID string
	// Vout Zero based output index within the transaction
	Vout int
}

// ValidCreateActionInput represents the input for a transaction action
type ValidCreateActionInput struct {
	Outpoint              OutPoint
	InputDescription      DescriptionString5to2000Bytes
	SequenceNumber        PositiveIntegerOrZero
	UnlockingScript       *HexString
	UnlockingScriptLength PositiveInteger
}

// ValidCreateActionOutput represents the output for a transaction action
type ValidCreateActionOutput struct {
	LockingScript      HexString
	Satoshis           SatoshiValue
	OutputDescription  DescriptionString5to2000Bytes
	Basket             *BasketStringUnder300Bytes
	CustomInstructions *string
	Tags               []BasketStringUnder300Bytes
}

// ValidProcessActionOptions represents options for processing an action
type ValidProcessActionOptions struct {
	AcceptDelayedBroadcast BooleanDefaultTrue
	ReturnTXIDOnly         BooleanDefaultFalse
	NoSend                 BooleanDefaultFalse
	SendWith               []TXIDHexString
}

// ValidCreateActionOptions extends ValidProcessActionOptions with additional options
type ValidCreateActionOptions struct {
	ValidProcessActionOptions
	SignAndProcess   bool
	TrustSelf        *string
	KnownTxids       []TXIDHexString
	NoSendChange     []OutPoint
	RandomizeOutputs bool
}

// ValidProcessActionArgs represents arguments for processing an action.
// It contains the core parameters needed to process a transaction.
type ValidProcessActionArgs struct {
	// Options contains configuration settings for how the action should be processed
	Options ValidProcessActionOptions
	// IsSendWith is true if a batch of transactions is included for processing
	IsSendWith bool
	// IsNewTx is true if there is a new transaction (not no inputs and no outputs)
	IsNewTx bool
	// IsRemixChange is true if this is a request to remix change
	// When true, IsNewTx will also be true and IsSendWith must be false
	IsRemixChange bool
	// IsNoSend is true if any new transaction should NOT be sent to the network
	IsNoSend bool
	// IsDelayed is true if options.AcceptDelayedBroadcast is true
	IsDelayed bool
}

// ValidCreateActionArgs represents the arguments for creating a transaction action
type ValidCreateActionArgs struct {
	Description                  DescriptionString5to2000Bytes
	InputBEEF                    BEEF
	Inputs                       []ValidCreateActionInput
	Outputs                      []ValidCreateActionOutput
	LockTime                     int
	Version                      int
	Labels                       []string
	IsSignAction                 bool
	RandomVals                   *[]int
	IncludeAllSourceTransactions bool

	// Below are args from ValidProcessActionArgs

	// IsSendWith is true if a batch of transactions is included for processing
	IsSendWith bool
	// IsNewTx is true if there is a new transaction (not no inputs and no outputs)
	IsNewTx bool
	// IsRemixChange is true if this is a request to remix change
	// When true, IsNewTx will also be true and IsSendWith must be false
	IsRemixChange bool
	// IsNoSend is true if any new transaction should NOT be sent to the network
	IsNoSend bool
	// IsDelayed is true if options.AcceptDelayedBroadcast is true
	IsDelayed bool
}

// AuthID represents the identity of the user making the request
type AuthID struct {
	IdentityKey string
	UserID      *int
	IsActive    *bool
}

// WalletStorageWriter is an interface for writing to the wallet storage
type WalletStorageWriter interface {
	CreateAction(auth AuthID, args ValidCreateActionArgs)
}
