package testabilities

import (
	"context"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/actions"
)

type MockFunder struct {
}

func (m *MockFunder) Fund(ctx context.Context, targetSat int64, currentTxSize uint64, numberOfDesiredUTXOs int, minimumDesiredUTXOValue uint64, userID int) (*actions.FundingResult, error) {
	return &actions.FundingResult{}, nil
}
