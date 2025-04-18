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
	Description string `gorm:"type:string"`
	Version     int
	LockTime    int
	TxID        *string
	InputBeef   []byte
	RawTx       []byte

	Outputs []Output `gorm:"foreignKey:TransactionID"`
	Labels  []Label  `gorm:"many2many:transaction_labels;"`
}
