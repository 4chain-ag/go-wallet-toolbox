package actions

import (
	"context"
	"fmt"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/satoshi"
	"iter"
	"log/slog"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/defs"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/logging"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/txutils"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/commission"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk/primitives"
	"github.com/go-softwarelab/common/pkg/seq"
	"github.com/go-softwarelab/common/pkg/seqerr"
)

const (
	derivationPrefixLength = 16
	referenceLength        = 12
)

type UTXO struct {
	TxID     string
	Vout     uint32
	Satoshis satoshi.Value
}

type FundingResult struct {
	AllocatedUTXOs []*UTXO
	ChangeCount    uint64
	ChangeAmount   satoshi.Value
	Fee            satoshi.Value
}

func (fr *FundingResult) TotalAllocated() (satoshi.Value, error) {
	total, err := satoshi.Sum(seq.Map(seq.FromSlice(fr.AllocatedUTXOs), func(utxo *UTXO) satoshi.Value {
		return utxo.Satoshis
	}))
	if err != nil {
		return 0, fmt.Errorf("failed to sum allocated UTXOs: %w", err)
	}

	return total, nil
}

type CreateActionParams struct {
	Version     int
	LockTime    int
	Description string
	Labels      []primitives.StringUnder300
	Outputs     iter.Seq[*wdk.ValidCreateActionOutput]
	Inputs      iter.Seq[*wdk.ValidCreateActionInput]
}

func FromValidCreateActionArgs(args *wdk.ValidCreateActionArgs) CreateActionParams {
	// TODO: use only the necessary fields (no redundant fields)
	return CreateActionParams{
		Version:     args.Version,
		LockTime:    args.LockTime,
		Description: string(args.Description),
		Labels:      args.Labels,
		Outputs:     seq.PointersFromSlice(args.Outputs),
		Inputs:      seq.PointersFromSlice(args.Inputs),
	}
}

type Funder interface {
	// Fund
	// @param targetSat - the target amount of satoshis to fund (total inputs - total outputs)
	// @param currentTxSize - the current size of the transaction in bytes (size of tx + current inputs + current outputs)
	// @param numberOfDesiredUTXOs - the number of UTXOs in basket #TakeFromBasket
	// @param minimumDesiredUTXOValue - the minimum value of UTXO in basket #TakeFromBasket
	// @param userID - the user ID
	Fund(ctx context.Context, targetSat satoshi.Value, currentTxSize uint64, basket *wdk.TableOutputBasket, userID int) (*FundingResult, error)
}

type BasketRepo interface {
	FindByName(ctx context.Context, userID int, name string) (*wdk.TableOutputBasket, error)
}

type TxRepo interface {
	CreateTransaction(ctx context.Context, transaction *wdk.NewTx) error
}

type create struct {
	logger        *slog.Logger
	funder        Funder
	basketRepo    BasketRepo
	txRepo        TxRepo
	commission    *commission.ScriptGenerator
	commissionCfg defs.Commission
}

func newCreateAction(logger *slog.Logger, funder Funder, commissionCfg defs.Commission, basketRepo BasketRepo, txRepo TxRepo) *create {
	logger = logging.Child(logger, "createAction")
	c := &create{
		logger:        logger,
		funder:        funder,
		basketRepo:    basketRepo,
		txRepo:        txRepo,
		commissionCfg: commissionCfg,
	}

	if commissionCfg.Enabled() {
		c.commission = commission.NewScriptGenerator(string(commissionCfg.PubKeyHex))
	}

	return c
}

