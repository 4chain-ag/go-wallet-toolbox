package results

import (
	"github.com/bsv-blockchain/go-sdk/transaction"
)

// PostBEEF is the success result of the single service PostBEEF method.
type PostBEEF struct {
	Notes
	TxIDResults []PostTxID
}

// PostBEEFError is the error result of the single service PostBEEF method.
type PostBEEFError struct {
	Cause error
	Notes
	TxIDResults []PostTxID
}

func (p *PostBEEFError) Error() string {
	return p.Cause.Error()
}

type ResultStatus string

const (
	// ResultStatusSuccess indicates that the result was a success.
	ResultStatusSuccess ResultStatus = "success"
	// ResultStatusError indicates that the result was an error.
	ResultStatusError ResultStatus = "error"
)

// PostTxID is the struct representing postTX result for particular TxID
type PostTxID struct {
	Result ResultStatus
	TxID   string
	// AlreadyKnown if true, the transaction was already known to this service. Usually treat as a success.
	// Potentially stop posting to additional transaction processors.
	AlreadyKnown bool
	// DoubleSpend is when service indicated this broadcast double spends at least one input
	// `competingTxs` may be an array of txids that were first seen spends of at least one input.
	DoubleSpend  bool
	BlockHash    string
	BlockHeight  int64
	MerklePath   *transaction.MerklePath
	CompetingTxs []string
	// TODO: Data type is object | string | PostTxResultForTxidError
	Data any
	Notes
	Error error
}
