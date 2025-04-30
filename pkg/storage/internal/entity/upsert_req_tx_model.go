package entity

import (
	"time"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
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
