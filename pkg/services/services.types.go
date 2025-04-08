package services

import (
	"time"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	"github.com/bsv-blockchain/go-sdk/transaction"
)

// Currency represents supported currency types
type Currency string

// Supported currency types
const (
	USD Currency = "USD"
	EUR Currency = "EUR"
	GBP Currency = "GBP"
)

// UtxoStatusOutputFormat represents supported utxo status output formats
type UtxoStatusOutputFormat string

// Supported utxo status output formats
const (
	HashLE UtxoStatusOutputFormat = "hashLE"
	HashBE UtxoStatusOutputFormat = "hashBE"
	Script UtxoStatusOutputFormat = "script"
)

// FiatExchangeRates is the rate struct for fiat currency
type FiatExchangeRates struct {
	Timestamp time.Time
	Rates     map[string]float64
	Base      string
}

// RawTxResult is result from RawTx method
type RawTxResult struct {
	// TxID is a transaction hash or rawTx
	TxID string
	// Name is the name of the service returning the rawTx or nil if no rawTx
	Name *string
	// RawTx are multiple proofs that may be returned when a transaction also appears in
	// one or more orphaned blocks
	RawTx []string
}

// BaseBlockHeader are fields of 80 byte serialized header in order whose double sha256 hash is a block's hash value
// and the next block's previousHash value.
// All block hash values and merkleRoot values are 32 byte hex string values with the byte order reversed from the serialized byte order.
type BaseBlockHeader struct {
	// Block header version value. Serialized length is 4 bytes.
	Version int
	// PreviousHash is a hash of previous block's block header. Serialized length is 32 bytes.
	PreviousHash string
	// MerkleRoot is root hash of the merkle tree of all transactions in this block. Serialized length is 32 bytes.
	MerkleRoot string
	// Time is block header time value. Serialized length is 4 bytes.
	Time int
	// Bits are block header bits value. Serialized length is 4 bytes.
	Bits int
	// Nonce is block header nonce value. Serialized length is 4 bytes.
	Nonce int
}

// BlockHeader is a base block header with its computed height and hash in its chain
type BlockHeader struct {
	BaseBlockHeader
	// Height is the of the header, starting from zero
	Height uint
	// Hash is the double sha256 hash of the serialized `BaseBlockHeader` fields
	Hash string
}

// MerklePathResult is result from MerklePath method
type MerklePathResult struct {
	// Name is the name of the service returning the proof, or undefined if no proof
	Name *string
	// MerklePath are multiple proofs may be returned when a transaction also appears in
	// one or more orphaned blocks
	MerklePath *transaction.MerklePath
	Header     *BlockHeader
	Notes      []wdk.ReqHistoryNote
}

// UtxoStatusDetails represents details about occurrences of an output script as a UTXO
type UtxoStatusDetails struct {
	// Height is the block height containing the matching unspent transaction output
	// Typically there will be only one, but future orphans can result in multiple values
	Height *int

	// Txid is the transaction hash (txid) of the transaction containing the matching unspent transaction output
	// Typically there will be only one, but future orphans can result in multiple values
	Txid *string

	// Index is the output index in the transaction containing of the matching unspent transaction output
	// Typically there will be only one, but future orphans can result in multiple values
	Index *int

	// Satoshis is the amount of the matching unspent transaction output
	// Typically there will be only one, but future orphans can result in multiple values
	Satoshis *uint
}

// UtxoStatusResult represents the result of a GetUtxoStatus operation
type UtxoStatusResult struct {
	// Name is the name of the service to which the transaction was submitted for processing
	Name string

	// IsUtxo is true if the output is associated with at least one unspent transaction output
	IsUtxo *bool

	// Details contains additional details about occurrences of this output script as a UTXO.
	// Normally there will be one item in the array but due to the possibility of orphan races
	// there could be more than one block in which it is a valid UTXO.
	Details []UtxoStatusDetails
}

// PostTxResultForTxID is the struct representing postTX result for particular TxID
type PostTxResultForTxID struct {
	TxID string
	// AlreadyKnown if true, the transaction was already known to this service. Usually treat as a success.
	// Potentially stop posting to additional transaction processors.
	AlreadyKnown bool
	// DoubleSpend is when service indicated this broadcast double spends at least one input
	// `competingTxs` may be an array of txids that were first seen spends of at least one input.
	DoubleSpend  bool
	BlockHash    *string
	BlockHeight  *int
	MerklePath   *transaction.MerklePath
	CompetingTxs []string
	// TODO: Data type is object | string | PostTxResultForTxidError
	Data  any
	Notes []wdk.ReqHistoryNote
}

// PostBeefResult are properties on array items of result returned from postBeef method
type PostBeefResult struct {
	// Name is the name of the service to which the transaction was submitted for processing
	Name        string
	TxIDResults []PostTxResultForTxID
	// Data is service response object. Use service name and status to infer type of object.
	Data  any
	Notes []wdk.ReqHistoryNote
}
