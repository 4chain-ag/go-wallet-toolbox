package actions

import (
	"context"
	"fmt"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/entity"
	"github.com/go-softwarelab/common/pkg/optional"
	"iter"
	"log/slog"
	"math/rand/v2"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/defs"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/logging"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/txutils"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/commission"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk/primitives"
	"github.com/go-softwarelab/common/pkg/must"
	"github.com/go-softwarelab/common/pkg/seq"
	"github.com/go-softwarelab/common/pkg/seqerr"
	"github.com/go-softwarelab/common/pkg/to"
)

const (
	derivationLength = 16
	referenceLength  = 12
)

type UTXO struct {
	TxID     string
	Vout     uint32
	Satoshis uint64
}

type FundingResult struct {
	AllocatedUTXOs []*UTXO
	ChangeCount    uint64
	ChangeAmount   uint64
	Fee            uint64
}

func (fr *FundingResult) TotalAllocated() uint64 {
	total := uint64(0)
	for _, utxo := range fr.AllocatedUTXOs {
		total += utxo.Satoshis
	}
	return total
}

type CreateActionParams struct {
	Version          int
	LockTime         int
	Description      string
	Labels           []primitives.StringUnder300
	Outputs          []wdk.ValidCreateActionOutput
	Inputs           []wdk.ValidCreateActionInput
	RandomizeOutputs bool
}

func FromValidCreateActionArgs(args *wdk.ValidCreateActionArgs) CreateActionParams {
	// TODO: use only the necessary fields (no redundant fields)
	return CreateActionParams{
		Version:          args.Version,
		LockTime:         args.LockTime,
		Description:      string(args.Description),
		Labels:           args.Labels,
		Outputs:          args.Outputs,
		Inputs:           args.Inputs,
		RandomizeOutputs: args.Options.RandomizeOutputs,
	}
}

type Funder interface {
	// Fund
	// @param targetSat - the target amount of satoshis to fund (total inputs - total outputs)
	// @param currentTxSize - the current size of the transaction in bytes (size of tx + current inputs + current outputs)
	// @param numberOfDesiredUTXOs - the number of UTXOs in basket #TakeFromBasket
	// @param minimumDesiredUTXOValue - the minimum value of UTXO in basket #TakeFromBasket
	// @param userID - the user ID
	Fund(ctx context.Context, targetSat int64, currentTxSize uint64, basket *wdk.TableOutputBasket, userID int) (*FundingResult, error)
}

type BasketRepo interface {
	FindByName(ctx context.Context, userID int, name string) (*wdk.TableOutputBasket, error)
}

