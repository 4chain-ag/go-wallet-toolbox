package wdk

import (
	"github.com/bsv-blockchain/go-sdk/transaction"
	"github.com/bsv-blockchain/go-sdk/transaction/chaintracker"
)

// WalletServices is an interface to access functionality implemented by external transaction processing services
type WalletServices interface {
	// ChainTracker returns service which requires `options.chaintracks` be valid.
	ChainTracker() chaintracker.ChainTracker

	// HeaderForHeight returns serialized block header for height on active chain
	HeaderForHeight(height int) ([]int, error)

	// Height returns the height of the active chain
	Height() int

	// BsvExchangeRate returns approximate exchange rate US Dollar / BSV, USD / BSV
	// This is the US Dollar price of one BSV
	BsvExchangeRate() float64

	// FiatExchangeRate returns  approximate exchange rate currency per base.
	FiatExchangeRate(currency Currency, base *Currency) float64

	// RawTx attempts to obtain the raw transaction bytes associated with a 32 byte transaction hash (txid).
	//
	// Cycles through configured transaction processing services attempting to get a valid response.
	//
	// On success:
	// Result txid is the requested transaction hash
	// Result rawTx will be an array containing raw transaction bytes.
	// Result name will be the responding service's identifying name.
	// Returns result without incrementing active service.
	//
	// On failure:
	// Result txid is the requested transaction hash
	// Result mapi will be the first mapi response obtained (service name and response), or null
	// Result error will be the first error thrown (service name and CwiError), or null
	// Increments to next configured service and tries again until all services have been tried.
	RawTx(txid string, useNext bool) (RawTxResult, error)

	// MerklePath attempts to obtain the merkle proof associated with a 32 byte transaction hash (txid).
	//
	// Cycles through configured transaction processing services attempting to get a valid response.
	//
	// On success:
	// Result txid is the requested transaction hash
	// Result proof will be the merkle proof.
	// Result name will be the responding service's identifying name.
	// Returns result without incrementing active service.
	//
	// On failure:
	// Result txid is the requested transaction hash
	// Result mapi will be the first mapi response obtained (service name and response), or null
	// Result error will be the first error thrown (service name and CwiError), or null
	// Increments to next configured service and tries again until all services have been tried.
	//
	MerklePath(txid string, useNext bool) (MerklePathResult, error)

	// PostBeef attempts to post beef with given txIDs
	PostBeef(beef *transaction.Beef, txids []string) ([]*PostBeefResult, error)

	// UtxoStatus attempts to determine the UTXO status of a transaction output.
	//
	// Cycles through configured transaction processing services attempting to get a valid response.
	UtxoStatus(
		output string,
		outputFormat UtxoStatusOutputFormat,
		useNext bool,
	) (UtxoStatusResult, error)

	// HashToHeader attempts to retrieve BlockHeader by its hash
	HashToHeader(hash string) (*BlockHeader, error)

	// NLockTimeIsFinal returns whether the locktime value allows the transaction to be mined at the current chain height
	// TODO: txOrLockTime type = string | number[] | BsvTransaction | number
	NLockTimeIsFinal(txOrLockTime any) bool
}
