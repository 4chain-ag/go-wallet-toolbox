package testabilities

import (
	"log/slog"
	"net/http"
	"testing"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/defs"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/logging"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/services"
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

type ServicesFixture interface {
	WhatsOnChain() WhatsOnChainFixture
	NewServices(opts ...services.Options) *services.WalletServices
}

type WhatsOnChainFixture interface {
	WillRespondWithRates(status int, content string, err error) WhatsOnChainFixture
}

type servicesFixture struct {
	t          testing.TB
	require    *require.Assertions
	logger     *slog.Logger
	services   *services.WalletServices
	httpClient *resty.Client
	transport  *httpmock.MockTransport
}

func (s *servicesFixture) WhatsOnChain() WhatsOnChainFixture {
	return s
}

func (s *servicesFixture) NewServices(opts ...services.Options) *services.WalletServices {
	s.t.Helper()

	walletServices := services.New(s.httpClient, s.logger, defs.NetworkTestnet, opts...)

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

	return &servicesFixture{
		t:          t,
		require:    require.New(t),
		logger:     logging.NewTestLogger(t),
		httpClient: client,
		transport:  transport,
	}
}
