package actions

import (
	"context"
	"fmt"
	"iter"
	"log/slog"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/logging"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/txutils"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	"github.com/go-softwarelab/common/pkg/seq"
	"github.com/go-softwarelab/common/pkg/seqerr"
)

type UTXO struct {
	TxID     string
	Vout     uint32
	Satoshis uint64
}

type FundingResult struct {
	AllocatedUTXOs []*UTXO
	ChangeCount    int
	ChangeAmount   uint64
	Fee            uint64
}

type CreateActionParams struct {
	Outputs iter.Seq[*wdk.ValidCreateActionOutput]
	Inputs  iter.Seq[*wdk.ValidCreateActionInput]
}

func FromValidCreateActionArgs(args *wdk.ValidCreateActionArgs) CreateActionParams {
	// TODO: use only the necessary fields (no redundant fields)
	return CreateActionParams{
		Outputs: seq.PointersFromSlice(args.Outputs),
		Inputs:  seq.PointersFromSlice(args.Inputs),
	}
}

type Funder interface {
	// Fund
	// @param targetSat - the target amount of satoshis to fund (total inputs - total outputs)
	// @param currentTxSize - the current size of the transaction in bytes (size of tx + current inputs + current outputs)
	// @param numberOfDesiredUTXOs - the number of UTXOs in basket #TakeFromBasket
	// @param minimumDesiredUTXOValue - the minimum value of UTXO in basket #TakeFromBasket
	// @param userID - the user ID
	Fund(ctx context.Context, targetSat int64, currentTxSize uint64, numberOfDesiredUTXOs int, minimumDesiredUTXOValue uint64, userID int) (*FundingResult, error)
}

type BasketRepo interface {
	FindByName(userID int, name string) (*wdk.TableOutputBasket, error)
}

type create struct {
	logger     *slog.Logger
	funder     Funder
	basketRepo BasketRepo
}

func newCreateAction(logger *slog.Logger, funder Funder, basketRepo BasketRepo) *create {
	logger = logging.Child(logger, "createAction")
	return &create{
		logger:     logger,
		funder:     funder,
		basketRepo: basketRepo,
	}
}

func (c *create) Create(auth wdk.AuthID, args CreateActionParams) (*wdk.StorageCreateActionResult, error) {
	basket, err := c.basketRepo.FindByName(*auth.UserID, wdk.BasketNameForChange)
	if err != nil {
		return nil, fmt.Errorf("failed to find basket: %w", err)
	}
	if basket == nil {
		return nil, fmt.Errorf("basket for change (%s) not found", wdk.BasketNameForChange)
	}

	initialTxSize, err := c.txSize(&args)
	if err != nil {
		return nil, err
	}

	_, err = c.funder.Fund(context.Background(), 0, initialTxSize, basket.NumberOfDesiredUTXOs, basket.MinimumDesiredUTXOValue, *auth.UserID)
	if err != nil {
		return nil, fmt.Errorf("funding failed: %w", err)
	}

	return &wdk.StorageCreateActionResult{}, nil
}

func (c *create) txSize(args *CreateActionParams) (uint64, error) {
	outputSizes := seqerr.MapSeq(args.Outputs, func(o *wdk.ValidCreateActionOutput) (uint64, error) {
		return o.ScriptLength()
	})

	inputSizes := seqerr.MapSeq(args.Inputs, func(o *wdk.ValidCreateActionInput) (uint64, error) {
		return o.ScriptLength()
	})

	txSize, err := txutils.TransactionSize(inputSizes, outputSizes)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate transaction size: %w", err)
	}

	return txSize, nil
}
