package funder

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/defs"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/logging"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/utils/to"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/actions"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/actions/funder/errfunder"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/database/models"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/paging"
)

type UTXORepository interface {
	FindAllUTXOs(ctx context.Context, userID int, page *paging.Page) ([]*models.UserUTXO, error)
}

type SQL struct {
	logger         *slog.Logger
	utxoRepository UTXORepository
	feeCalculator  *feeCalc
}

func NewSQL(logger *slog.Logger, utxoRepository UTXORepository, feeModel defs.FeeModel) *SQL {
	logger = logging.Child(logger, "funderSQL")
	feeCalculator := newFeeCalculator(feeModel)

	return &SQL{
		logger:         logger,
		utxoRepository: utxoRepository,
		feeCalculator:  feeCalculator,
	}
}

// Fund
// @param targetSat - the target amount of satoshis to fund (total inputs - total outputs)
// @param currentTxSize - the current size of the transaction in bytes (size of tx + current inputs + current outputs)
// @param numberOfDesiredUTXOs - the number of UTXOs in basket #TakeFromBasket
// @param minimumDesiredUTXOValue - the minimum value of UTXO in basket #TakeFromBasket
// @param userID - the user ID.
func (f *SQL) Fund(ctx context.Context, targetSat int64, currentTxSize int64, numberOfDesiredUTXOs int, minimumDesiredUTXOValue uint64, userID int) (*actions.FundingResult, error) {
	txSize, err := to.UInt64(currentTxSize)
	if err != nil {
		return nil, fmt.Errorf("invalid currentTxSize: %w", err)
	}

	txSats, err := to.UInt64(targetSat)
	if err != nil {
		return nil, fmt.Errorf("invalid targetSat: %w", err)
	}

	page := &paging.Page{
		Size:   1000,
		SortBy: "satoshis",
	}

	collector := newCollector(txSats, txSize, f.feeCalculator)

	for {
		utxos, err := f.utxoRepository.FindAllUTXOs(ctx, userID, page)
		if err != nil {
			return nil, fmt.Errorf("couldn't get users utxos: %w", err)
		} else if len(utxos) == 0 && !collector.IsFunded() {
			return nil, errfunder.NotEnoughFunds
		}

		for _, utxo := range utxos {
			err = collector.Allocate(utxo)
			if err != nil {
				return nil, fmt.Errorf("failed to collect utxo: %w", err)
			}
		}

		if collector.IsFunded() {
			break
		}
		page.Next()
	}

	result, err := collector.GetResult()
	if err != nil {
		return nil, fmt.Errorf("failed to get result: %w", err)
	}
	return result, nil
}

type utxoCollector struct {
	txSats      uint64
	txSize      uint64
	satsCovered uint64

	result        *actions.FundingResult
	feeCalculator *feeCalc
}

func newCollector(txSats, txSize uint64, feeCalculator *feeCalc) *utxoCollector {
	return &utxoCollector{
		txSats:        txSats,
		txSize:        txSize,
		feeCalculator: feeCalculator,

		result: &actions.FundingResult{
			AllocatedUTXOs: make([]*actions.UTXO, 0),
			ChangeAmount:   0,
			ChangeCount:    0,
			Fee:            0,
		},
	}
}

func (c *utxoCollector) Allocate(utxo *models.UserUTXO) (err error) {
	c.addToAllocated(utxo)

	err = c.increaseSize(utxo.EstimatedInputSize)
	if err != nil {
		return fmt.Errorf("failed to increase tx size: %w", err)
	}

	c.increaseValue(utxo.Satoshis)
	return nil
}

func (c *utxoCollector) addToAllocated(utxo *models.UserUTXO) {
	c.result.AllocatedUTXOs = append(c.result.AllocatedUTXOs, &actions.UTXO{
		TxID:     utxo.TxID,
		Vout:     utxo.Vout,
		Satoshis: utxo.Satoshis,
	})
}

func (c *utxoCollector) increaseSize(size uint64) (err error) {
	c.txSize += size
	c.result.Fee, err = c.feeCalculator.Calculate(c.txSize)
	if err != nil {
		return fmt.Errorf("failed to calculate fee: %w", err)
	}
	return nil
}

func (c *utxoCollector) increaseValue(satoshis uint64) {
	c.satsCovered += satoshis
}

func (c *utxoCollector) IsFunded() bool {
	return c.satsCovered >= c.satsToCover()
}

func (c *utxoCollector) satsToCover() uint64 {
	return c.txSats + c.result.Fee
}

func (c *utxoCollector) GetResult() (*actions.FundingResult, error) {
	return c.result, nil
}
