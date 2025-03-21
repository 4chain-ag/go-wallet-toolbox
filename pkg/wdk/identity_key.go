package wdk

import (
	"fmt"
	primitives "github.com/bsv-blockchain/go-sdk/primitives/ec"
)

func IdentityKey(privKey string) (string, error) {
	rootKey, err := primitives.PrivateKeyFromHex(privKey)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %w", err)
	}

	return rootKey.PubKey().ToDERHex(), nil
}
