package testabilities

import (
	"testing"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/satoshi"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/actions/funder/testabilities"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/testabilities/testusers"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
)

type faucetFixture struct {
	t             testing.TB
	user          testusers.User
	funderFixture testabilities.FunderFixture
	basket        *wdk.TableOutputBasket
}

func (f *faucetFixture) TopUp(satoshis satoshi.Value) {
	f.t.Helper()

	f.funderFixture.UTXO().
		OwnedBy(f.user).
		InBasket(f.basket).
		WithSatoshis(satoshis.Int64()).
		P2PKH().
		Stored()
}
