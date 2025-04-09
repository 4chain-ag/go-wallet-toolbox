package services_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/services"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/services/internal/providers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mock data structures for testing
type MockRoundTripper struct {
	mockResponse *http.Response
	mockError    error
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.mockResponse, m.mockError
}

func newMockHTTPClient(response *http.Response, err error) *http.Client {
	return &http.Client{
		Transport: &MockRoundTripper{
			mockResponse: response,
			mockError:    err,
		},
	}
}

func TestUpdateBsvExchangeRateSuccess(t *testing.T) {
	t.Run("returns cached exchange rate if within update threshold", func(t *testing.T) {
		// given:
		httpClient := newMockHTTPClient(nil, nil)
		cachedRate := providers.BSVExchangeRate{
			Timestamp: time.Now().Add(-5 * time.Minute),
			Base:      "USD",
			Rate:      100.0,
		}

		// and:
		woc := services.New(httpClient, "test", services.WithBsvExchangeRate(cachedRate))

		// when:
		result, err := woc.BsvExchangeRate()

		// then:
		assert.NoError(t, err)
		assert.Equal(t, cachedRate.Rate, result)
	})

	t.Run("returns updated exchange rate when outside threshold", func(t *testing.T) {
		// given:
		responseBody := providers.BSVExchangeRateResponse{
			Time:     123456,
			Rate:     50.5,
			Currency: "USD",
		}
		bodyBytes, err := json.Marshal(responseBody)
		// then:
		require.NoError(t, err)

		// and given:
		httpClient := newMockHTTPClient(&http.Response{
			Status:     "200 OK",
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(bodyBytes)),
		}, nil)
		woc := services.New(httpClient, "test",
			services.WithBsvExchangeRate(providers.BSVExchangeRate{
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
		httpClient := newMockHTTPClient(nil, assert.AnError)

		woc := services.New(httpClient, "test")

		_, err := woc.BsvExchangeRate()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to fetch exchange rate")
	})

	t.Run("returns error if HTTP response is not 200", func(t *testing.T) {
		httpClient := newMockHTTPClient(&http.Response{
			Status:     "500 Internal Server Error",
			StatusCode: http.StatusInternalServerError,
			Body:       io.NopCloser(bytes.NewBufferString("")),
		}, nil)

		woc := services.New(httpClient, "test")

		_, err := woc.BsvExchangeRate()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to retrieve successful response from WOC")
	})

	t.Run("returns error if response JSON is invalid", func(t *testing.T) {
		httpClient := newMockHTTPClient(&http.Response{
			Status:     "200 OK",
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString("not a json")),
		}, nil)

		woc := services.New(httpClient, "test")

		_, err := woc.BsvExchangeRate()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to decode exchange rate response")
	})

	t.Run("returns error if currency is not USD", func(t *testing.T) {
		response := providers.BSVExchangeRateResponse{
			Time:     123456,
			Rate:     50.5,
			Currency: "EUR", // Not USD
		}
		bodyBytes, _ := json.Marshal(response)

		httpClient := newMockHTTPClient(&http.Response{
			Status:     "200 OK",
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(bodyBytes)),
		}, nil)

		woc := services.New(httpClient, "test")

		_, err := woc.BsvExchangeRate()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported currency")
	})
}
