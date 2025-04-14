package services

import (
	"fmt"
	"log/slog"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/defs"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/transactions/utils"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/services/configuration"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/services/internal"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/services/internal/servicequeue"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/services/internal/whatsonchain"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	"github.com/bsv-blockchain/go-sdk/transaction"
	"github.com/bsv-blockchain/go-sdk/transaction/chaintracker"
	"github.com/go-resty/resty/v2"
)

// RawTxResultService is an interface for services that implement RawTx method
type RawTxResultService interface {
	RawTx(txID string, chain defs.BSVNetwork) (*internal.RawTxResult, error)
}

// WalletServices is a struct that contains services used by a wallet
type WalletServices struct {
	httpClient   *resty.Client
	logger       *slog.Logger
	chain        defs.BSVNetwork
	config       *configuration.WalletServices
	whatsonchain *whatsonchain.WhatsOnChain

	rawTxServices *servicequeue.ServicesQueue[RawTxResultService]
	// getMerklePathServices: ServiceCollection<sdk.GetMerklePathService>
	// getRawTxServices: ServiceCollection<sdk.GetRawTxService>
	// postBeefServices: ServiceCollection<sdk.PostBeefService>
	// getUtxoStatusServices: ServiceCollection<sdk.GetUtxoStatusService>
	// updateFiatExchangeRateServices: ServiceCollection<sdk.UpdateFiatExchangeRateService>
}

// New will return a new WalletServices
func New(httpClient *resty.Client, logger *slog.Logger, config configuration.WalletServices) *WalletServices {
	if httpClient == nil {
		panic("httpClient is required")
	}

	woc := whatsonchain.New(httpClient, logger, config.Chain, config.WhatsOnChain)

	rawTxResultServices := servicequeue.New(
		servicequeue.ServiceExecutor[RawTxResultService]{
			Name:    "WhatsOnChain",
			Service: woc,
		})

	return &WalletServices{
		httpClient:    httpClient,
		chain:         config.Chain,
		config:        &config,
		logger:        logger,
		whatsonchain:  woc,
		rawTxServices: rawTxResultServices,
	}
}

// RawTx attempts to obtain the raw transaction bytes associated with a 32 byte transaction hash (txid).
//
// Cycles through configured transaction processing services attempting to get a valid response.
//
// On success:
// Result txid is the requested transaction hash
// Result rawTx will be an array containing raw transaction bytes.
// Result name will be the responding service's identifying name.
// Returns result without incrementing active service.
//
// On failure:
// Result txid is the requested transaction hash
// Result mapi will be the first mapi response obtained (service name and response), or null
// Result error will be the first error thrown (service name and CwiError), or null
// Increments to next configured service and tries again until all services have been tried.
func (s *WalletServices) RawTx(txID string) (internal.RawTxResult, error) {
	var lastErr error

	for tries := 0; tries < s.rawTxServices.Count(); tries++ {
		srv, err := s.rawTxServices.Current()
		if err != nil {
			return internal.RawTxResult{}, fmt.Errorf("error while getting current service: %w", err)
		}

		resp, err := srv.Service.RawTx(txID, s.chain)
		if err != nil {
			lastErr = fmt.Errorf("error while getting raw tx result from service: %s => %w", srv.Name, err)
			s.logger.Error(lastErr.Error())
			s.rawTxServices.Next()
			continue
		}

		if resp == nil {
			lastErr = fmt.Errorf("transaction with txID: %s not found", txID)
			s.rawTxServices.Next()
			continue
		}

		txIDFromRawTx := utils.TransactionIDFromRawTx(resp.RawTx)
		if txID != txIDFromRawTx {
			lastErr = fmt.Errorf("computed txid %s doesn't match requested value %s", txIDFromRawTx, txID)
			s.rawTxServices.Next()
			continue
		}

		return *resp, nil
	}

	if lastErr != nil {
		return internal.RawTxResult{}, fmt.Errorf("all services failed, last error: %w", lastErr)
	}

	return internal.RawTxResult{}, fmt.Errorf("internal error during getting RawTx")
}

// ChainTracker returns service, which requires `options.chaintracks` be valid.
func (s *WalletServices) ChainTracker() chaintracker.ChainTracker {
	panic("Not implemented yet")
}

// HeaderForHeight returns serialized block header for height on active chain
func (s *WalletServices) HeaderForHeight(height int64) ([]int64, error) {
	panic("Not implemented yet")
}

// Height returns the height of the active chain
func (s *WalletServices) Height() int64 {
	panic("Not implemented yet")
}

// BsvExchangeRate returns approximate exchange rate US Dollar / BSV, USD / BSV
// This is the US Dollar price of one BSV
func (s *WalletServices) BsvExchangeRate() (float64, error) {
	bsvExchangeRate, err := s.whatsonchain.UpdateBsvExchangeRate()
	if err != nil {
		return 0, fmt.Errorf("error during bsvExchangeRate: %w", err)
	}

	return bsvExchangeRate.Rate, nil
}

// FiatExchangeRate returns approximate exchange rate currency per base.
func (s *WalletServices) FiatExchangeRate(currency wdk.Currency, base *wdk.Currency) float64 {
	panic("Not implemented yet")
}

// MerklePath attempts to obtain the merkle proof associated with a 32 byte transaction hash (txid).
//
// Cycles through configured transaction processing services attempting to get a valid response.
//
// On success:
// Result txid is the requested transaction hash
// Result proof will be the merkle proof.
// Result name will be the responding service's identifying name.
// Returns result without incrementing active service.
//
// On failure:
// Result txid is the requested transaction hash
// Result mapi will be the first mapi response obtained (service name and response), or null
// Result error will be the first error thrown (service name and CwiError), or null
// Increments to next configured service and tries again until all services have been tried.
func (s *WalletServices) MerklePath(txid string, useNext bool) (MerklePathResult, error) {
	panic("Not implemented yet")
}

// PostBeef attempts to post beef with given txIDs
func (s *WalletServices) PostBeef(beef *transaction.Beef, txids []string) ([]*PostBeefResult, error) {
	panic("Not implemented yet")
}

// UtxoStatus attempts to determine the UTXO status of a transaction output.
//
// Cycles through configured transaction processing services attempting to get a valid response.
func (s *WalletServices) UtxoStatus(
	output string,
	outputFormat UtxoStatusOutputFormat,
	useNext bool,
) (UtxoStatusResult, error) {
	panic("Not implemented yet")
}

// HashToHeader attempts to retrieve BlockHeader by its hash
func (s *WalletServices) HashToHeader(hash string) (*BlockHeader, error) {
	panic("Not implemented yet")
}

// NLockTimeIsFinal returns whether the locktime value allows the transaction to be mined at the current chain height
// TODO: txOrLockTime type = string | number[] | BsvTransaction | number
func (s *WalletServices) NLockTimeIsFinal(txOrLockTime any) bool {
	panic("Not implemented yet")
}
