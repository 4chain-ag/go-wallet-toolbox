package whatsonchain

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/defs"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/logging"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	"github.com/go-resty/resty/v2"
	"github.com/go-softwarelab/common/pkg/to"
)

// bsvExchangeRateResponse is the response from WhatsOnChain for bsv exchange range
type bsvExchangeRateResponse struct {
	Time     int     `json:"time"`
	Rate     float64 `json:"rate"`
	Currency string  `json:"currency"`
}

type WhatsOnChain struct {
	httpClient *resty.Client
	url        string
	apiKey     string
	logger     *slog.Logger

	bsvExchangeRate   *wdk.BSVExchangeRate // TODO: possibly handle by some caching structure/redis
	bsvUpdateInterval *time.Duration
}

func New(httpClient *resty.Client, logger *slog.Logger, network defs.BSVNetwork, config WhatsOnChainConfiguration) *WhatsOnChain {
	if httpClient == nil {
		panic("httpClient is required")
	}

	logger = logging.Child(logger, "serviceWhatsOnChain")
	client := httpClient.Clone().
		SetRetryCount(Retries).
		SetRetryWaitTime(RetriesWaitTime).
		SetRetryMaxWaitTime(Retries * RetriesWaitTime)

	return &WhatsOnChain{
		httpClient:        client,
		apiKey:            config.ApiKey,
		url:               fmt.Sprintf("https://api.whatsonchain.com/v1/bsv/%s", network),
		logger:            logger,
		bsvExchangeRate:   config.BsvExchangeRate,
		bsvUpdateInterval: config.BsvUpdateInterval,
	}
}

func (woc *WhatsOnChain) prepareRequest(req *resty.Request) *resty.Request {
	req.SetHeader("Accept", "application/json")

	if woc.apiKey != "" {
		req.SetHeader("Authorization", woc.apiKey)
	}

	return req
}

func (woc *WhatsOnChain) UpdateBsvExchangeRate() (wdk.BSVExchangeRate, error) {
	if woc.bsvExchangeRate != nil {
		updateInterval := to.IfThen(woc.bsvUpdateInterval != nil, *woc.bsvUpdateInterval).ElseThen(DefaultBSVExchangeUpdateInterval)
		// Calculate the threshold time by subtracting updateMsecs from the current time
		thresholdTime := time.Now().Add(-updateInterval)

		// Check if the rate timestamp is newer than the threshold time
		if woc.bsvExchangeRate.Timestamp.After(thresholdTime) {
			return *woc.bsvExchangeRate, nil
		}
	}

	var exchangeRateResponse bsvExchangeRateResponse
	req := woc.httpClient.
		R().
		AddRetryCondition(func(res *resty.Response, err error) bool {
			return res.StatusCode() == http.StatusTooManyRequests
		})
	woc.prepareRequest(req)

	res, err := req.
		SetResult(&exchangeRateResponse).
		Get(fmt.Sprintf("%s/exchangerate", woc.url))
	if err != nil {
		return wdk.BSVExchangeRate{}, fmt.Errorf("failed to fetch exchange rate: %w", err)
	}

	if res.StatusCode() != http.StatusOK {
		return wdk.BSVExchangeRate{}, fmt.Errorf("failed to retrieve successful response from WOC. Actual status: %d", res.StatusCode())
	}

	if exchangeRateResponse.Currency != string(wdk.USD) {
		return wdk.BSVExchangeRate{}, fmt.Errorf("unsupported currency returned from Whats On Chain")
	}

	woc.bsvExchangeRate = &wdk.BSVExchangeRate{
		Timestamp: time.Now(),
		Base:      wdk.USD,
		Rate:      exchangeRateResponse.Rate,
	}

	return *woc.bsvExchangeRate, nil
}
