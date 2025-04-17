package funder

import (
	"context"
	"fmt"
	"iter"
	"log/slog"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/defs"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/logging"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/satoshi"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/txutils"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/actions"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/actions/funder/errfunder"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/database/models"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/paging"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	"github.com/go-softwarelab/common/pkg/must"
	"github.com/go-softwarelab/common/pkg/seqerr"
	"github.com/go-softwarelab/common/pkg/to"
)

var changeOutputSize = txutils.P2PKHOutputSize

const utxoBatchSize = 1000

type UTXORepository interface {
	FindFreeUTXOs(ctx context.Context, userID int, basketID int, page *paging.Page) ([]*models.UserUTXO, error)
	CountUTXOs(ctx context.Context, userID int, basketID int) (int64, error)
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
func (f *SQL) Fund(ctx context.Context, targetSat satoshi.Value, currentTxSize uint64, basket *wdk.TableOutputBasket, userID int) (*actions.FundingResult, error) {
	existing, err := f.utxoRepository.CountUTXOs(ctx, userID, basket.BasketID)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate desired utxo number in basket: %w", err)
	}

	collector, err := newCollector(targetSat, currentTxSize, basket.NumberOfDesiredUTXOs-existing, basket.MinimumDesiredUTXOValue, f.feeCalculator)
	if err != nil {
		return nil, fmt.Errorf("failed to start collecting utxo: %w", err)
	}

	utxos := f.loadUTXOs(ctx, userID, basket.BasketID)

	err = collector.Allocate(utxos)
	if err != nil {
		return nil, fmt.Errorf("failed to allocate utxos: %w", err)
	}

	return collector.GetResult()
}

func (f *SQL) loadUTXOs(ctx context.Context, userID int, basketID int) iter.Seq2[*models.UserUTXO, error] {
	batches := seqerr.ProduceWithArg(
		func(page *paging.Page) ([]*models.UserUTXO, *paging.Page, error) {
			utxos, err := f.utxoRepository.FindFreeUTXOs(ctx, userID, basketID, page)
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
	txSats satoshi.Value
	txSize uint64

	fee           satoshi.Value
	feeCalculator *feeCalc

	satsCovered    satoshi.Value
	allocatedUTXOs []*actions.UTXO

	numberOfDesiredUTXOs    uint64
	minimumDesiredUTXOValue uint64
	changeOutputsCount      uint64
	minimumChange           uint64
}

func newCollector(txSats satoshi.Value, txSize uint64, numberOfDesiredUTXOs int64, minimumDesiredUTXOValue uint64, feeCalculator *feeCalc) (c *utxoCollector, err error) {
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

	c.calculateMinimumChange()

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

	err = c.increaseValue(satoshi.MustFrom(utxo.Satoshis))
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
		Satoshis: satoshi.MustFrom(utxo.Satoshis),
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

func (c *utxoCollector) increaseValue(sats satoshi.Value) error {
	var err error
	c.satsCovered, err = satoshi.Add(c.satsCovered, sats)
	if err != nil {
		return fmt.Errorf("cannot increase tx value: %w", err)
	}
	return nil
}

func (c *utxoCollector) satsToCover() satoshi.Value {
	return satoshi.MustAdd(c.txSats, c.fee)
}

func (c *utxoCollector) change() satoshi.Value {
	return satoshi.MustSubtract(c.satsCovered, c.satsToCover())
}

func (c *utxoCollector) prepareResult() (*actions.FundingResult, error) {
	changeAmount := c.change()

	// If adding a change output increases the fee to the point where no change remains,
	// the change outputs are discarded, and the additional amount is given as a higher fee to the miner.
	if changeAmount == 0 {
		c.changeOutputsCount = 0
	}

	return &actions.FundingResult{
		AllocatedUTXOs: c.allocatedUTXOs,
		Fee:            c.fee,
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

// calculateMinimumChange determines the minimum change amount based on the **Desired** minimum UTXO value.
// The "desired" minimum UTXO value represents the user's preference for common UTXO values in the basket.
// In contrast, the minimum change is the threshold below which a new UTXO is not created.
func (c *utxoCollector) calculateMinimumChange() {
	c.minimumChange = c.minimumDesiredUTXOValue / 4
}