func (c *create) Create(ctx context.Context, userID int, params CreateActionParams) (*wdk.StorageCreateActionResult, error) {
	basket, err := c.basketRepo.FindByName(ctx, userID, wdk.BasketNameForChange)
	if err != nil {
		return nil, fmt.Errorf("failed to find basket: %w", err)
	}
	if basket == nil {
		return nil, fmt.Errorf("basket for change (%s) not found", wdk.BasketNameForChange)
	}

	xoutputs := params.Outputs
	xinputs := params.Inputs

	var commOut *serviceChargeOutput
	if c.commission != nil {
		commOut, err = c.createCommissionOutput()
		if err != nil {
			return nil, fmt.Errorf("failed to collect outputs: %w", err)
		}
		xoutputs = seq.Append(xoutputs, &commOut.ValidCreateActionOutput)
	}

	initialTxSize, err := c.txSize(xinputs, xoutputs)
	if err != nil {
		return nil, err
	}

	targetSat, err := c.targetSat(xinputs, xoutputs)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate target satoshis: %w", err)
	}

	funding, err := c.funder.Fund(ctx, targetSat, initialTxSize, basket, userID)
	if err != nil {
		return nil, fmt.Errorf("funding failed: %w", err)
	}

	changeDist := txutils.NewChangeDistribution(satoshi.MustFrom(basket.MinimumDesiredUTXOValue), txutils.Rand).
		Distribute(funding.ChangeCount, funding.ChangeAmount)

	// TODO: convert change values into outputs
	_ = changeDist

	derivationPrefix, reference, err := c.randomValues()
	if err != nil {
		return nil, err
	}

	totalAllocated, err := funding.TotalAllocated()
	if err != nil {
		return nil, err
	}

	err = c.txRepo.CreateTransaction(ctx, &wdk.NewTx{
		UserID:      userID,
		Version:     params.Version,
		LockTime:    params.LockTime,
		Status:      wdk.TxStatusUnsigned,
		Reference:   reference,
		IsOutgoing:  true,
		Description: params.Description,
		Satoshis:    satoshi.MustSubtract(funding.ChangeAmount, totalAllocated).Int64(),
		Labels:      params.Labels,

		// TODO: inputBEEF
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	return &wdk.StorageCreateActionResult{
		Reference:        reference,
		Version:          params.Version,
		LockTime:         params.LockTime,
		DerivationPrefix: derivationPrefix,
	}, nil
}

type serviceChargeOutput struct {
	wdk.ValidCreateActionOutput
	KeyOffset string
}

func (c *create) createCommissionOutput() (*serviceChargeOutput, error) {
	lockingScript, keyOffset, err := c.commission.Generate()
	if err != nil {
		return nil, fmt.Errorf("failed to generate commission script: %w", err)
	}

	return &serviceChargeOutput{
		ValidCreateActionOutput: wdk.ValidCreateActionOutput{
			LockingScript:     primitives.HexString(lockingScript),
			Satoshis:          primitives.SatoshiValue(c.commissionCfg.Satoshis),
			OutputDescription: "Storage Service Charge",
		},
		KeyOffset: keyOffset,
	}, nil
}

func (c *create) targetSat(_ iter.Seq[*wdk.ValidCreateActionInput], xoutputs iter.Seq[*wdk.ValidCreateActionOutput]) (satoshi.Value, error) {
	providedInputs := satoshi.Zero()
	// TODO: sum provided inputs satoshis - but first the values should be found

	providedOutputs, err := satoshi.Sum(seq.Map(xoutputs, func(output *wdk.ValidCreateActionOutput) primitives.SatoshiValue {
		return output.Satoshis
	}))
	if err != nil {
		return 0, fmt.Errorf("failed to sum provided outputs' satoshis: %w", err)
	}

	sub, err := satoshi.Subtract(providedOutputs, providedInputs)
	if err != nil {
		return 0, fmt.Errorf("failed to subtract commission from provided outputs: %w", err)
	}

	return sub, nil
}

func (c *create) txSize(xinputs iter.Seq[*wdk.ValidCreateActionInput], xoutputs iter.Seq[*wdk.ValidCreateActionOutput]) (uint64, error) {
	inputSizes := seqerr.MapSeq(xinputs, func(o *wdk.ValidCreateActionInput) (uint64, error) {
		return o.ScriptLength()
	})

	outputSizes := seqerr.MapSeq(xoutputs, func(o *wdk.ValidCreateActionOutput) (uint64, error) {
		return o.ScriptLength()
	})

	txSize, err := txutils.TransactionSize(inputSizes, outputSizes)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate transaction size: %w", err)
	}

	return txSize, nil
}

func (c *create) randomValues() (derivationPrefix string, reference string, err error) {
	derivationPrefix, err = txutils.RandomBase64(derivationPrefixLength)
	if err != nil {
		err = fmt.Errorf("failed to generate random derivation prefix: %w", err)
		return
	}

	reference, err = txutils.RandomBase64(referenceLength)
	if err != nil {
		err = fmt.Errorf("failed to generate random reference: %w", err)
		return
	}

	return
}
