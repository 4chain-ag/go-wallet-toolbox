package actions

import (
	"context"
	"fmt"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/logging"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	"github.com/bsv-blockchain/go-sdk/transaction"
	"github.com/go-softwarelab/common/pkg/must"
	"log/slog"
)

type process struct {
	logger     *slog.Logger
	txRepo     TransactionsRepo
	outputRepo OutputRepo
}

func newProcessAction(logger *slog.Logger, txRepo TransactionsRepo, outputRepo OutputRepo) *process {
	logger = logging.Child(logger, "processAction")
	return &process{
		logger:     logger,
		txRepo:     txRepo,
		outputRepo: outputRepo,
	}
}

func (p *process) Process(ctx context.Context, userID int, args *wdk.ProcessActionArgs) (*wdk.ProcessActionResult, error) {
	if args.IsNewTx {
		err := p.processNewTx(ctx, userID, args)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (p *process) processNewTx(ctx context.Context, userID int, args *wdk.ProcessActionArgs) error {
	tx, err := transaction.NewTransactionFromBytes(args.RawTx)
	if err != nil {
		return fmt.Errorf("failed to build transaction object from raw tx bytes: %w", err)
	}
	txID := tx.TxID().String()

	if txID != string(*args.TxID) {
		return fmt.Errorf("txID mismatch: provided %s, calculated from raw tx: %s", *args.TxID, txID)
	}

	// TODO: Services::nLockTimeIsFinal(tx)

	tableTx, err := p.txRepo.FindTransactionByReference(ctx, userID, *args.Reference)
	if err != nil {
		return fmt.Errorf("failed to find transaction by reference: %w", err)
	}

	err = p.validateNewTxTable(*args.Reference, tableTx)
	if err != nil {
		return err
	}

	_, outputs, err := p.outputRepo.FindInputsAndOutputsOfTransaction(ctx, userID, txID)
	if err != nil {
		return fmt.Errorf("failed to find inputs and outputs of transaction: %w", err)
	}

	err = p.validateNewTxOutputs(tx, outputs)
	if err != nil {
		return err
	}

	// TODO: Commission; but it requires Commission table (it needs to be created & new rows added during "createAction"

	// TODO: Add db transactionID to ProvenTxReq.Notify

	// TODO: Remove too long locking scripts (len > storage.maxOutputScript)

	return nil
}

func (p *process) validateNewTxTable(reference string, tableTx *wdk.TableTransaction) error {
	if tableTx == nil {
		return fmt.Errorf("transaction with reference (%s) not found in the database", reference)
	}
	if !tableTx.IsOutgoing {
		return fmt.Errorf("transaction with reference (%s) is not outgoing", reference)
	}
	if len(tableTx.InputBEEF) == 0 {
		return fmt.Errorf("transaction with reference (%s) has no inputBEEF", reference)
	}
	if tableTx.Status != wdk.TxStatusUnsigned && tableTx.Status != wdk.TxStatusUnprocessed {
		return fmt.Errorf("transaction with reference (%s) is not in a valid status for processing", reference)
	}
	return nil
}

func (p *process) validateNewTxOutputs(tx *transaction.Transaction, outputs []*wdk.TableOutput) error {
	for _, output := range outputs {
		voutInt := must.ConvertToIntFromUnsigned(output.Vout)
		if output.LockingScript != nil {
			if voutInt >= len(tx.Outputs) {
				return fmt.Errorf("output index %d is out of range of provided tx outputs count %d", voutInt, len(tx.Outputs))
			}
			fromDB := *output.LockingScript
			providedInArgs := tx.Outputs[voutInt].LockingScript.String()
			if providedInArgs != fromDB {
				return fmt.Errorf("locking script mismatch as vout: %d, provided %s, calculated from raw tx: %s", voutInt, providedInArgs, fromDB)
			}
		}
	}
	return nil
}
