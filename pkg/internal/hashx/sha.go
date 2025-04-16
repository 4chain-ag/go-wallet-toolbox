package hashx

import (
	"slices"

	crypto "github.com/bsv-blockchain/go-sdk/primitives/hash"
)

func DoubleSha256LE(data []byte) []byte {
	return crypto.Sha256d(data)
}

func DoubleSha256BE(data []byte) []byte {
	doubleHash := DoubleSha256LE(data)
	slices.Reverse(doubleHash)

	return doubleHash
}
