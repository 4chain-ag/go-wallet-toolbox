package testabilities

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

type WhatsOnChainFixture interface {
	WillRespondWithRates(status int, content string, err error)

	WillRespondWithRawTx(status int, txID, rawTx string, err error)
}

type wocFixture struct {
	testing.TB
	transport *httpmock.MockTransport
}

func NewWoCFixture(t testing.TB) WhatsOnChainFixture {
	return NewWoCFixtureWithTransport(t, httpmock.NewMockTransport())
}

func NewWoCFixtureWithTransport(t testing.TB, transport *httpmock.MockTransport) WhatsOnChainFixture {
	require.NotNil(t, transport, "transport must be provided")
	return &wocFixture{
		TB:        t,
		transport: transport,
	}

}

func (f *wocFixture) WillRespondWithRates(status int, content string, err error) {
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

	f.transport.RegisterResponder("GET", "https://api.whatsonchain.com/v1/bsv/test/exchangerate", responder(status, content))
}

func (f *wocFixture) WillRespondWithRawTx(status int, txID, rawTx string, err error) {
	responder := func(status int, content string) func(req *http.Request) (*http.Response, error) {
		return func(req *http.Request) (*http.Response, error) {
			if err != nil {
				return nil, err
			}
			res := httpmock.NewStringResponse(status, content)
			res.Header.Set("Content-Type", "text/plain")
			return res, nil
		}
	}

	url := fmt.Sprintf("https://api.whatsonchain.com/v1/bsv/test/tx/%s/hex", txID)
	f.transport.RegisterResponder("GET", url, responder(status, rawTx))
}
