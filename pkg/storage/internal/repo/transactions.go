package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/txutils"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/database/models"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/database/scopes"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/entity"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk/primitives"
	"github.com/go-softwarelab/common/pkg/must"
	"github.com/go-softwarelab/common/pkg/to"
	"gorm.io/gorm"
)

type Transactions struct {
	db *gorm.DB
}

func NewTransactions(db *gorm.DB) *Transactions {
	return &Transactions{db: db}
}

func (txs *Transactions) CreateTransaction(ctx context.Context, newTx *entity.NewTx) error {
	basketNameToID := make(map[string]int)

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
			TxID:        newTx.TxID,
			Labels:      nil,
		}
		for _, output := range newTx.Outputs {
			var basketID *int
			out := models.Output{
				Vout:               output.Vout,
				UserID:             newTx.UserID,
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
				SenderIdentityKey:  output.SenderIdentityKey,
			}

			if output.Basket != nil {
				if cachedBasketID, ok := basketNameToID[*output.Basket]; ok {
					basketID = to.Ptr(cachedBasketID)
				} else {
					var basket models.OutputBasket
					err := tx.
						Where(models.OutputBasket{UserID: newTx.UserID, Name: *output.Basket}).
						FirstOrCreate(&basket).Error
					if err != nil {
						return fmt.Errorf("failed to find output basket: %w", err)
					}

					basketNameToID[*output.Basket] = basket.BasketID
					basketID = to.Ptr(basket.BasketID)
				}

				out.BasketID = basketID
			}

			if out.Spendable && out.Change && basketID != nil {
				if out.Satoshis == 0 {
					return fmt.Errorf("change output with zero satoshis")
				}
				sats, err := to.UInt64(out.Satoshis)
				if err != nil {
					return fmt.Errorf("failed to convert satoshis to uint64: %w", err)
				}

				out.UserUTXO = &models.UserUTXO{
					UserID:             newTx.UserID,
					BasketID:           *basketID,
					Satoshis:           sats,
					EstimatedInputSize: txutils.EstimatedInputSizeByType(output.Type),
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

		for _, reservedOutputID := range newTx.ReservedOutputIDs {
			model.ReservedUtxos = append(model.ReservedUtxos, models.UserUTXO{
				UserID:   newTx.UserID,
				OutputID: reservedOutputID,
			})
		}

		return tx.Create(model).Error
	})
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}
	return nil
}

func (txs *Transactions) FindTransactionByUserIDAndTxID(ctx context.Context, userID int, txID string) (*wdk.TableTransaction, error) {
	var transaction models.Transaction
	err := txs.db.WithContext(ctx).Scopes(scopes.UserID(userID)).Where("tx_id = ?", txID).First(&transaction).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find transaction: %w", err)
	}

	return &wdk.TableTransaction{
		CreatedAt:     transaction.CreatedAt,
		UpdatedAt:     transaction.UpdatedAt,
		TransactionID: transaction.ID,
		UserID:        transaction.UserID,
		Status:        transaction.Status,
		Reference:     primitives.Base64String(transaction.Reference),
		IsOutgoing:    transaction.IsOutgoing,
		Satoshis:      primitives.SatoshiValue(must.ConvertToUInt64(transaction.Satoshis)),
		Description:   transaction.Description,
		Version:       to.Ptr(transaction.Version),
		LockTime:      to.Ptr(transaction.LockTime),
		TxID:          transaction.TxID,
		InputBEEF:     transaction.InputBeef,
		RawTx:         transaction.RawTx,
	}, nil

}
