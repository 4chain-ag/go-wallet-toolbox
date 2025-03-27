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
	if auth.UserID == nil {
		return nil, fmt.Errorf("missing user ID")
	}
	if err := c.validateArgs(&args); err != nil {
		return nil, err
	}

	result, err := c.funder.Fund(context.Background(), nil, *auth.UserID)
	if err != nil {
		return nil, fmt.Errorf("funding failed: %w", err)
	}
	panic(fmt.Errorf("not implemented %v", result))
}

func (c *create) validateArgs(args *wdk.ValidCreateActionArgs) error {
	deducedIsSendWith := len(args.Options.SendWith) > 0
	if args.IsSendWith != deducedIsSendWith {
		return fmt.Errorf("inconsistent IsSendWith with Options.SendWith")
	}

	deducedIsRemixChange := !args.IsSendWith && len(args.Inputs) == 0 && len(args.Outputs) == 0
	if args.IsRemixChange != deducedIsRemixChange {
		return fmt.Errorf("inconsistent IsRemixChange with IsSendWith and Inputs and Outputs")
	}

	deducedIsNewTx := args.IsRemixChange || len(args.Inputs) > 0 || len(args.Outputs) > 0
	if args.IsNewTx != deducedIsNewTx {
		return fmt.Errorf("inconsistent IsNewTx with IsRemixChange and Inputs and Outputs")
	}

	if !args.IsNewTx {
		return fmt.Errorf("create action is meant to create a new transaction")
	}

	deducedIsSignAction := args.IsNewTx && !args.Options.SignAndProcess && containsNilUnlockingScript(args.Inputs)
	if args.IsSignAction != deducedIsSignAction {
		return fmt.Errorf("inconsistent IsSignAction with IsNewTx and Options.SignAndProcess and Inputs.UnlockingScript")
	}

	deducedIsDelayed := bool(args.Options.AcceptDelayedBroadcast)
	if args.IsDelayed != deducedIsDelayed {
		return fmt.Errorf("inconsistent IsDelayed with Options.AcceptDelayedBroadcast")
	}

	deducedIsNoSend := bool(args.Options.NoSend)
	if args.IsNoSend != deducedIsNoSend {
		return fmt.Errorf("inconsistent IsNoSend with Options.NoSend")
	}

	return nil
}

func containsNilUnlockingScript(inputs []wdk.ValidCreateActionInput) bool {
	for _, input := range inputs {
		if input.UnlockingScript == nil {
			return true
		}
	}
	return false
}
