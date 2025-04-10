package wdk

import "github.com/4chain-ag/go-wallet-toolbox/pkg/wdk/primitives"

// NewTx represents all the information necessary to store a transaction with additional information like labels, tags, inputs, and outputs.
// This meant to be used for createAction
type NewTx struct {
	UserID int

	Version     int
	LockTime    int
	Status      TxStatus
	Reference   string
	Satoshis    uint64
	IsOutgoing  bool
	InputBeef   []byte
	Description string

	Labels []primitives.StringUnder300
}
