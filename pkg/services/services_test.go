package services_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/services"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/services/internal/whatsonchain"
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func newMockHTTPClient(status int, content string, err error) *resty.Client {
	transport := httpmock.NewMockTransport()
	client := resty.New()
	client.GetClient().Transport = transport

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

	// transport.RegisterResponder("GET", fmt.Sprintf("%s/v1/tx/%s", arcURL, minedTxID), responder(http.StatusOK, `{
	transport.RegisterResponder("GET", "https://api.whatsonchain.com/v1/bsv/test/exchangerate", responder(status, content))

	return client
}

func TestUpdateBsvExchangeRateSuccess(t *testing.T) {
	t.Run("returns cached exchange rate if within update threshold", func(t *testing.T) {
		// given:
		httpClient := newMockHTTPClient(200, "", nil)
		cachedRate := whatsonchain.BSVExchangeRate{
			Timestamp: time.Now().Add(-5 * time.Minute),
			Base:      "USD",
			Rate:      100.0,
		}

		// and:
		woc := services.New(httpClient, nil, "test", services.WithBsvExchangeRate(cachedRate))

		// when:
		result, err := woc.BsvExchangeRate()

		// then:
		assert.NoError(t, err)
		assert.Equal(t, cachedRate.Rate, result)
	})

	t.Run("returns updated exchange rate when outside threshold", func(t *testing.T) {
		// given:
		httpClient := newMockHTTPClient(200, `{
			"time": 123456,
			"rate": 50.5,
			"currency": "USD"
		}`, nil)

		woc := services.New(httpClient, nil, "test",
			services.WithBsvExchangeRate(whatsonchain.BSVExchangeRate{
				Timestamp: time.Now().Add(-16 * time.Minute),
				Base:      "USD",
				Rate:      100.0,
			}))

		// when:
		result, err := woc.BsvExchangeRate()

		// then:
		assert.NoError(t, err)
		assert.Equal(t, 50.5, result)
	})
}

func TestUpdateBsvExchangeRateFail(t *testing.T) {
	t.Run("returns error if HTTP request fails", func(t *testing.T) {
		httpClient := newMockHTTPClient(200, `{
			"time": 123456,
			"rate": 50.5,
			"currency": "USD"
		}`, assert.AnError)

		woc := services.New(httpClient, nil, "test")

		_, err := woc.BsvExchangeRate()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to fetch exchange rate")
	})

	t.Run("returns error if HTTP response is not 200", func(t *testing.T) {
		httpClient := newMockHTTPClient(500, `{
			"time": 123456,
			"rate": 50.5,
			"currency": "EUR"
		}`, nil)

		woc := services.New(httpClient, nil, "test")

		_, err := woc.BsvExchangeRate()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to retrieve successful response from WOC")
	})

	t.Run("returns error if currency is not USD", func(t *testing.T) {
		httpClient := newMockHTTPClient(200, `{
			"time": 123456,
			"rate": 50.5,
			"currency": "EUR"
		}`, nil)

		woc := services.New(httpClient, nil, "test")

		_, err := woc.BsvExchangeRate()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported currency")
	})
}
