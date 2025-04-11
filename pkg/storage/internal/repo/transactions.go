package repo

import (
	"context"
	"fmt"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/database/models"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	"gorm.io/gorm"
)

type Transactions struct {
	db *gorm.DB
}

func NewTransactions(db *gorm.DB) *Transactions {
	return &Transactions{db: db}
}

func (txs *Transactions) CreateTransaction(ctx context.Context, newTx *wdk.NewTx) error {
	err := txs.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		model := &models.Transaction{
			UserID:      newTx.UserID,
			Status:      newTx.Status,
			Reference:   newTx.Reference,
			IsOutgoing:  newTx.IsOutgoing,
			Satoshis:    newTx.Satoshis,
			Description: newTx.Description,
			Version:     newTx.Version,
			LockTime:    newTx.LockTime,
			InputBeef:   newTx.InputBeef,
			RawTx:       nil,
			TxID:        nil,
			Labels:      nil,
		}
		for _, label := range newTx.Labels {
			model.Labels = append(model.Labels, models.Label{
				Name:   string(label),
				UserID: newTx.UserID,
			})
		}

		return tx.Create(model).Error
	})
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}
	return nil
}