type TxRepo interface {
	CreateTransaction(ctx context.Context, transaction *entity.NewTx) error
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

	xoutputs := seq.PointersFromSlice(params.Outputs)
	xinputs := seq.PointersFromSlice(params.Inputs)

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

	changeDistribution := txutils.NewChangeDistribution(basket.MinimumDesiredUTXOValue, txutils.Rand).
		Distribute(funding.ChangeCount, funding.ChangeAmount)

	derivationPrefix, reference, err := c.randomValues()
	if err != nil {
		return nil, err
	}

	newOutputs, err := c.newOutputs(
		changeDistribution,
		funding.ChangeCount,
		derivationPrefix,
		params.Outputs,
		commOut,
		params.RandomizeOutputs,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create new outputs: %w", err)
	}

	err = c.txRepo.CreateTransaction(ctx, &entity.NewTx{
		UserID:      userID,
		Version:     params.Version,
		LockTime:    params.LockTime,
		Status:      wdk.TxStatusUnsigned,
		Reference:   reference,
		IsOutgoing:  true,
		Description: params.Description,
		Satoshis:    must.ConvertToInt64FromUnsigned(funding.ChangeAmount) - must.ConvertToInt64FromUnsigned(funding.TotalAllocated()),
		Outputs:     newOutputs,
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
		Outputs:          c.resultOutputs(newOutputs),
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

func (c *create) targetSat(_ iter.Seq[*wdk.ValidCreateActionInput], xoutputs iter.Seq[*wdk.ValidCreateActionOutput]) (int64, error) {
	providedInputs := int64(0)
	// TODO: sum provided inputs satoshis - but first the values should be found

	providedOutputs := int64(0)
	for output := range xoutputs {
		satInt64, err := to.Int64FromUnsigned(output.Satoshis)
		if err != nil {
			return 0, fmt.Errorf("failed to convert satoshis to int64: %w", err)
		}
		providedOutputs += satInt64
	}

	return providedOutputs - providedInputs, nil
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
	derivationPrefix, err = randomDerivation()
	if err != nil {
		return
	}

	reference, err = txutils.RandomBase64(referenceLength)
	if err != nil {
		err = fmt.Errorf("failed to generate random reference: %w", err)
		return
	}

	return
}

func (c *create) newOutputs(
	changeDistribution iter.Seq[uint64],
	changeCount uint64,
	derivationPrefix string,
	providedOutputs []wdk.ValidCreateActionOutput,
	commissionOutput *serviceChargeOutput,
	randomizeOutputs bool,
) ([]*entity.NewOutput, error) {
	length := must.ConvertToIntFromUnsigned(changeCount) + len(providedOutputs)
	if commissionOutput != nil {
		length++
	}
	len32 := must.ConvertToUInt32(length)

	all := make([]*entity.NewOutput, 0, len32)

	for satoshis := range changeDistribution {
		derivationSuffix, err := randomDerivation()
		if err != nil {
			return nil, fmt.Errorf("failed to generate random derivation suffix: %w", err)
		}

		all = append(all, &entity.NewOutput{
			Satoshis:         must.ConvertToInt64FromUnsigned(satoshis),
			Basket:           to.Ptr(wdk.BasketNameForChange),
			Spendable:        true,
			Change:           true,
			ProvidedBy:       wdk.ProvidedByStorage,
			Type:             wdk.OutputTypeP2PKH,
			DerivationPrefix: to.Ptr(derivationPrefix),
			DerivationSuffix: to.Ptr(derivationSuffix),
		})
	}

	for _, output := range providedOutputs {
		all = append(all, &entity.NewOutput{
			Satoshis:           must.ConvertToInt64FromUnsigned(output.Satoshis),
			Basket:             (*string)(output.Basket),
			Spendable:          true,
			Change:             false,
			ProvidedBy:         wdk.ProvidedByYou,
			Type:               wdk.OutputTypeCustom,
			LockingScript:      &output.LockingScript,
			CustomInstructions: output.CustomInstructions,
			Description:        string(output.OutputDescription),
		})
	}

	if commissionOutput != nil {
		all = append(all, &entity.NewOutput{
			LockingScript: to.Ptr(commissionOutput.LockingScript),
			Satoshis:      must.ConvertToInt64FromUnsigned(commissionOutput.Satoshis),
			Basket:        nil,
			Spendable:     false,
			Change:        false,
			ProvidedBy:    wdk.ProvidedByStorage,
			Type:          wdk.OutputTypeCustom,
			Purpose:       wdk.StorageCommissionPurpose,
		})
	}

	if randomizeOutputs {
		rand.Shuffle(len(all), func(i, j int) {
			all[i], all[j] = all[j], all[i]
		})
	}

	for vout := uint32(0); vout < len32; vout++ {
		all[vout].Vout = vout
	}

	return all, nil
}

func (c *create) resultOutputs(newOutputs []*entity.NewOutput) []wdk.StorageCreateTransactionSdkOutput {
	resultOutputs := make([]wdk.StorageCreateTransactionSdkOutput, len(newOutputs))
	for i, output := range newOutputs {

		resultOutputs[i] = wdk.StorageCreateTransactionSdkOutput{
			Vout:             output.Vout,
			ProvidedBy:       output.ProvidedBy,
			Purpose:          output.Purpose,
			DerivationSuffix: output.DerivationSuffix,
			ValidCreateActionOutput: wdk.ValidCreateActionOutput{
				Satoshis:           primitives.SatoshiValue(must.ConvertToUInt64(output.Satoshis)),
				OutputDescription:  primitives.String5to2000Bytes(output.Description),
				CustomInstructions: output.CustomInstructions,
				LockingScript:      optional.OfPtr(output.LockingScript).OrZeroValue(),
			},
		}

		if output.Basket != nil {
			resultOutputs[i].Basket = to.Ptr(primitives.StringUnder300(*output.Basket))
		}
	}

	return resultOutputs
}

func randomDerivation() (string, error) {
	suffix, err := txutils.RandomBase64(derivationLength)
	if err != nil {
		return "", fmt.Errorf("failed to generate random derivation: %w", err)
	}

	return suffix, nil
}
