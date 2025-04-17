package testabilities

import (
	"testing"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/satoshi"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/actions/funder/testabilities"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/testabilities/testusers"
)

type faucetFixture struct {
	t             testing.TB
	user          testusers.User
	funderFixture testabilities.FunderFixture
}

func (f *faucetFixture) TopUp(satoshis satoshi.Value) {
	f.t.Helper()

	f.funderFixture.UTXO().OwnedBy(f.user).WithSatoshis(satoshis.Int64()).P2PKH().Stored()
}
