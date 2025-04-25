package actions

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/logging"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/satoshi"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/txutils"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/entity"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk/primitives"
	"github.com/bsv-blockchain/go-sdk/transaction"
	"github.com/go-softwarelab/common/pkg/must"
	"github.com/go-softwarelab/common/pkg/to"
)

type TransactionsRepo interface {
	CreateTransaction(ctx context.Context, transaction *entity.NewTx) error
	FindTransactionByUserIDAndTxID(ctx context.Context, userID int, txID string) (*wdk.TableTransaction, error)
}

type internalize struct {
	logger *slog.Logger
	txRepo TransactionsRepo
}

func newInternalizeAction(logger *slog.Logger, txRepo TransactionsRepo) *internalize {
	logger = logging.Child(logger, "internalizeAction")
	return &internalize{
		logger: logger,
		txRepo: txRepo,
	}
}

func (in *internalize) Internalize(ctx context.Context, userID int, args *wdk.InternalizeActionArgs) (*wdk.InternalizeActionResult, error) {
	tx, err := transaction.NewTransactionFromBEEF(args.Tx)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction from BEEF: %w", err)
	}
	txID := tx.TxID().String()

	satoshis := satoshi.Zero()

	storedTx, err := in.txRepo.FindTransactionByUserIDAndTxID(ctx, userID, txID)
	if err != nil {
		return nil, fmt.Errorf("failed to find transaction by userID and txID: %w", err)
	}
	if storedTx != nil {
		panic("not implemented yet") // TODO: Implement internalize action for known transaction (with merge)
	}

	var newOutputs []*entity.NewOutput
	outputsCountU64, err := to.UInt32(len(tx.Outputs))
	if err != nil {
		return nil, fmt.Errorf("failed to convert outputs count to uint32: %w", err)
	}
	for _, outputSpec := range args.Outputs {
		if outputSpec.OutputIndex >= outputsCountU64 {
			return nil, fmt.Errorf("output index %d is out of range of provided tx outputs count %d", outputSpec.OutputIndex, outputsCountU64)
		}

		output := tx.Outputs[outputSpec.OutputIndex]

		switch outputSpec.Protocol {
		case wdk.WalletPaymentProtocol:
			satoshis = satoshi.MustAdd(satoshis, output.Satoshis)

			remittance := outputSpec.PaymentRemittance
			newOutputs = append(newOutputs, &entity.NewOutput{
				Vout:              outputSpec.OutputIndex,
				Spendable:         true,
				LockingScript:     to.Ptr(primitives.HexString(output.LockingScript.String())), //TODO: Check is LockingScript can't be []byte
				Basket:            to.Ptr(wdk.BasketNameForChange),                             // TODO: check if the basket exists before
				Satoshis:          satoshi.MustFrom(output.Satoshis),
				SenderIdentityKey: to.Ptr(string(remittance.SenderIdentityKey)),
				Type:              wdk.OutputTypeP2PKH,
				ProvidedBy:        wdk.ProvidedByStorage,
				Purpose:           "change",
				Change:            true,
				DerivationPrefix:  to.Ptr(string(remittance.DerivationPrefix)),
				DerivationSuffix:  to.Ptr(string(remittance.DerivationSuffix)),
			})

		case wdk.BasketInsertionProtocol:
			remittance := outputSpec.InsertionRemittance
			newOutputs = append(newOutputs, &entity.NewOutput{
				Vout:               outputSpec.OutputIndex,
				Spendable:          true,
				LockingScript:      to.Ptr(primitives.HexString(output.LockingScript.String())), //TODO: Check is LockingScript can't be []byte
				Basket:             to.Ptr(string(remittance.Basket)),
				Satoshis:           satoshi.MustFrom(output.Satoshis),
				Type:               wdk.OutputTypeCustom,
				CustomInstructions: remittance.CustomInstructions,
				Change:             false,
				ProvidedBy:         wdk.ProvidedByYou,
			})
		}
	}

	reference, err := txutils.RandomBase64(referenceLength)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random reference: %w", err)
	}

	err = in.txRepo.CreateTransaction(ctx, &entity.NewTx{
		UserID:      userID,
		Version:     must.ConvertToIntFromUnsigned(tx.Version),  // TODO: Refactor Version fields to be uint32
		LockTime:    must.ConvertToIntFromUnsigned(tx.LockTime), // TODO: Refactor LockTime fields to be uint32
		Status:      wdk.TxStatusUnproven,
		Reference:   reference,
		IsOutgoing:  false,
		Description: string(args.Description),
		Satoshis:    satoshis.Int64(),
		TxID:        to.Ptr(txID),
		Outputs:     newOutputs,
		Labels:      args.Labels,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	return &wdk.InternalizeActionResult{
		Accepted: true,
		IsMerge:  false,
		TxID:     txID,
		Satoshis: primitives.SatoshiValue(satoshis.MustUInt64()),
	}, nil
}
