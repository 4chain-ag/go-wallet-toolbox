package whatsonchain

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/defs"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/logging"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/services/internal"
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
	logger = logging.Child(logger, "serviceWhatsOnChain")

	return &WhatsOnChain{
		httpClient: httpClient,
		apiKey:     apiKey,
		url:        fmt.Sprintf("https://api.whatsonchain.com/v1/bsv/%s", network),
		logger:     logger,
	}
}

func (woc *WhatsOnChain) prepareRequest(req *resty.Request) *resty.Request {
	req.SetHeader("Accept", "application/json")

	if woc.apiKey != "" {
		req.SetHeader("Authorization", woc.apiKey)
	}

	return req
}

func (woc *WhatsOnChain) UpdateBsvExchangeRate(exchangeRate *BSVExchangeRate, bsvUpdateDuration *time.Duration) (BSVExchangeRate, error) {
	if exchangeRate != nil {
		updateInterval := to.IfThen(bsvUpdateDuration != nil, *bsvUpdateDuration).ElseThen(internal.FifteenMinutes)
		// Calculate the threshold time by subtracting updateMsecs from the current time
		thresholdTime := time.Now().Add(-updateInterval)

		// Check if the rate timestamp is newer than the threshold time
		if exchangeRate.Timestamp.After(thresholdTime) {
			return *exchangeRate, nil
		}
	}

	var exchangeRateResponse BSVExchangeRateResponse
	req := woc.httpClient.Clone().
		SetRetryCount(internal.Retries).
		SetRetryWaitTime(internal.RetriesWaitTime).
		SetRetryMaxWaitTime(internal.Retries * internal.RetriesWaitTime).
		AddRetryCondition(func(res *resty.Response, err error) bool {
			return res.Status() == "Too Many Requests"
		}).
		R()
	woc.prepareRequest(req)

	res, err := req.
		SetResult(&exchangeRateResponse).
		Get(fmt.Sprintf("%s/exchangerate", woc.url))
	if err != nil {
		return BSVExchangeRate{}, fmt.Errorf("failed to fetch exchange rate: %w", err)
	}

	if res.StatusCode() != http.StatusOK {
		return BSVExchangeRate{}, fmt.Errorf("failed to retrieve successful response from WOC. Actual status: %d", res.StatusCode())
	}

	if exchangeRateResponse.Currency != string(internal.USD) {
		return BSVExchangeRate{}, fmt.Errorf("unsupported currency returned from Whats On Chain")
	}

	return BSVExchangeRate{
		Timestamp: time.Now(),
		Base:      internal.USD,
		Rate:      exchangeRateResponse.Rate,
	}, nil
}
