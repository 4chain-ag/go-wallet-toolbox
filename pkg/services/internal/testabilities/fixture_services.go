package testabilities

import (
	"fmt"
	"log/slog"
	"net/http"
	"testing"
	"time"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/defs"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/logging"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/services"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/services/configuration"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/services/internal/whatsonchain"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	"github.com/go-resty/resty/v2"
	"github.com/go-softwarelab/common/pkg/to"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

type ServicesFixture interface {
	WhatsOnChain() WhatsOnChainFixture

	Services() WalletServicesFixture
	NewServicesWithConfig(config configuration.WalletServices) *services.WalletServices
}

type WalletServicesFixture interface {
	WithDefaultConfig() *services.WalletServices

	WithBsvExchangeRate(exchangeRate wdk.BSVExchangeRate) *services.WalletServices
}
type WhatsOnChainFixture interface {
	WillRespondWithRates(status int, content string, err error) WhatsOnChainFixture
}

type servicesFixture struct {
	t                    testing.TB
	require              *require.Assertions
	logger               *slog.Logger
	services             *services.WalletServices
	httpClient           *resty.Client
	transport            *httpmock.MockTransport
	walletServicesConfig *configuration.WalletServices
}

func (s *servicesFixture) WhatsOnChain() WhatsOnChainFixture {
	return s
}

func (s *servicesFixture) WithDefaultConfig() *services.WalletServices {
	s.t.Helper()

	walletServices := services.New(s.httpClient, s.logger, *s.walletServicesConfig)
	s.services = walletServices

	return s.services
}

func (s *servicesFixture) WithBsvExchangeRate(exchangeRate wdk.BSVExchangeRate) *services.WalletServices {
	s.t.Helper()
	s.walletServicesConfig.WhatsOnChain.BSVExchangeRate = exchangeRate

	walletServices := services.New(s.httpClient, s.logger, *s.walletServicesConfig)
	s.services = walletServices

	return s.services
}

func (s *servicesFixture) Services() WalletServicesFixture {
	return s
}

func (s *servicesFixture) NewServicesWithConfig(config configuration.WalletServices) *services.WalletServices {
	s.t.Helper()

	walletServices := services.New(s.httpClient, s.logger, config)

	s.services = walletServices

	return s.services
}

func (s *servicesFixture) WillRespondWithRates(status int, content string, err error) WhatsOnChainFixture {
	responder := func(status int, content string) func(req *http.Request) (*http.Response, error) {
		return func(req *http.Request) (*http.Response, error) {
			if err != nil {
				return nil, err
			}
			res := httpmock.NewStringResponse(status, content)
			res.Header.Set("Content-Type", "application/json")
			return res, nil
		}
	}

	s.transport.RegisterResponder("GET", "https://api.whatsonchain.com/v1/bsv/test/exchangerate", responder(status, content))

	return s
}

func Given(t testing.TB) ServicesFixture {
	transport := httpmock.NewMockTransport()
	client := resty.New()
	client.GetClient().Transport = transport

	servicesConfig := servicesCfg(defs.NetworkTestnet)

	return &servicesFixture{
		t:                    t,
		require:              require.New(t),
		logger:               logging.NewTestLogger(t),
		httpClient:           client,
		transport:            transport,
		walletServicesConfig: &servicesConfig,
	}
}

func servicesCfg(chain defs.BSVNetwork) configuration.WalletServices {
	var taalApiKey string
	var port int
	var arcUrl string

	if chain == defs.NetworkMainnet {
		taalApiKey = "mainnet_9596de07e92300c6287e4393594ae39c"
		port = 8084
		arcUrl = "https://api.taal.com/arc"
	} else {
		taalApiKey = "testnet_0e6cf72133b43ea2d7861da2a38684e3"
		port = 8083
		arcUrl = "https://arc-test.taal.com/arc"
	}

	return configuration.WalletServices{
		Chain:      chain,
		TaalAPIKey: taalApiKey,
		WhatsOnChain: configuration.WhatsOnChain{
			BSVUpdateInterval: to.Ptr(whatsonchain.DefaultBSVExchangeUpdateInterval),
		},
		FiatExchangeRates: wdk.FiatExchangeRates{
			Timestamp: time.Date(2023, time.December, 13, 0, 0, 0, 0, time.UTC),
			Base:      wdk.USD,
			Rates: map[wdk.Currency]float64{
				wdk.USD: 1,
				wdk.GBP: 0.8,
				wdk.EUR: 0.93,
			},
		},
		FiatUpdateInterval:              to.Ptr(whatsonchain.DefaultFiatExchangeUpdateInterval),
		DisableMapiCallback:             true, // rely on WalletMonitor by default
		ExchangeratesApiKey:             "bd539d2ff492bcb5619d5f27726a766f",
		ChaintracksFiatExchangeRatesUrl: fmt.Sprintf("https://npm-registry.babbage.systems:%d/getFiatExchangeRates", port),
		Chaintracks:                     nil, // TODO: implement me
		ArcURL:                          arcUrl,
		ArcConfig:                       nil, // TODO: implement me
	}
}
