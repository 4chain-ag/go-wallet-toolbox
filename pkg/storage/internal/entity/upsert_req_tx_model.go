package entity

import (
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	"time"
)

type NewHistoryNote struct {
	When  *time.Time
	What  string
	Attrs map[string]any
}

type UpsertProvenTxReq struct {
	InputBeef           []byte
	RawTx               []byte
	TxID                string
	Status              wdk.ProvenTxReqStatus
	HistoryNote         NewHistoryNote
	NotifyTransactionID uint
}
