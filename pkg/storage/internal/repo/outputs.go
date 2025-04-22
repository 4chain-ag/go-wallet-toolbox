package repo

import (
	"context"
	"fmt"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/database/models"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	"github.com/go-softwarelab/common/pkg/seq"
	"github.com/go-softwarelab/common/pkg/slices"
	"gorm.io/gorm"
	"iter"
)

type Outputs struct {
	db *gorm.DB
}

func NewOutputs(db *gorm.DB) *Outputs {
	return &Outputs{db: db}
}

func (o *Outputs) FindOutputs(ctx context.Context, outputIDs iter.Seq[uint]) ([]*wdk.TableOutput, error) {
	count := seq.Count(outputIDs)
	if count == 0 {
		return nil, nil
	}
	idsClause := make([][]any, 0, count)
	for outputID := range outputIDs {
		idsClause = append(idsClause, []any{outputID})
	}

	var outputs []*models.Output
	err := o.db.WithContext(ctx).
		Model(models.Output{}).
		Preload("Transaction.TxID").
		Where("id IN ?", idsClause).
		Find(&outputs).Error

	if err != nil {
		return nil, fmt.Errorf("failed to find outputs: %w", err)
	}

	return slices.Map(outputs, func(output *models.Output) *wdk.TableOutput {
		return &wdk.TableOutput{
			CreatedAt:          output.CreatedAt,
			UpdatedAt:          output.UpdatedAt,
			OutputID:           output.ID,
			UserID:             output.UserID,
			TransactionID:      output.TransactionID,
			BasketID:           output.BasketID,
			Spendable:          output.Spendable,
			Change:             output.Change,
			OutputDescription:  output.Description,
			Vout:               output.Vout,
			Satoshis:           output.Satoshis,
			ProvidedBy:         output.ProvidedBy,
			Purpose:            output.Purpose,
			Type:               output.Type,
			Txid:               output.Transaction.TxID,
			DerivationPrefix:   output.DerivationPrefix,
			DerivationSuffix:   output.DerivationSuffix,
			CustomInstructions: output.CustomInstructions,
			LockingScript:      output.LockingScript,
		}
	}), nil
}
