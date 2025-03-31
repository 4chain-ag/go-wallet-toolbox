package funder

import (
	"context"
	"log/slog"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/defs"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/logging"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/actions"
	"gorm.io/gorm"
)

type SQL struct {
	logger   *slog.Logger
	db       *gorm.DB
	feeModel defs.FeeModel
}

func NewSQL(logger *slog.Logger, db *gorm.DB, feeModel defs.FeeModel) *SQL {
	logger = logging.Child(logger, "funderSQL")
	return &SQL{
		logger:   logger,
		db:       db,
		feeModel: feeModel,
	}
}

// Fund
// @param targetSat - the target amount of satoshis to fund (total inputs - total outputs)
// @param currentTxSize - the current size of the transaction in bytes (size of tx + current inputs + current outputs)
// @param numberOfDesiredUTXOs - the number of UTXOs in basket #TakeFromBasket
// @param minimumDesiredUTXOValue - the minimum value of UTXO in basket #TakeFromBasket
// @param userID - the user ID.
func (f *SQL) Fund(ctx context.Context, targetSat int64, currentTxSize int64, numberOfDesiredUTXOs int, minimumDesiredUTXOValue uint64, userID int) (*actions.FundingResult, error) {
	panic("not implemented")
}
