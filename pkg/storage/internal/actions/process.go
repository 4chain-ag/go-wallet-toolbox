package actions

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/logging"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/entity"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/history"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	"github.com/bsv-blockchain/go-sdk/transaction"
	"github.com/go-softwarelab/common/pkg/must"
)

type process struct {
	logger       *slog.Logger
	txRepo       TransactionsRepo
	outputRepo   OutputRepo
	provenTxRepo ProvenTxRepo
}

func newProcessAction(logger *slog.Logger, txRepo TransactionsRepo, outputRepo OutputRepo, provenTxRepo ProvenTxRepo) *process {
	logger = logging.Child(logger, "processAction")
	return &process{
		logger:       logger,
		txRepo:       txRepo,
		outputRepo:   outputRepo,
		provenTxRepo: provenTxRepo,
	}
}

func (p *process) Process(ctx context.Context, userID int, args *wdk.ProcessActionArgs) (*wdk.ProcessActionResult, error) {
	if args.IsNewTx {
		err := p.processNewTx(ctx, userID, args)
		if err != nil {
			return nil, err
		}
	}

	if args.IsSendWith {
		panic("not implemented yet")
	}

	if args.IsDelayed {
		panic("not implemented yet")
	}

	if args.IsNoSend {
		return &wdk.ProcessActionResult{
			SendWithResults: make([]wdk.SendWithResult, 0),
		}, nil
	}

	_, err := p.broadcastSingleTx(ctx, string(*args.TxID))
	if err != nil {
		return nil, fmt.Errorf("failed to broadcast transaction: %w", err)
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

	err = p.validateStateOfTableTx(*args.Reference, tableTx)
	if err != nil {
		return err
	}

	_, outputs, err := p.outputRepo.FindInputsAndOutputsOfTransaction(ctx, tableTx.TransactionID)
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

	newTxStatus, newReqStatus := p.newStatuses(args)

	err = p.txRepo.UpdateTransaction(ctx, entity.UpdatedTx{
		UserID:        userID,
		TransactionID: tableTx.TransactionID,
		Spendable:     true,
		TxID:          txID,
		TxStatus:      newTxStatus,
		ReqTxStatus:   newReqStatus,
		RawTx:         args.RawTx,
		InputBeef:     tableTx.InputBEEF,
	}, history.ProcessActionHistoryNote, history.UserIDHistoryAttr(userID))
	if err != nil {
		return fmt.Errorf("failed to update transaction: %w", err)
	}

	return nil
}

func (p *process) validateStateOfTableTx(reference string, tableTx *wdk.TableTransaction) error {
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
		if output.Change {
			continue
		}

		if output.LockingScript == nil {
			return fmt.Errorf("locking script is nil for output %d", output.OutputID)
		}

		voutInt := must.ConvertToIntFromUnsigned(output.Vout)
		if voutInt >= len(tx.Outputs) {
			return fmt.Errorf("output index %d is out of range of provided tx outputs count %d", voutInt, len(tx.Outputs))
		}

		fromDB := *output.LockingScript
		providedInArgs := tx.Outputs[voutInt].LockingScript.String()
		if providedInArgs != fromDB {
			return fmt.Errorf("locking script mismatch at vout: %d, provided %s, calculated from raw tx: %s", voutInt, providedInArgs, fromDB)
		}
	}
	return nil
}

func (p *process) newStatuses(args *wdk.ProcessActionArgs) (txStatus wdk.TxStatus, reqStatus wdk.ProvenTxReqStatus) {
	switch {
	case args.IsNoSend:
		reqStatus = wdk.ProvenTxStatusNoSend
		txStatus = wdk.TxStatusNoSend
	case args.IsDelayed:
		reqStatus = wdk.ProvenTxStatusUnsent
		txStatus = wdk.TxStatusUnprocessed
	default:
		reqStatus = wdk.ProvenTxStatusUnprocessed
		txStatus = wdk.TxStatusUnprocessed
	}

	return
}

func (p *process) broadcastSingleTx(ctx context.Context, txID string) (wdk.SendWithResultStatus, error) {
	sendStatus, err := p.sendStatusByReqTxStatus(ctx, txID)
	if err != nil {
		return "", err
	}

	if sendStatus != wdk.SendWithResultStatusSending {
		return sendStatus, nil
	}

	beef, err := p.provenTxRepo.BuildValidBEEF(ctx, txID, wdk.ProvenTxReqStatusesForSourceTransactions)
	if err != nil {
		return "", fmt.Errorf("failed to build valid BEEF: %w", err)
	}

	// TODO: SPV of the beef

	_ = beef // TODO Services::PostBEEF
	return wdk.SendWithResultStatusSending, nil
}

func (p *process) sendStatusByReqTxStatus(ctx context.Context, txID string) (wdk.SendWithResultStatus, error) {
	reqTxStatus, err := p.provenTxRepo.FindProvenTxStatus(ctx, txID)
	if err != nil {
		return "", fmt.Errorf("failed to find proven tx status: %w", err)
	}

	switch reqTxStatus.BroadcastStatus() {
	case wdk.TxReqBroadcastReadyToSend:
		return wdk.SendWithResultStatusSending, nil
	case wdk.TxReqBroadcastError:
		return wdk.SendWithResultStatusFailed, nil
	case wdk.TxReqBroadcastAlreadySent:
		return wdk.SendWithResultStatusUnproven, nil
	case wdk.TxReqBroadcastUnknown:
		fallthrough
	default:
		return "", fmt.Errorf("unknown broadcast status")
	}
}
