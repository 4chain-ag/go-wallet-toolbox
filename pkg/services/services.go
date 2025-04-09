package services

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/defs"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/logging"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/services/internal/providers"
	"github.com/bsv-blockchain/go-sdk/transaction"
	"github.com/bsv-blockchain/go-sdk/transaction/chaintracker"
)

// WalletServices is a struct that contains services used by a wallet
type WalletServices struct {
	httpClient   *http.Client
	logger       *slog.Logger
	chain        defs.BSVNetwork
	options      *WalletServicesOptions
	whatsonchain *providers.WhatsOnChain

	// getMerklePathServices: ServiceCollection<sdk.GetMerklePathService>
	// getRawTxServices: ServiceCollection<sdk.GetRawTxService>
	// postBeefServices: ServiceCollection<sdk.PostBeefService>
	// getUtxoStatusServices: ServiceCollection<sdk.GetUtxoStatusService>
	// updateFiatExchangeRateServices: ServiceCollection<sdk.UpdateFiatExchangeRateService>
}

// New will return a new WalletServices
func New(httpClient *http.Client, logger *slog.Logger, chain defs.BSVNetwork, opts ...Options) *WalletServices {
	if httpClient == nil {
		panic("httpClient is required")
	}

	options := defaultServicesOptions(chain)
	for _, opt := range opts {
		opt(options)
	}

	return &WalletServices{
		httpClient:   httpClient,
		chain:        chain,
		options:      options,
		logger:       logger,
		whatsonchain: providers.NewWhatsOnChain(httpClient, logging.Child(logger, "service_whats_on_chain"), options.WhatsOnChainApiKey, chain),
	}
}

// ChainTracker returns service which requires `options.chaintracks` be valid.
func (s *WalletServices) ChainTracker() chaintracker.ChainTracker {
	panic("Not implemented yet")
}

// HeaderForHeight returns serialized block header for height on active chain
func (s *WalletServices) HeaderForHeight(height int) ([]int, error) {
	panic("Not implemented yet")
}

// Height returns the height of the active chain
func (s *WalletServices) Height() int {
	panic("Not implemented yet")
}

// BsvExchangeRate returns approximate exchange rate US Dollar / BSV, USD / BSV
// This is the US Dollar price of one BSV
func (s *WalletServices) BsvExchangeRate() (float64, error) {
	bsvExchangeRate, err := s.whatsonchain.UpdateBsvExchangeRate(
		&s.options.BsvExchangeRate,
		&s.options.BsvUpdateMsecs,
	)
	if err != nil {
		return 0, fmt.Errorf("error during bsvExchangeRate: %w", err)
	}

	s.options.BsvExchangeRate = bsvExchangeRate
	return bsvExchangeRate.Rate, nil
}

// FiatExchangeRate returns  approximate exchange rate currency per base.
func (s *WalletServices) FiatExchangeRate(currency Currency, base *Currency) float64 {
	panic("Not implemented yet")
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
func (s *WalletServices) RawTx(txid string, useNext bool) (RawTxResult, error) {
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
