package actions

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/logging"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
)

type UTXO struct {
	TxID     string
	Vout     uint32
	Satoshis uint64
}

type FundingResult struct {
	AllocatedUTXOs []UTXO
	ChangeCount    int
	ChangeAmount   uint64
	Fee            uint64
}

type CreateActionParams struct {
}

func FromValidCreateActionArgs(args *wdk.ValidCreateActionArgs) CreateActionParams {
	// TODO: use only the necessary fields (no redundant fields)
	return CreateActionParams{}
}

type Funder interface {
	// Fund
	// @param targetSat - the target amount of satoshis to fund (total inputs - total outputs)
	// @param currentTxSize - the current size of the transaction in bytes (size of tx + current inputs + current outputs)
	// @param numberOfDesiredUTXOs - the number of UTXOs in basket #TakeFromBasket
	// @param minimumDesiredUTXOValue - the minimum value of UTXO in basket #TakeFromBasket
	// @param userID - the user ID
	Fund(ctx context.Context, targetSat int64, currentTxSize int64, numberOfDesiredUTXOs int, minimumDesiredUTXOValue uint64, userID int) (*FundingResult, error)
}

type create struct {
	logger *slog.Logger
	funder Funder
}

func newCreateAction(logger *slog.Logger, funder Funder) *create {
	logger = logging.Child(logger, "createAction")
	return &create{
		logger: logger,
		funder: funder,
	}
}

func (c *create) Create(auth wdk.AuthID, args CreateActionParams) (*wdk.StorageCreateActionResult, error) {
	result, err := c.funder.Fund(context.Background(), 0, 0, 0, 0, *auth.UserID)
	if err != nil {
		return nil, fmt.Errorf("funding failed: %w", err)
	}
	panic(fmt.Errorf("not implemented %v", result))
}
