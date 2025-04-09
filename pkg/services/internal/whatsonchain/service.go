package whatsonchain

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/defs"
	"github.com/go-resty/resty/v2"
	"github.com/go-softwarelab/common/pkg/to"
)

type WhatsOnChain struct {
	httpClient *resty.Client
	url        string
	apiKey     string
	logger     *slog.Logger
}

func New(httpClient *resty.Client, logger *slog.Logger, apiKey string, network defs.BSVNetwork) *WhatsOnChain {
	return &WhatsOnChain{
		httpClient: httpClient,
		apiKey:     apiKey,
		url:        fmt.Sprintf("https://api.whatsonchain.com/v1/bsv/%s", network),
		logger:     logger,
	}
}

func (woc *WhatsOnChain) UpdateBsvExchangeRate(exchangeRate *BSVExchangeRate, bsvUpdateMsecs *int) (BSVExchangeRate, error) {
	if exchangeRate != nil {
		fifteenMinutesMs := 1000 * 60 * 15
		updateMscecs := to.IfThen(bsvUpdateMsecs != nil, *bsvUpdateMsecs).ElseThen(fifteenMinutesMs)
		// Calculate the threshold time by subtracting updateMsecs from the current time
		thresholdTime := time.Now().Add(-time.Duration(updateMscecs) * time.Millisecond)

		// Check if the rate timestamp is newer than the threshold time
		if exchangeRate.Timestamp.After(thresholdTime) {
			return *exchangeRate, nil
		}
	}

	for retry := 0; retry < 2; retry++ {
		var exchangeRateResponse BSVExchangeRateResponse
		req := woc.httpClient.R().
			SetHeader("Accept", "application/json")
		if woc.apiKey != "" {
			req.SetHeader("Authorization", woc.apiKey)
		}

		res, err := req.
			SetResult(&exchangeRateResponse).
			Get(fmt.Sprintf("%s/exchangerate", woc.url))
		if err != nil {
			return BSVExchangeRate{}, fmt.Errorf("failed to fetch exchange rate: %w", err)
		}

		if res.Status() == "Too Many Requests" && retry < 2 {
			time.Sleep(2000 * time.Millisecond)
			continue
		}

		if res.StatusCode() != http.StatusOK {
			return BSVExchangeRate{}, fmt.Errorf("failed to retrieve successful response from WOC. Actual status: %d", res.StatusCode())
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
