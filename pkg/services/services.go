package services

import (
	"github.com/4chain-ag/go-wallet-toolbox/pkg/defs"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/services/internal/providers"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	"github.com/bsv-blockchain/go-sdk/transaction"
	"github.com/bsv-blockchain/go-sdk/transaction/chaintracker"
)

type Services struct {
	chain        defs.BSVNetwork
	options      WalletServicesOptions
	whatsonchain *providers.WhatsOnChain
	arc          *providers.ARC
	bitails      *providers.Bitails

	// getMerklePathServices: ServiceCollection<sdk.GetMerklePathService>
	// getRawTxServices: ServiceCollection<sdk.GetRawTxService>
	// postBeefServices: ServiceCollection<sdk.PostBeefService>
	// getUtxoStatusServices: ServiceCollection<sdk.GetUtxoStatusService>
	// updateFiatExchangeRateServices: ServiceCollection<sdk.UpdateFiatExchangeRateService>
}

func New(chain defs.BSVNetwork, opts ...Options) *Services {
	options := defaultServicesOptions(chain)
	for _, opt := range opts {
		opt(&options)
	}

	return &Services{
		chain:        chain,
		options:      options,
		whatsonchain: providers.NewWhatsOnChain(),
		arc:          providers.NewARC(),
		bitails:      providers.NewBitails(),
	}
}

func (s *Services) ChainTracker() chaintracker.ChainTracker {
	panic("Not implemented yet")
}

func (s *Services) HeaderForHeight(height int) ([]int, error) {
	panic("Not implemented yet")
}

func (s *Services) Height() int {
	panic("Not implemented yet")
}

func (s *Services) BsvExchangeRate() float64 {
	panic("Not implemented yet")
}

func (s *Services) FiatExchangeRate(currency wdk.Currency, base *wdk.Currency) float64 {
	panic("Not implemented yet")
}

func (s *Services) RawTx(txid string, useNext bool) (wdk.RawTxResult, error) {
	panic("Not implemented yet")
}

func (s *Services) MerklePath(txid string, useNext bool) (wdk.MerklePathResult, error) {
	panic("Not implemented yet")
}

func (s *Services) PostBeef(beef *transaction.Beef, txids []string) ([]*wdk.PostBeefResult, error) {
	panic("Not implemented yet")
}

func (s *Services) UtxoStatus(
	output string,
	outputFormat wdk.UtxoStatusOutputFormat,
	useNext bool,
) (wdk.UtxoStatusResult, error) {
	panic("Not implemented yet")
}

func (s *Services) HashToHeader(hash string) (*wdk.BlockHeader, error) {
	panic("Not implemented yet")
}

func (s *Services) NLockTimeIsFinal(txOrLockTime any) bool {
	panic("Not implemented yet")
}
