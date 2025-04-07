package funder

import (
	"github.com/go-softwarelab/common/pkg/seq"
	"github.com/stretchr/testify/require"
	"testing"
)

func mockZeroRandomizer(_ uint64) uint64 {
	return 0
}

func TestChangeDistribution(t *testing.T) {
	tests := map[string]struct {
		initialValue uint64
		randomizer   func(uint64) uint64
		count        uint64
		amount       uint64
		expected     []uint64
	}{
		"not saturated: reminder + (count-1) * initialValue": {
			initialValue: 1000,
			randomizer:   mockZeroRandomizer,
			count:        6,
			amount:       5500,
			expected:     []uint64{500, 1000, 1000, 1000, 1000, 1000},
		},
		"not saturated: initialValue/4 + (count-1) * initialValue": {
			initialValue: 1000,
			randomizer:   mockZeroRandomizer,
			count:        6,
			amount:       5250,
			expected:     []uint64{250, 1000, 1000, 1000, 1000, 1000},
		},
		"not saturated: 1 + (count-1) * initialValue": {
			// NOTE: This case should not happen
			// Change output should not be below "fistValue" argument (which is by default = initialValue / 4)
			// In real life, provided number of outputs (count) should be decreased
			initialValue: 1000,
			randomizer:   mockZeroRandomizer,
			count:        6,
			amount:       5001,
			expected:     []uint64{1, 1000, 1000, 1000, 1000, 1000},
		},
		"reduced count: (count-1) * initialValue": {
			// NOTE: This case should not happen
			// Algorithm could not find a solution for target count, and automatically reduced the number of outputs
			// In real life, provided number of outputs (count) should be decreased
			initialValue: 1000,
			randomizer:   mockZeroRandomizer,
			count:        6,
			amount:       5000,
			expected:     []uint64{1000, 1000, 1000, 1000, 1000},
		},
		"not saturated, reduced count: (count-1) * initialValue": {
			// NOTE: This case should not happen
			// Algorithm could not find a solution for target count, and automatically reduced the number of outputs
			// In real life, provided number of outputs (count) should be decreased
			initialValue: 1000,
			randomizer:   mockZeroRandomizer,
			count:        6,
			amount:       4999,
			expected:     []uint64{999, 1000, 1000, 1000, 1000},
		},
		"equally saturated: (count) * initialValue": {
			// NOTE: This case should not happen
			// Algorithm could not find a solution for target count, and automatically reduced the number of outputs
			// In real life, provided number of outputs (count) should be decreased
			initialValue: 1000,
			randomizer:   mockZeroRandomizer,
			count:        6,
			amount:       6000,
			expected:     []uint64{1000, 1000, 1000, 1000, 1000, 1000},
		},
		"saturated: equal distribution": {
			// NOTE: This case should not happen
			// Algorithm could not find a solution for target count, and automatically reduced the number of outputs
			// In real life, provided number of outputs (count) should be decreased
			initialValue: 1000,
			randomizer:   mockZeroRandomizer,
			count:        6,
			amount:       7200,
			expected:     []uint64{1200, 1200, 1200, 1200, 1200, 1200},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// given:
			dist := NewChangeDistribution(test.initialValue, test.randomizer)

			// when:
			values := dist.Distribute(test.count, test.amount)

			// then:
			require.EqualValues(t, test.expected, seq.Collect(values))
		})
	}
}

func TestChangeDistribution2(t *testing.T) {
	// given:
	initialValue := uint64(1000)
	randomizer := mockZeroRandomizer

	// and:
	count := uint64(6)
	amount := uint64(5500)

	// and:
	dist := NewChangeDistribution(initialValue, randomizer)

	// when:
	values := dist.Distribute(count, amount)

	// then:
	require.EqualValues(t, []uint64{500, 1000, 1000, 1000, 1000, 1000}, seq.Collect(values))
}
