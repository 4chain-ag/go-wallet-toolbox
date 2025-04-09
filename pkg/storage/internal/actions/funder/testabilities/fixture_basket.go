package testabilities

import (
	"testing"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/testabilities/testusers"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
)

const (
	desiredUTXONumberToPreferSingleChange = 1
	testDesiredUTXOValue                  = 1000
)

type BasketFixture interface {
	ThatPrefersSingleChange() *wdk.TableOutputBasket
	WithNumberOfDesiredUTXOs(i int) *wdk.TableOutputBasket
}

type basketFixture struct {
	testing.TB
	user testusers.User
}

func newBasketFixture(t testing.TB, user testusers.User) *basketFixture {
	return &basketFixture{
		TB:   t,
		user: user,
	}
}

func (f *basketFixture) ThatPrefersSingleChange() *wdk.TableOutputBasket {
	return f.WithNumberOfDesiredUTXOs(desiredUTXONumberToPreferSingleChange)
}

func (f *basketFixture) WithNumberOfDesiredUTXOs(number int) *wdk.TableOutputBasket {
	return &wdk.TableOutputBasket{
		BasketID: 1,
		UserID:   f.user.ID,
		BasketConfiguration: wdk.BasketConfiguration{
			Name:                    "default",
			NumberOfDesiredUTXOs:    number,
			MinimumDesiredUTXOValue: testDesiredUTXOValue,
		},
		CreatedAt: exampleDate,
		UpdatedAt: exampleDate,
		IsDeleted: false,
	}
}
