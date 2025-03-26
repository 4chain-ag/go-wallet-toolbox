package funder

import (
	"context"
	"log/slog"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/logging"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/actions"
	sdk "github.com/bsv-blockchain/go-sdk/transaction"
	"gorm.io/gorm"
)

type SQL struct {
	logger *slog.Logger
	db     *gorm.DB
}

func NewSQL(logger *slog.Logger, db *gorm.DB) *SQL {
	logger = logging.Child(logger, "funderSQL")
	return &SQL{
		logger: logger,
		db:     db,
	}
}

func (f *SQL) Fund(ctx context.Context, tx *sdk.Transaction, userID int) (*actions.FundingResult, error) {
	panic("not implemented")
}
