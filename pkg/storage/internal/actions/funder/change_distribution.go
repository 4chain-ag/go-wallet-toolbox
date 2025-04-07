package funder

import (
	"crypto/rand"
	"fmt"
	"iter"
	"math/big"

	"github.com/go-softwarelab/common/pkg/seq"
	"github.com/go-softwarelab/common/pkg/to"
)

type Randomizer func(max uint64) uint64

type ChangeDistribution struct {
	initialValue uint64
	randomizer   Randomizer
}

func NewChangeDistribution(initialValue uint64, randomizer Randomizer) *ChangeDistribution {
	return &ChangeDistribution{
		initialValue: initialValue,
		randomizer:   randomizer,
	}
}

func (d *ChangeDistribution) Distribute(count uint64, amount uint64) iter.Seq[uint64] {
	if count == 0 || amount == 0 {
		return seq.Of[uint64]()
	}
	if count == 1 {
		return seq.Of(amount)
	}

	countSignedInt := int(count) //nolint:gosec // count is always > 1 at this point

	saturationThreshold := count * d.initialValue
	if saturationThreshold < amount {
		base := amount / count
		reminder := amount % count
		countMinusOne := countSignedInt - 1

		distribution := seq.Concat(
			seq.Of[uint64](base+reminder),
			seq.Repeat(base, countMinusOne),
		)

		noise := d.randomNoise(count, distribution)

		var i uint64
		var v uint64
		distribution = seq.Map(distribution, func(current uint64) uint64 {
			v = current - noise[i] + noise[count-i-1]
			i++
			return v
		})

		return distribution
	}

	if saturationThreshold == amount {
		return seq.Repeat(d.initialValue, countSignedInt)
	}

	// not saturated - at least one output is less than initialValue:
	for i := uint64(1); i < count; i++ {
		saturatedOutputs := count - i
		b := saturatedOutputs * d.initialValue
		if amount > b {
			return seq.Concat(
				seq.Of[uint64](amount-b),
				seq.Repeat(d.initialValue, int(saturatedOutputs)), //nolint:gosec // saturatedOutputs is always > 0 at this point
			)
		}
	}

	return seq.Of(amount)
}

func (d *ChangeDistribution) randomNoise(count uint64, distribution iter.Seq[uint64]) []uint64 {
	noise := make([]uint64, 0, count)
	for current := range distribution {
		randomRange := current - d.initialValue
		var randomized uint64
		if randomRange != 0 {
			randomized = d.randomizer(randomRange)
		}
		noise = append(noise, randomized)
	}
	return noise
}

func Rand(max uint64) uint64 {
	maxI64, err := to.Int64FromUnsigned(max)
	if err != nil {
		panic(fmt.Errorf("rand: cannot convert max value to signed int: %w", err))
	}
	nBig, err := rand.Int(rand.Reader, big.NewInt(maxI64))
	if err != nil {
		panic(fmt.Errorf("failed to generate random number: %w", err))
	}
	return nBig.Uint64()
}
