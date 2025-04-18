package repo

import (
	"context"
	"fmt"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/database/models"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/entity"
	"gorm.io/gorm"
)

type Transactions struct {
	db *gorm.DB
}

func NewTransactions(db *gorm.DB) *Transactions {
	return &Transactions{db: db}
}

func (txs *Transactions) CreateTransaction(ctx context.Context, newTx *entity.NewTx) error {
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
		for _, output := range newTx.Outputs {
			out := models.Output{
				Vout:               output.Vout,
				Satoshis:           output.Satoshis.Int64(),
				Spendable:          output.Spendable,
				Change:             output.Change,
				ProvidedBy:         string(output.ProvidedBy),
				Description:        output.Description,
				Purpose:            output.Purpose,
				Type:               string(output.Type),
				DerivationPrefix:   output.DerivationPrefix,
				DerivationSuffix:   output.DerivationSuffix,
				LockingScript:      (*string)(output.LockingScript),
				CustomInstructions: output.CustomInstructions,
			}

			if output.Basket != nil {
				out.Basket = &models.OutputBasket{
					Name:   *output.Basket,
					UserID: newTx.UserID,
				}
			}

			model.Outputs = append(model.Outputs, out)
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
