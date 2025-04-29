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
	err := txs.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		basketMaker := newCachedBasketMaker(tx, newTx.UserID)

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
		for _, newOut := range newTx.Outputs {
			var basketID *int
			if newOut.Basket != nil {
				var err error
				basketID, err = basketMaker.findOrCreate(ctx, *newOut.Basket, wdk.DefaultNumberOfDesiredUTXOs, wdk.DefaultMinimumDesiredUTXOValue)
				if err != nil {
					return fmt.Errorf("failed to find or create output basket: %w", err)
				}
			}

			output, err := txs.makeNewOutput(ctx, newTx.UserID, newOut, basketID)
			if err != nil {
				return err
			}

			model.Outputs = append(model.Outputs, output)
		}
		for _, label := range newTx.Labels {
			model.Labels = append(model.Labels, &models.Label{
				Name:   string(label),
				UserID: newTx.UserID,
			})
		}

		for _, reservedOutputID := range newTx.ReservedOutputIDs {
			model.ReservedUtxos = append(model.ReservedUtxos, &models.UserUTXO{
				UserID:   newTx.UserID,
				OutputID: reservedOutputID,
			})
		}

		if err := txs.markReservedOutputsAsNotSpendable(tx, newTx.UserID, newTx.ReservedOutputIDs); err != nil {
			return err
		}

		return tx.Create(model).Error
	})
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}
	return nil
}

func (txs *Transactions) makeNewOutput(ctx context.Context, userID int, output *entity.NewOutput, basketID *int) (*models.Output, error) {
	out := models.Output{
		Vout:               output.Vout,
		UserID:             userID,
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
		BasketID:           basketID,
	}

	if out.Spendable && out.Change {
		if basketID == nil {
			return nil, fmt.Errorf("basket ID is nil for change output")
		}
		if out.Satoshis == 0 {
			return nil, fmt.Errorf("change output with zero satoshis")
		}
		sats, err := to.UInt64(out.Satoshis)
		if err != nil {
			return nil, fmt.Errorf("failed to convert satoshis to uint64: %w", err)
		}

		out.UserUTXO = &models.UserUTXO{
			UserID:             userID,
			BasketID:           *basketID,
			Satoshis:           sats,
			EstimatedInputSize: txutils.EstimatedInputSizeByType(output.Type),
		}
	}
	return &out, nil
}

func (txs *Transactions) markReservedOutputsAsNotSpendable(tx *gorm.DB, userID int, outputIDs []uint) error {
	if len(outputIDs) == 0 {
		return nil
	}

	err := tx.Model(&models.Output{}).
		Where("id IN ?", outputIDs).
		Where("user_id = ?", userID).
		Update("spendable", false).Error
	if err != nil {
		return fmt.Errorf("failed to mark reserved outputs as not spendable: %w", err)
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
