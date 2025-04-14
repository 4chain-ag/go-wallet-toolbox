package txutils

import (
	"github.com/go-softwarelab/common/pkg/seq"
	"github.com/go-softwarelab/common/types"
	"iter"
)

func NewRandomSequence[T any, E types.Integer](length E, randomizer func() (T, error)) (iter.Seq[T], error) {
	list := make([]T, length)
	for i := E(0); i < length; i++ {
		item, err := randomizer()
		if err != nil {
			return nil, err
		}
		list[i] = item
	}

	return seq.FromSlice(list), nil
}
