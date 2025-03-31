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

func ExampleIntFromUnsigned() {
	// Converting within range
	val, err := to.IntFromUnsigned(uint64(42))
	fmt.Printf("%T(%d), Error: %v\n", val, val, err)

	// Output:
	// int(42), Error: <nil>
}

func ExampleInt16() {
	// Converting within range
	val, err := to.Int16(42)
	fmt.Printf("%T(%d), Error: %v\n", val, val, err)

	// Converting out of range
	valOOR, errOOR := to.Int16(40000)
	fmt.Printf("%T(%d), Error: %v\n", valOOR, valOOR, errors.Is(errOOR, to.ErrValueOutOfRange))

	// Output:
	// int16(42), Error: <nil>
	// int16(0), Error: true
}

func ExampleInt32() {
	// Converting within range
	val, err := to.Int32(42)
	fmt.Printf("%T(%d), Error: %v\n", val, val, err)

	// Output:
	// int32(42), Error: <nil>
}

func ExampleInt64() {
	// Converting within range
	val, err := to.Int64(42)
	fmt.Printf("%T(%d), Error: %v\n", val, val, err)

	// Output:
	// int64(42), Error: <nil>
}

func ExampleUInt() {
	// Valid conversion
	val, err := to.UInt(42)
	fmt.Printf("%T(%d), Error: %v\n", val, val, err)

	// Negative number
	valNeg, errNeg := to.UInt(-5)
	fmt.Printf("%T(%d), Error: %v\n", valNeg, valNeg, errors.Is(errNeg, to.ErrValueOutOfRange))

	// Output:
	// uint(42), Error: <nil>
	// uint(0), Error: true
}

func ExampleUInt8() {
	// Valid conversion
	val, err := to.UInt8(42)
	fmt.Printf("%T(%d), Error: %v\n", val, val, err)

	// Out of range
	valOOR, errOOR := to.UInt8(300)
	fmt.Printf("%T(%d), Error: %v\n", valOOR, valOOR, errors.Is(errOOR, to.ErrValueOutOfRange))

	// Output:
	// uint8(42), Error: <nil>
	// uint8(0), Error: true
}

func ExampleUInt16() {
	// Valid conversion
	val, err := to.UInt16(42)
	fmt.Printf("%T(%d), Error: %v\n", val, val, err)

	// Output:
	// uint16(42), Error: <nil>
}

func ExampleUInt64() {
	// Valid conversion
	val, err := to.UInt64(42)
	fmt.Printf("%T(%d), Error: %v\n", val, val, err)

	// Negative number
	valNeg, errNeg := to.UInt64(-5)
	fmt.Printf("%T(%d), Error: %v\n", valNeg, valNeg, errors.Is(errNeg, to.ErrValueOutOfRange))

	// Output:
	// uint64(42), Error: <nil>
	// uint64(0), Error: true
}

func ExampleFloat64FromInteger() {
	// Converting to float64
	val, err := to.Float64FromInteger(42)
	fmt.Printf("%T(%.1f), Error: %v\n", val, val, err)

	// Output:
	// float64(42.0), Error: <nil>
}
