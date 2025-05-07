package testabilities

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	sdk "github.com/bsv-blockchain/go-sdk/transaction"
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const ArcURL = "https://api.taal.com/arc"
const ArcToken = "mainnet_9596de07e92300c6287e4393594ae39c"
const DeploymentID = "go-wallet-toolbox-test"

type ArcFixture interface {
	IsUpAndRunning()
	HttpClient() *resty.Client
}

type arcFixture struct {
	testing.TB
	transport *httpmock.MockTransport
}

func NewArcFixture(t testing.TB) ArcFixture {
	transport := httpmock.NewMockTransport()
	return NewArcFixtureWithTransport(t, transport)
}

func NewArcFixtureWithTransport(t testing.TB, transport *httpmock.MockTransport) ArcFixture {
	require.NotNil(t, transport, "http.RoundTripper must be provided")

	return &arcFixture{
		TB:        t,
		transport: transport,
	}
}

func (f *arcFixture) HttpClient() *resty.Client {
	client := resty.New()
	client.SetTransport(f.transport)
	return client
}

func (f *arcFixture) IsUpAndRunning() {
	f.transport.RegisterResponder("POST", ArcURL+"/v1/tx", func(req *http.Request) (*http.Response, error) {
		b, err := io.ReadAll(req.Body)
		if !assert.NoError(f, err) {
			return nil, err
		}

		var body map[string]any
		err = json.Unmarshal(b, &body)
		if !assert.NoError(f, err) {
			return nil, err
		}

		rawTx := body["rawTx"]
		if !assert.NotNil(f, rawTx) {
			return httpmock.NewJsonResponse(400, map[string]any{
				"detail":    "The request seems to be malformed and cannot be processed",
				"extraInfo": "error parsing transactions from request: no transaction found - empty request body",
				"instance":  nil,
				"status":    400,
				"title":     "Bad request",
				"txid":      nil,
				"type":      "https://bitcoin-sv.github.io/arc/#/errors?id=_400",
			})
		}
		txHex := rawTx.(string)

		tx, err := sdk.NewTransactionFromBEEFHex(txHex)
		if !assert.NoError(f, err) {
			return httpmock.NewJsonResponse(463, map[string]any{
				"detail":    "Transaction is malformed and cannot be processed",
				"extraInfo": err.Error(),
				"instance":  nil,
				"status":    463,
				"title":     "Malformed transaction",
				"txid":      nil,
				"type":      "https://bitcoin-sv.github.io/arc/#/errors?id=_463",
			})
		}

		id := tx.TxID()

		content := map[string]any{
			"blockHash":    "",
			"blockHeight":  0,
			"competingTxs": nil,
			"extraInfo":    "",
			"merklePath":   "",
			"timestamp":    time.Now().Format("2006-01-02T15:04:05.999999999Z"),
			"txStatus":     "SEEN_ON_NETWORK",
			"txid":         id,
		}

		responder, err := httpmock.NewJsonResponder(200, content)
		require.NoError(f, err)

		f.transport.RegisterResponder("GET", ArcURL+"/v1/tx/"+id.String(), responder)
		return responder(req)
	})
}
