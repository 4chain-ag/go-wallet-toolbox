package funder

import (
	"fmt"
	"math"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/defs"
	"github.com/go-softwarelab/common/pkg/to"
)

type feeCalc struct {
	bytes float64
	value uint64
}

func newFeeCalculator(model defs.FeeModel) *feeCalc {
	if model.Type != defs.SatPerKB {
		panic("unsupported fee model")
	}

	feeValue, err := to.UInt64(model.Value)
	if err != nil {
		panic("invalid fee model value: " + err.Error())
	}
	return &feeCalc{
		value: feeValue,
		bytes: 1000,
	}
}

func (f *feeCalc) Calculate(txSize uint64) (uint64, error) {
	size, err := to.Float64FromUnsigned(txSize)
	if err != nil {
		return 0, fmt.Errorf("invalid transaction size: %s", err.Error())
	}

	multiplier, err := to.UInt64(math.Ceil(size / f.bytes))
	if err != nil {
		return 0, fmt.Errorf("failed to calculate size / feeModel.bytes: %w", err)
	}

	fee, err := to.UInt64(multiplier * f.value)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate fee value: %w", err)
	}

	return fee, nil
}
