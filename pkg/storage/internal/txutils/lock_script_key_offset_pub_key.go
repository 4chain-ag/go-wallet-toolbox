package txutils

import (
	"fmt"
	primitives "github.com/bsv-blockchain/go-sdk/primitives/ec"
	crypto "github.com/bsv-blockchain/go-sdk/primitives/hash"
	"github.com/bsv-blockchain/go-sdk/script"
	"github.com/bsv-blockchain/go-sdk/transaction/template/p2pkh"
)

// LockingScriptWithKeyOffset is a tool to generate a locking script with key offset.
type LockingScriptWithKeyOffset struct {
	offsetPrivGenerator func() (*primitives.PrivateKey, error)
}

// NewLockingScriptWithKeyOffset creates a new instance of LockingScriptWithKeyOffset.
func NewLockingScriptWithKeyOffset() *LockingScriptWithKeyOffset {
	return &LockingScriptWithKeyOffset{
		offsetPrivGenerator: randomPrivateKey,
	}
}

// Generate creates a locking script and randomizes a key offset (WIF formatted private key) from the given public key.
// NOTE: It is used to add Service Charge output to the transaction.
func (l *LockingScriptWithKeyOffset) Generate(pubKey string) (lockingScript string, keyOffset string, err error) {
	offsetPub, keyOffset, err := l.offsetPubKey(pubKey)

	address, err := script.NewAddressFromPublicKey(offsetPub, true)
	if err != nil {
		return "", "", fmt.Errorf("failed to create address from public key: %v", err)
	}

	lockingScriptObj, err := p2pkh.Lock(address)
	if err != nil {
		return "", "", fmt.Errorf("failed to create locking script: %v", err)
	}

	return lockingScriptObj.String(), keyOffset, nil

}

func (l *LockingScriptWithKeyOffset) offsetPubKey(pubKey string) (offsetPubKey *primitives.PublicKey, keyOffset string, err error) {
	pub, err := primitives.PublicKeyFromString(pubKey)
	if err != nil {
		return nil, "", fmt.Errorf("failed to parse public key: %v", err)
	}

	hashedSecret, keyOffset, err := l.keyOffsetToHashedSecret(pub)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get hashed secret: %v", err)
	}
	newPointX, newPointY := primitives.S256().ScalarBaseMult(hashedSecret)
	newPubKeyX, newPubKeyY := primitives.S256().Add(newPointX, newPointY, pub.X, pub.Y)
	offsetPubKey = &primitives.PublicKey{
		Curve: primitives.S256(),
		X:     newPubKeyX,
		Y:     newPubKeyY,
	}

	return offsetPubKey, keyOffset, nil
}

func (l *LockingScriptWithKeyOffset) keyOffsetToHashedSecret(pub *primitives.PublicKey) (hashedSecret []byte, keyOffset string, err error) {
	offset, err := l.offsetPrivGenerator()
	if err != nil {
		return nil, "", fmt.Errorf("failed to create new private key for keyOffset: %v", err)
	}

	sharedSecret, err := offset.DeriveSharedSecret(pub)
	if err != nil {
		return nil, "", fmt.Errorf("failed to derive shared secret: %v", err)
	}
	hashedSecret = crypto.Sha256(sharedSecret.ToDER())

	return hashedSecret, offset.Wif(), nil
}

func randomPrivateKey() (*primitives.PrivateKey, error) {
	privKey, err := primitives.NewPrivateKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %v", err)
	}
	return privKey, nil
}
