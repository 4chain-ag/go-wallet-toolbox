package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/database/models"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/entity"
	"gorm.io/gorm"
)

type ProvenTxReq struct {
	db *gorm.DB
}

func NewProvenTxReqRepo(db *gorm.DB) *ProvenTxReq {
	return &ProvenTxReq{db: db}
}

func (p *ProvenTxReq) UpsertProvenTxReq(ctx context.Context, req *entity.UpsertProvenTxReq) error {
	err := p.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var model models.ProvenTxReq
		err := tx.First(&model, req.TxID).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("cannot upsert proven tx req: %w", err)
		}

		model.Status = req.Status
		model.TxID = req.TxID
		model.RawTx = req.RawTx
		model.InputBeef = req.InputBeef

		// TODO:

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to upsert proven tx req: %w", err)
	}
	return nil
}
