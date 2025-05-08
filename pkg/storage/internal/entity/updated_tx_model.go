package entity

import "github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"

type UpdatedTx struct {
	Spendable   bool
	TxID        string
	TxStatus    wdk.TxStatus
	ReqTxStatus wdk.ProvenTxReqStatus
	InputBeef   []byte
	RawTx       []byte
}
