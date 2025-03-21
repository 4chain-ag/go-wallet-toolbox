package wdk

import (
	"fmt"

	primitives "github.com/bsv-blockchain/go-sdk/primitives/ec"
)

// IdentityKey returns public key (DER HEX) from provided private key (HEX).
func IdentityKey(privKey string) (string, error) {
	rootKey, err := primitives.PrivateKeyFromHex(privKey)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %w", err)
	}

	return rootKey.PubKey().ToDERHex(), nil
}
