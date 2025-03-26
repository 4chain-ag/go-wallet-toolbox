package actions

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/logging"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	sdk "github.com/bsv-blockchain/go-sdk/transaction"
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

type Funder interface {
	Fund(ctx context.Context, tx *sdk.Transaction, userID int) (*FundingResult, error)
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

func (c *create) Create(auth wdk.AuthID, args wdk.ValidCreateActionArgs) (*wdk.StorageCreateActionResult, error) {
	result, err := c.funder.Fund(context.Background(), nil, *auth.UserID)
	if err != nil {
		return nil, fmt.Errorf("funding failed: %w", err)
	}
	panic(fmt.Errorf("not implemented %v", result))
}
