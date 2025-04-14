package models

import (
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model

	UserID      int
	Status      wdk.TxStatus
	Reference   string
	IsOutgoing  bool
	Satoshis    int64
	Description string
	Version     int
	LockTime    int
	TxID        *string
	InputBeef   []byte
	RawTx       []byte

	Labels []Label `gorm:"many2many:transaction_labels;"`
}
