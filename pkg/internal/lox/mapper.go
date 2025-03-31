package lox

// Iteratee represents a function that processes an item of type T at a specific index,
// returning a result of type R. It is useful for operations requiring both the item and its index.
type Iteratee[T, R any] func(item T, index int) R

// NoIndexIteratee represents a function that processes an item of type T
// without requiring its index, returning a result of type R.
// It is useful for operations where the index is irrelevant.
type NoIndexIteratee[T, R any] func(item T) R

// MappingFn allows to pass a mapper function that doesn't require index
// to "github.com/samber/lo" Map() function
func MappingFn[T, R any](mapper NoIndexIteratee[T, R]) Iteratee[T, R] {
	return func(item T, _ int) R {
		return mapper(item)
	}
}
