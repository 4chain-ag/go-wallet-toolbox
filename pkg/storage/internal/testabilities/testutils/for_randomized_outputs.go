package testutils

import (
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk/primitives"
	"github.com/go-softwarelab/common/pkg/seq"
	"github.com/stretchr/testify/require"
	"slices"
	"testing"
)

func FindOutput[T any](
	t *testing.T,
	outputs []T,
	finder func(p T) bool,
) (*T, uint32) {
	t.Helper()
	index := slices.IndexFunc(outputs, finder)
	require.GreaterOrEqual(t, index, 0)

	return &outputs[index], uint32(index)
}

func CountOutputsWithCondition[T any](
	t *testing.T,
	outputs []T,
	finder func(p T) bool,
) int {
	t.Helper()

	return seq.Count(seq.Filter(seq.FromSlice(outputs), finder))
}

func SumOutputsWithCondition[T any](
	t *testing.T,
	outputs []T,
	getter func(p T) primitives.SatoshiValue,
	finder func(p T) bool,
) primitives.SatoshiValue {
	t.Helper()

	sum := primitives.SatoshiValue(0)
	for _, output := range outputs {
		if finder(output) {
			sum += getter(output)
		}
	}
	return sum
}

func ForEveryOutput[T any](
	t *testing.T,
	outputs []T,
	finder func(p T) bool,
	validator func(p T),
) {
	t.Helper()

	for _, output := range outputs {
		if finder(output) {
			validator(output)
		}
	}
}
