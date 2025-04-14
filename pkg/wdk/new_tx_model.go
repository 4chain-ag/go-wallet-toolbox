package wdk

import (
	"iter"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk/primitives"
)

// NewTx represents all the information necessary to store a transaction with additional information like labels, tags, inputs, and outputs.
// This meant to be used for createAction
type NewTx struct {
	UserID int

	Version     int
	LockTime    int
	Status      TxStatus
	Reference   string
	Satoshis    int64
	IsOutgoing  bool
	InputBeef   []byte
	Description string

	Outputs iter.Seq[*NewOutput]

	Labels []primitives.StringUnder300
}

// NewOutput represents an output of a new transaction.
type NewOutput struct {
	Satoshis         int64
	Basket           string
	Spendable        bool
	Change           bool
	ProvidedBy       ProvidedBy
	Purpose          string
	Type             string
	DerivationPrefix string
	DerivationSuffix string
}
