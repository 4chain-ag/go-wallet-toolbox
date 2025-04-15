package txutils

import (
	"encoding/hex"
	"slices"

	crypto "github.com/bsv-blockchain/go-sdk/primitives/hash"
)

// TransactionIDFromRawTx will return a transactionID from the rawTx
func TransactionIDFromRawTx(rawTx []byte) string {
	hash := doubleSha256BE(rawTx)
	transactionID := hex.EncodeToString(hash)

	return transactionID
}

func doubleSha256BE(data []byte) []byte {
	doubleHash := crypto.Sha256d(data)
	slices.Reverse(doubleHash)

	return doubleHash
}
