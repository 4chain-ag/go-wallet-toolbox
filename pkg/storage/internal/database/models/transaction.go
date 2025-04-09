package models

import (
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model

	UserID      uint
	Status      wdk.TxStatus
	Reference   string
	IsOutgoing  bool
	Satoshis    uint64
	Description string
	Version     int
	lockTime    *int
	TxID        *string
	InputBeef   []byte
	RawTx       []byte
}
