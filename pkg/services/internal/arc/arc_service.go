package arc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/logging"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/services/internal"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/services/internal/httpx"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/services/results"
	"github.com/bsv-blockchain/go-sdk/transaction"
	"github.com/go-resty/resty/v2"
	"github.com/go-softwarelab/common/pkg/is"
	"github.com/go-softwarelab/common/pkg/seq"
	"github.com/go-softwarelab/common/pkg/seq2"
	"github.com/go-softwarelab/common/pkg/types"
)

// Custom ARC defined http status codes
const (
	StatusNotExtendedFormat             = 460
	StatusFeeTooLow                     = 465
	StatusCumulativeFeeValidationFailed = 473
)

var ErrProblematicStatus = errors.New("arc returned problematic status")

type Service struct {
	logger       *slog.Logger
	httpClient   *resty.Client
	config       Config
	broadcastURL string
	queryTxURL   string
}

// NewARCService creates a new arc service.
func NewARCService(logger *slog.Logger, httpClient *resty.Client, config Config) *Service {
	if httpClient == nil {
		httpClient = resty.New()
	}

	headers := httpx.NewHeaders().
		AcceptJSON().
		ContentTypeJSON().
		UserAgent().Value("go-wallet-toolbox").
		Authorization().IfNotEmpty(config.Token).
		Set("XDeployment-ID").OrDefault(config.DeploymentID, "go-wallet-toolbox#"+time.Now().Format("20060102150405"))

	httpClient = httpClient.Clone().
		SetHeaders(headers)

	service := &Service{
		logger:       logging.Child(logger, "arc"),
		httpClient:   httpClient,
		config:       config,
		broadcastURL: config.URL + "/v1/tx",
		queryTxURL:   config.URL + "/v1/tx/{txID}",
	}

	return service
}

// PostBeef attempts to post beef with given txIDs
func (s *Service) PostBeef(ctx context.Context, beef *transaction.Beef, txIDs []string) (*results.PostBEEF, error) {
	beefTxs := seq2.Values(seq2.FromMap(beef.Transactions))
	canBeSerializedToBEEFV1 := seq.Every(beefTxs, func(tx *transaction.BeefTx) bool {
		return tx.DataFormat != transaction.TxIDOnly && tx.Transaction != nil
	})

	if !canBeSerializedToBEEFV1 {
		return nil, fmt.Errorf("arc is not supporting beef v2 and provided beef cannot be converted to v1")
	}

	beefHex, err := toHex(beef)
	if err != nil {
		return nil, err
	}

	res, err := s.broadcast(ctx, beefHex)
	if err != nil {
		return nil, fmt.Errorf("failed to broadcast beef: %w", err)
	}

	txIDsToGetStatus := seq.Filter(seq.FromSlice(txIDs), func(txID string) bool {
		return res.TxID != txID
	})

	txsData := internal.MapParallel(ctx, txIDsToGetStatus, s.getTransactionData)

	txsData = seq.Prepend(txsData, internal.NewNamedResult(res.TxID, types.SuccessResult(res)))

	resultsForTxID := seq.Map(txsData, func(it *internal.NamedResult[*TXInfo]) results.PostTxID {
		if it.IsError() {
			return results.PostTxID{
				TxID:   it.Name(),
				Result: results.ResultStatusError,
				Error:  it.MustGetError(),
			}
		}
		info := it.MustGetValue()

		result := results.PostTxID{
			TxID:         it.Name(),
			DoubleSpend:  info.TXStatus == DoubleSpendAttempted,
			BlockHash:    info.BlockHash,
			BlockHeight:  info.BlockHeight,
			CompetingTxs: info.CompetingTxs,
			Data:         info,
		}

		if is.NotBlankString(info.MerklePath) {
			result.MerklePath, err = transaction.NewMerklePathFromHex(info.MerklePath)
			if err != nil {
				result.Error = err
				result.Result = results.ResultStatusError
			}
		}

		dataBytes, err := json.Marshal(info)
		if err != nil {
			result.Data = fmt.Sprintf("%+v", info)
		}
		result.Data = string(dataBytes)

		return result
	})

	return &results.PostBEEF{
		TxIDResults: seq.Collect(resultsForTxID),
	}, nil

}

func (s *Service) getTransactionData(ctx context.Context, txID string) *internal.NamedResult[*TXInfo] {
	txInfo, err := s.queryTransaction(ctx, txID)
	if err != nil {
		return internal.NewNamedResult(txID, types.FailureResult[*TXInfo](fmt.Errorf("arc query tx %s failed: %w", txID, err)))
	}

	if txInfo == nil {
		return internal.NewNamedResult(txID, types.FailureResult[*TXInfo](fmt.Errorf("not found tx %s in arc", txID)))
	}

	if txInfo.TxID != txID {
		return internal.NewNamedResult(txID, types.FailureResult[*TXInfo](fmt.Errorf("got response for tx %s while querying for %s", txInfo.TxID, txID)))
	}

	return internal.NewNamedResult(txID, types.SuccessResult(txInfo))
}

func toHex(beef *transaction.Beef) (string, error) {
	// This is a temporary solution until go-sdk properly implements BEEF serialization
	// It searches for the subject transaction in transaction.Beef and serializes this one to BEEF hex.
	// For now, it's not supporting more than one subject transaction.
	idToTx := seq2.FromMap(beef.Transactions)

	// inDegree will contain the number of transactions for which the given tx is a parent
	inDegree := seq2.CollectToMap(seq2.MapValues(idToTx, func(tx *transaction.BeefTx) int { return 0 }))

	// txsNotMined we are not interested in inputs of mined transactions
	txsNotMined := seq.Filter(seq2.Values(idToTx), func(tx *transaction.BeefTx) bool {
		return tx.Transaction.MerklePath == nil
	})

	inputs := seq.FlattenSlices(seq.Map(txsNotMined, func(tx *transaction.BeefTx) []*transaction.TransactionInput {
		return tx.Transaction.Inputs
	}))

	inputsIds := seq.Map(inputs, func(input *transaction.TransactionInput) string {
		return input.SourceTXID.String()
	})

	seq.ForEach(inputsIds, func(inputTxID string) {
		if _, ok := inDegree[inputTxID]; !ok {
			panic(fmt.Sprintf("unexpected input txid %s, this shouldn't ever happen", inputTxID))
		}
		inDegree[inputTxID]++
	})

	txIDsWithoutChildren := seq2.FilterByValue(seq2.FromMap(inDegree), is.Zero)

	subjectTxs := seq.Collect(seq2.Keys(txIDsWithoutChildren))
	if len(subjectTxs) != 1 {
		return "", fmt.Errorf("expected only one subject tx, but got %d", len(subjectTxs))
	}

	subjectTx, ok := beef.Transactions[subjectTxs[0]]
	if !ok {
		return "", fmt.Errorf("expected to find subject tx %s in beef, but it was not found, this shouldn't ever happen", subjectTxs[0])
	}

	beefHex, err := subjectTx.Transaction.BEEFHex()
	if err != nil {
		return "", fmt.Errorf("failed to convert subject tx into BEEF hex: %w", err)
	}
	return beefHex, nil
}
