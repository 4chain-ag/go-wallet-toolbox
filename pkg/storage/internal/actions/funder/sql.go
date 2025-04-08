package funder

import (
	"context"
	"fmt"
	"iter"
	"log/slog"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/defs"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/logging"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/actions"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/actions/funder/errfunder"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/database/models"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/paging"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/txutils"
	"github.com/go-softwarelab/common/pkg/must"
	"github.com/go-softwarelab/common/pkg/seqerr"
	"github.com/go-softwarelab/common/pkg/to"
)

var changeOutputSize = txutils.P2PKHOutputSize

const utxoBatchSize = 1000

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
func (f *SQL) Fund(ctx context.Context, targetSat int64, currentTxSize uint64, numberOfDesiredUTXOs int, minimumDesiredUTXOValue uint64, userID int) (*actions.FundingResult, error) {
	collector, err := newCollector(targetSat, currentTxSize, numberOfDesiredUTXOs, minimumDesiredUTXOValue, f.feeCalculator)
	if err != nil {
		return nil, fmt.Errorf("failed to start collecting utxo: %w", err)
	}

	utxos := f.loadUTXOs(ctx, userID)

	err = collector.Allocate(utxos)
	if err != nil {
		return nil, fmt.Errorf("failed to allocate utxos: %w", err)
	}

	return collector.GetResult()

}

func (f *SQL) loadUTXOs(ctx context.Context, userID int) iter.Seq2[*models.UserUTXO, error] {
	batches := seqerr.ProduceWithArg(
		func(page *paging.Page) ([]*models.UserUTXO, *paging.Page, error) {
			utxos, err := f.utxoRepository.FindAllUTXOs(ctx, userID, page)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to load utxos: %w", err)
			}
			page.Next()
			return utxos, page, nil
		},
		&paging.Page{
			Limit:  utxoBatchSize,
			SortBy: "satoshis",
		})

	return seqerr.FlattenSlices(batches)
}

type utxoCollector struct {
	txSats int64
	txSize uint64

	fee           int64
	feeCalculator *feeCalc

	satsCovered    int64
	allocatedUTXOs []*actions.UTXO

	numberOfDesiredUTXOs    uint64
	minimumDesiredUTXOValue uint64
	changeOutputsCount      uint64
	minimumChange           uint64
}

func newCollector(txSats int64, txSize uint64, numberOfDesiredUTXOs int, minimumDesiredUTXOValue uint64, feeCalculator *feeCalc) (c *utxoCollector, err error) {
	c = &utxoCollector{
		txSats:                  txSats,
		minimumDesiredUTXOValue: minimumDesiredUTXOValue,
		feeCalculator:           feeCalculator,
		allocatedUTXOs:          make([]*actions.UTXO, 0),
	}

	err = c.increaseSize(txSize)
	if err != nil {
		return nil, fmt.Errorf("failed to increase transaction size: %w", err)
	}

	c.numberOfDesiredUTXOs = must.ConvertToUInt64(to.NoLessThan(numberOfDesiredUTXOs, 1))
	c.minimumChange = minimumDesiredUTXOValue / 4

	err = c.calculateChangeOutputs()
	if err != nil {
		return nil, fmt.Errorf("failed to calculate change outputs: %w", err)
	}

	return c, nil
}

func (c *utxoCollector) Allocate(utxos iter.Seq2[*models.UserUTXO, error]) error {
	utxos = seqerr.TakeUntilTrue(utxos, c.IsFunded)
	err := seqerr.ForEach(utxos, c.allocateUTXO)
	if err != nil {
		return fmt.Errorf("failed to allocate utxo: %w", err)
	}
	return nil
}

func (c *utxoCollector) IsFunded() bool {
	return c.satsCovered >= c.satsToCover()
}

func (c *utxoCollector) GetResult() (*actions.FundingResult, error) {
	if c.IsFunded() {
		return c.prepareResult()
	}
	return nil, errfunder.NotEnoughFunds
}

func (c *utxoCollector) allocateUTXO(utxo *models.UserUTXO) (err error) {
	c.addToAllocated(utxo)

	err = c.increaseSize(utxo.EstimatedInputSize)
	if err != nil {
		return fmt.Errorf("failed to increase tx size: %w", err)
	}

	err = c.increaseValue(utxo.Satoshis)
	if err != nil {
		return fmt.Errorf("failed to increase tx value: %w", err)
	}

	err = c.calculateChangeOutputs()
	if err != nil {
		return fmt.Errorf("failed to calculate change outputs: %w", err)
	}

	return nil
}

func (c *utxoCollector) addToAllocated(utxo *models.UserUTXO) {
	c.allocatedUTXOs = append(c.allocatedUTXOs, &actions.UTXO{
		TxID:     utxo.TxID,
		Vout:     utxo.Vout,
		Satoshis: utxo.Satoshis,
	})
}

func (c *utxoCollector) increaseSize(size uint64) (err error) {
	c.txSize += size
	c.fee, err = c.feeCalculator.Calculate(c.txSize)
	if err != nil {
		return fmt.Errorf("failed to calculate fee: %w", err)
	}
	return nil
}

func (c *utxoCollector) increaseValue(sats uint64) error {
	satoshis, err := to.Int64FromUnsigned(sats)
	if err != nil {
		return fmt.Errorf("utxo satoshis value int64: %w", err)
	}
	c.satsCovered += satoshis
	return nil
}

func (c *utxoCollector) satsToCover() int64 {
	return c.txSats + c.fee
}

func (c *utxoCollector) change() int64 {
	return c.satsCovered - c.satsToCover()
}

func (c *utxoCollector) prepareResult() (*actions.FundingResult, error) {
	fee, err := to.UInt64(c.fee)
	if err != nil {
		return nil, fmt.Errorf("cannot convert fee to uint64: %w", err)
	}

	changeAmount, err := to.UInt64(c.change())
	if err != nil {
		return nil, fmt.Errorf("cannot convert change amount to uint64: %w", err)
	}

	// if adding change will increase the fee so there is no change anymore,
	// we're dropping the changes and passing higher fee to miner.
	if changeAmount == 0 {
		c.changeOutputsCount = 0
	}

	return &actions.FundingResult{
		AllocatedUTXOs: c.allocatedUTXOs,
		Fee:            fee,
		ChangeAmount:   changeAmount,
		ChangeCount:    c.changeOutputsCount,
	}, nil
}

func (c *utxoCollector) calculateChangeOutputs() error {
	change := c.change()
	if change <= 0 {
		return nil
	}

	c.calculateChangeCount(must.ConvertToUInt64(change))

	err := c.increaseSize(c.changeOutputsCount * changeOutputSize)
	if err != nil {
		return fmt.Errorf("failed to increase transaction size: %w", err)
	}

	return nil
}

func (c *utxoCollector) calculateChangeCount(changeVal uint64) {
	c.changeOutputsCount = changeVal/c.minimumDesiredUTXOValue + 1

	if changeVal%c.minimumDesiredUTXOValue < c.minimumChange {
		c.changeOutputsCount -= 1
	}

	c.changeOutputsCount = to.ValueBetween(c.changeOutputsCount, 1, c.numberOfDesiredUTXOs)
}
