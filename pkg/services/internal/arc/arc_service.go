package arc

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/logging"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/services/internal/httpx"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/services/results"
	"github.com/bsv-blockchain/go-sdk/transaction"
	"github.com/go-resty/resty/v2"
	"github.com/go-softwarelab/common/pkg/is"
	"github.com/go-softwarelab/common/pkg/seq"
	"github.com/go-softwarelab/common/pkg/seq2"
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
	}

	return service
}

// PostBeef attempts to post beef with given txIDs
func (s *Service) PostBeef(ctx context.Context, beef *transaction.Beef, txids []string) (*results.PostBEEF, error) {
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

	return &results.PostBEEF{
		TxIDResults: []results.PostTxID{
			{
				TxID: res.TxID,
			},
		},
	}, nil

	// TODO:
	// for each txids
	// 	if txid == broadcasted tx.ID
	// 		get tx

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

func (s *Service) broadcast(ctx context.Context, txHex string) (*TXInfo, error) {
	result := &TXInfo{}
	arcErr := &APIError{}

	headers := httpx.NewHeaders().
		Set("X-CallbackUrl").IfNotEmpty(s.config.CallbackURL).
		Set("X-CallbackToken").IfNotEmpty(s.config.CallbackToken).
		Set("X-WaitFor").IfNotEmpty(s.config.WaitFor)

	req := s.httpClient.R().
		SetContext(ctx).
		SetHeaders(headers).
		SetResult(result).
		SetError(arcErr)

	req.SetBody(requestBody{
		RawTx: txHex,
	})

	response, err := req.Post(s.broadcastURL)

	if err != nil {
		var netError net.Error
		if errors.As(err, &netError) {
			return nil, fmt.Errorf("arc is unreachable: %w", netError)
		}
		return nil, fmt.Errorf("failed to send request to arc: %w", err)
	}

	switch response.StatusCode() {
	case http.StatusOK:
		if result.TXStatus.IsProblematic() {
			return nil, fmt.Errorf("%w: tx status: %s", ErrProblematicStatus, result.TXStatus)
		}
		return result, nil
	case http.StatusUnauthorized, http.StatusForbidden, http.StatusNotFound:
		return nil, fmt.Errorf("arc returned unauthorized: %w", arcErr)
	case StatusNotExtendedFormat:
		return nil, fmt.Errorf("arc expects transaction in extended format: %w", arcErr)
	case StatusFeeTooLow, StatusCumulativeFeeValidationFailed:
		return nil, fmt.Errorf("arc rejected transaction because of wrong fee: %w", arcErr)
	default:
		return nil, fmt.Errorf("arc cannot process provided transaction: %w", arcErr)
	}
}

type requestBody struct {
	// Even though the name suggests that it is a raw transaction,
	// it is actually a hex encoded transaction
	// and can be in Raw, Extended Format or BEEF format.
	RawTx string `json:"rawTx"`
}
