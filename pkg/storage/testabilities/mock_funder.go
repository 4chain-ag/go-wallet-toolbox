package testabilities

import (
	"context"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/actions"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
)

type MockFunder struct {
}

func (m *MockFunder) Fund(ctx context.Context, targetSat int64, currentTxSize uint64, basket *wdk.TableOutputBasket, userID int) (*actions.FundingResult, error) {
	return &actions.FundingResult{}, nil
}
