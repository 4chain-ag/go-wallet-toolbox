package testabilities

import (
	"fmt"
	"net/http"

	"github.com/jarcoal/httpmock"
)

type knownTransaction struct {
	txid        string
	status      string
	blockHeight uint32
	blockHash   string
	merklePath  string
}

func (t *knownTransaction) toResponse() (*http.Response, error) {
	return httpmock.NewJsonResponse(t.toResponseContent())
}

func (t *knownTransaction) toResponseOrError() (*http.Response, error) {
	if t == nil {
		return nil, fmt.Errorf("unexpectedly cannot find transaction in known transactions")
	}
	return t.toResponse()
}

func (t *knownTransaction) toResponseContent() (int, map[string]any) {
	if t == nil {
		return http.StatusNotFound, map[string]any{
			"detail":    "The requested resource could not be found",
			"extraInfo": "transaction not found",
			"instance":  nil,
			"status":    http.StatusNotFound,
			"title":     "Not found",
			"txid":      nil,
			"type":      "https://bitcoin-sv.github.io/arc/#/errors?id=_404",
		}
	}

	return http.StatusOK, map[string]any{
		"blockHash":    t.blockHash,
		"blockHeight":  t.blockHeight,
		"competingTxs": nil,
		"extraInfo":    "",
		"merklePath":   t.merklePath,
		"timestamp":    timestamp,
		"txStatus":     t.status,
		"txid":         t.txid,
	}
}
