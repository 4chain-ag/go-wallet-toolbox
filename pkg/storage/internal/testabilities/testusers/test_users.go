package testusers

import (
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
	primitives "github.com/bsv-blockchain/go-sdk/primitives/ec"
	"github.com/stretchr/testify/require"
	"testing"
)

// NOTE: Testabilities can modify user IDs, to match ID with database

type User struct {
	// Name of the user just for in tests logging purpose
	Name    string
	ID      int
	PrivKey string
}

var Alice = User{
	Name:    "Alice",
	ID:      1,
	PrivKey: "143ab18a84d3b25e1a13cefa90038411e5d2014590a2a4a57263d1593c8dee1c",
}

var Bob = User{
	Name:    "Bob",
	ID:      2,
	PrivKey: "0881208859876fc227d71bfb8b91814462c5164b6fee27e614798f6e85d2547d",
}

func (u User) AuthID() wdk.AuthID {
	return wdk.AuthID{
		UserID: &u.ID,
	}
}

func (u User) PubKey(t *testing.T) string {
	t.Helper()

	priv, err := primitives.PrivateKeyFromHex(u.PrivKey)
	require.NoError(t, err)

	return priv.PubKey().ToDERHex()
}

func All() []*User {
	return []*User{&Alice, &Bob}
}
