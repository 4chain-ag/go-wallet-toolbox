package utils

import (
	"encoding/hex"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/hashx"
)

// TransactionIDFromRawTx will return a transactionID from the rawTx
func TransactionIDFromRawTx(rawTx []byte) string {
	hash := hashx.DoubleSha256BE(rawTx)
	transactionID := hex.EncodeToString(hash)

	return transactionID
}
