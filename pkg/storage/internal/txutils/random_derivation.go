package txutils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func RandomDerivation(length uint64) (string, error) {
	if length == 0 {
		return "", fmt.Errorf("length cannot be zero")
	}

	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", fmt.Errorf("failed to read random bytes: %w", err)
	}

	return base64.StdEncoding.EncodeToString(randomBytes), nil
}
