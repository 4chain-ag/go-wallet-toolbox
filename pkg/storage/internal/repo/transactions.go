package repo

import (
	"context"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	"gorm.io/gorm"
)

type Transactions struct {
	db *gorm.DB
}

func NewTransactions(db *gorm.DB) *Transactions {
	return &Transactions{db: db}
}

func (tx *Transactions) CreateTransaction(ctx context.Context, newTxModel *wdk.NewTxModel) error {
	// TODO
	return nil
}
