package funder

import (
	"crypto/rand"
	"fmt"
	"github.com/go-softwarelab/common/pkg/seq"
	"iter"
	"math/big"
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
	if count == 0 {
		return seq.Of[uint64]()
	}
	if count == 1 {
		return seq.Of(amount)
	}

	a := count * d.initialValue
	if a < amount {
		base := amount / count
		reminder := amount % count

		distribution := seq.Concat(
			seq.Of[uint64](base+reminder),
			seq.Repeat(base, int(count-1)),
		)

		noise := seq.ToSlice(
			seq.Map(distribution, func(current uint64) uint64 {
				randomRange := current - d.initialValue
				if randomRange == 0 {
					return 0
				}
				return d.randomizer(randomRange)
			}),
			make([]uint64, 0, count),
		)

		var i uint64
		var v uint64
		distribution = seq.Map(distribution, func(current uint64) uint64 {
			v = current - noise[i] + noise[count-i-1]
			i++
			return v
		})

		return distribution
	}

	if a == amount {
		return seq.Repeat(d.initialValue, int(count))
	}

	// a > amount
	for i := uint64(1); i < count; i++ {
		j := count - i
		b := j * d.initialValue
		if amount > b {
			return seq.Concat(
				seq.Of[uint64](amount-b),
				seq.Repeat(d.initialValue, int(j)),
			)
		}
	}

	return seq.Of(amount)
}

func random(max uint64) uint64 {
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		panic(fmt.Errorf("failed to generate random number: %w", err))
	}
	return nBig.Uint64()
}
