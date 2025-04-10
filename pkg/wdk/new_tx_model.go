package wdk

import "github.com/4chain-ag/go-wallet-toolbox/pkg/wdk/primitives"

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
