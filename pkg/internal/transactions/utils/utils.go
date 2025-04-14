package utils

import (
	"encoding/hex"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/hashx"
)

func TransactionIDFromRawTx(rawTx []byte) string {
	hash := hashx.DoubleSha256BE(rawTx)
	transactionID := hex.EncodeToString(hash)

	return transactionID
}
