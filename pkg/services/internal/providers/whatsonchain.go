package providers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/defs"

	"github.com/go-softwarelab/common/pkg/to"
)

type WhatsOnChain struct {
	httpClient *http.Client
	url        string
	apiKey     string
}

func NewWhatsOnChain(httpClient *http.Client, apiKey string, network defs.BSVNetwork) *WhatsOnChain {
	return &WhatsOnChain{
		httpClient: httpClient,
		apiKey:     apiKey,
		url:        fmt.Sprintf("https://api.whatsonchain.com/v1/bsv/%s", network),
	}
}

func (woc *WhatsOnChain) UpdateBsvExchangeRate(exchangeRate *BSVExchangeRate, bsvUpdateMsecs *int) (BSVExchangeRate, error) {
	if exchangeRate != nil {
		fifteenMinutesMs := 1000 * 60 * 15
		updateMscecs := to.IfThen(bsvUpdateMsecs != nil, *bsvUpdateMsecs).ElseThen(fifteenMinutesMs)
		if time.Since(exchangeRate.Timestamp) < time.Duration(updateMscecs) {
			return *exchangeRate, nil
		}
	}

	for retry := 0; retry < 2; retry++ {
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/exchangerate", woc.url), nil)
		if err != nil {
			return BSVExchangeRate{}, fmt.Errorf("failed to prepare request for fetch exchange rate: %w", err)
		}
		req.Header.Add("Accept", "application/json")
		if woc.apiKey != "" {
			req.Header.Add("Authorization", woc.apiKey)
		}

		res, err := woc.httpClient.Do(req)
		if err != nil {
			return BSVExchangeRate{}, fmt.Errorf("failed to fetch exchange rate: %w", err)
		}

		if res.Status == "Too Many Requests" && retry < 2 {
			time.Sleep(2000 * time.Millisecond)
			continue
		}

		if res.StatusCode != http.StatusOK {
			return BSVExchangeRate{}, fmt.Errorf("failed to retrieve successful response from WOC. Actual status: %d", res.StatusCode)
		}

		exchangeRateResponse := &BSVExchangeRateResponse{}
		if err := json.NewDecoder(res.Body).Decode(exchangeRateResponse); err != nil {
			return BSVExchangeRate{}, fmt.Errorf("failed to decode exchange rate response: %w", err)
		}
		if exchangeRateResponse.Currency != "USD" {
			return BSVExchangeRate{}, fmt.Errorf("unsupported currency returned from Whats On Chain")
		}

		return BSVExchangeRate{
			Timestamp: time.Now(),
			Base:      "USD",
			Rate:      exchangeRateResponse.Rate,
		}, nil

	}

	return BSVExchangeRate{}, fmt.Errorf("failed to update exchange rate: internal error")
}
