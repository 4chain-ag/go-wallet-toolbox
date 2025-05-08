package testabilities

import (
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	sdk "github.com/bsv-blockchain/go-sdk/transaction"
	"github.com/go-resty/resty/v2"
	"github.com/go-softwarelab/common/pkg/optional"
	"github.com/go-softwarelab/common/pkg/seq"
	"github.com/go-softwarelab/common/pkg/seq2"
	"github.com/go-softwarelab/common/pkg/to"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const ArcURL = "https://api.taal.com/arc"
const ArcToken = "mainnet_9596de07e92300c6287e4393594ae39c"
const DeploymentID = "go-wallet-toolbox-test"

var timestamp = time.Date(2018, time.November, 10, 23, 0, 0, 0, time.UTC).Format("2006-01-02T15:04:05.999999999Z")

type ArcFixture interface {
	IsUpAndRunning()
	HttpClient() *resty.Client
	TxInfoJSON(id string) string
	WillAlwaysReturnStatus(httpStatus int)
}

type arcFixture struct {
	testing.TB
	transport         *httpmock.MockTransport
	knownTransactions map[string]*knownTransaction
}

func NewArcFixture(t testing.TB) ArcFixture {
	transport := httpmock.NewMockTransport()
	return NewArcFixtureWithTransport(t, transport)
}

func NewArcFixtureWithTransport(t testing.TB, transport *httpmock.MockTransport) ArcFixture {
	require.NotNil(t, transport, "http.RoundTripper must be provided")

	return &arcFixture{
		TB:                t,
		transport:         transport,
		knownTransactions: make(map[string]*knownTransaction),
	}
}

func (f *arcFixture) HttpClient() *resty.Client {
	client := resty.New()
	client.SetTransport(f.transport)
	return client
}

func (f *arcFixture) WillAlwaysReturnStatus(httpStatus int) {
	var details string
	switch httpStatus {
	case 400:
		details = "The request seems to be malformed and cannot be processed"
	case 401:
		details = "The request is not authorized"
	case 403:
		details = "The request is not authorized"
	case 404:
		details = "The requested resource could not be found"
	case 500:
		details = "The server encountered an internal error and was unable to complete your request"
	}

	f.transport.RegisterResponder("POST", "=~"+ArcURL+"/v1/tx.*", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewJsonResponse(httpStatus, map[string]any{
			"error":     details,
			"extraInfo": "",
			"instance":  nil,
			"status":    httpStatus,
			"title":     http.StatusText(httpStatus),
			"txid":      nil,
			"type":      "https://bitcoin-sv.github.io/arc/#/errors?id=_" + to.StringFromInteger(httpStatus),
		})
	})
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

		f.store(txHex)

		return f.knownTransactions[tx.TxID().String()].toResponseOrError()
	})

	f.transport.RegisterResponder("GET", "=~"+ArcURL+"/v1/tx/.*", func(req *http.Request) (*http.Response, error) {
		txid := req.URL.String()[len(ArcURL+"/v1/tx/"):]
		return f.knownTransactions[txid].toResponse()
	})
}

func (f *arcFixture) TxInfoJSON(id string) string {
	_, content := f.knownTransactions[id].toResponseContent()
	b, err := json.Marshal(content)
	require.NoError(f, err, "failed to marshal response content")
	return string(b)
}

func (f *arcFixture) store(txHex string) {
	beefBytes, err := hex.DecodeString(txHex)
	require.NoError(f, err, "failed to decode BEEF hex")

	beef, err := sdk.NewBeefFromBytes(beefBytes)
	require.NoError(f, err, "failed to create BEEF from bytes")

	transactions := seq2.FromMap(beef.Transactions)
	includedTransactions := seq2.MapTo(transactions, func(txID string, tx *sdk.BeefTx) *knownTransaction {
		merklePath := optional.OfPtr(tx.Transaction.MerklePath)

		return &knownTransaction{
			txid:        txID,
			status:      to.IfThen(merklePath.IsEmpty(), "SEEN_ON_NETWORK").ElseThen("MINED"),
			blockHeight: optional.Map(merklePath, func(it sdk.MerklePath) uint32 { return it.BlockHeight }).OrZeroValue(),
			blockHash: optional.Map(merklePath, func(it sdk.MerklePath) string {
				root, err := it.ComputeRootHex(&txID)
				require.NoError(f, err, "failed to compute root: wrong test setup")
				return root
			}).OrZeroValue(),
			merklePath: optional.Map(merklePath, func(it sdk.MerklePath) string { return it.Hex() }).OrZeroValue(),
		}
	})

	seq.ForEach(includedTransactions, func(it *knownTransaction) {
		f.knownTransactions[it.txid] = it
	})
}
