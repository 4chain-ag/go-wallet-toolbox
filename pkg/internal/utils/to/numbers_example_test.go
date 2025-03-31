package to_test

import (
	"errors"
	"fmt"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/utils/to"
)

func ExampleInt() {
	// Converting within range
	val, err := to.Int(int16(42))
	fmt.Printf("%T(%d), Error: %v\n", val, val, err)

	// Output:
	// int(42), Error: <nil>
}

func ExampleInt8() {
	// Converting within range
	val, err := to.Int8(42)
	fmt.Printf("%T(%d), Error: %v\n", val, val, err)

	// Converting out of range
	valOOR, errOOR := to.Int8(1000)
	fmt.Printf("%T(%d), Error: %v\n", valOOR, valOOR, errors.Is(errOOR, to.ErrValueOutOfRange))

	// Output:
	// int8(42), Error: <nil>
	// int8(0), Error: true
}

func ExampleUInt32() {
	// Valid conversion
	val, err := to.UInt32(42)
	fmt.Printf("%T(%d), Error: %v\n", val, val, err)

	// Negative number
	valNeg, errNeg := to.UInt32(-5)
	fmt.Printf("%T(%d), Error: %v\n", valNeg, valNeg, errors.Is(errNeg, to.ErrValueOutOfRange))

	// Output:
	// uint32(42), Error: <nil>
	// uint32(0), Error: true
}
