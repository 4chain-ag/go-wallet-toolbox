package to

import (
	"fmt"
	"math"
	"strconv"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/utils/types"
)

// ErrValueOutOfRange is returned when the value is out of range of the target type.
var ErrValueOutOfRange = fmt.Errorf("%w to convert", strconv.ErrRange)

const (
	maxIntForUnsigned   = uint64(math.MaxInt)
	maxInt8ForUnsigned  = uint64(math.MaxInt8)
	maxInt16ForUnsigned = uint64(math.MaxInt16)
	maxInt32ForUnsigned = uint64(math.MaxInt32)
	maxInt64ForUnsigned = uint64(math.MaxInt64)
)

// Int will convert any integer to int, with range checks
func Int[V types.SignedNumber](value V) (int, error) {
	valToCompare := int64(value)
	if valToCompare < math.MinInt || valToCompare > math.MaxInt {
		return 0, fmt.Errorf("%v %w to int", value, ErrValueOutOfRange)
	}
	return int(value), nil
}

// IntFromUnsigned will convert any unsigned integer to int, with range checks.
func IntFromUnsigned[V types.Unsigned](value V) (int, error) {
	if uint64(value) > maxIntForUnsigned {
		return 0, fmt.Errorf("%d %w to int", value, ErrValueOutOfRange)
	}
	return int(value), nil
}

// Int8 will convert any integer to int8, with range checks
func Int8[V types.SignedNumber](value V) (int8, error) {
	valToCompare := int64(value)
	if valToCompare < math.MinInt8 || valToCompare > math.MaxInt8 {
		return 0, fmt.Errorf("%v %w to int8", value, ErrValueOutOfRange)
	}
	return int8(value), nil
}

// Int8FromUnsigned will convert any unsigned integer to int8, with range checks.
func Int8FromUnsigned[V types.Unsigned](value V) (int8, error) {
	if uint64(value) > maxInt8ForUnsigned {
		return 0, fmt.Errorf("%d %w to int8", value, ErrValueOutOfRange)
	}
	return int8(value), nil
}

// Int16 will convert any integer to int16, with range checks
func Int16[V types.SignedNumber](value V) (int16, error) {
	valToCompare := int64(value)
	if valToCompare < math.MinInt16 || valToCompare > math.MaxInt16 {
		return 0, fmt.Errorf("%v %w to int16", value, ErrValueOutOfRange)
	}
	return int16(value), nil
}

// Int16FromUnsigned will convert any unsigned integer to int16, with range checks.
func Int16FromUnsigned[V types.Unsigned](value V) (int16, error) {
	if uint64(value) > maxInt16ForUnsigned {
		return 0, fmt.Errorf("%d %w to int16", value, ErrValueOutOfRange)
	}
	return int16(value), nil
}

// Int32 will convert any integer to int32, with range checks
func Int32[V types.SignedNumber](value V) (int32, error) {
	valToCompare := int64(value)
	if valToCompare < math.MinInt32 || valToCompare > math.MaxInt32 {
		return 0, fmt.Errorf("%v %w to int32", value, ErrValueOutOfRange)
	}
	return int32(value), nil
}

// Int32FromUnsigned will convert any unsigned integer to int32, with range checks.
func Int32FromUnsigned[V types.Unsigned](value V) (int32, error) {
	if uint64(value) > maxInt32ForUnsigned {
		return 0, fmt.Errorf("%d %w to int32", value, ErrValueOutOfRange)
	}
	return int32(value), nil
}

// Int64 will convert any integer to int64, with range checks
func Int64[V types.SignedNumber](value V) (int64, error) {
	return int64(value), nil
}

// Int64FromUnsigned will convert any unsigned integer to int64, with range checks.
func Int64FromUnsigned[V types.Unsigned](value V) (int64, error) {
	if uint64(value) > maxInt64ForUnsigned {
		return 0, fmt.Errorf("%d %w to int64", value, ErrValueOutOfRange)
	}
	return int64(value), nil
}

// UInt will convert any integer to uint, with range checks
func UInt[V types.Integer](value V) (uint, error) {
	if value < 0 || uint64(value) > math.MaxUint {
		return 0, fmt.Errorf("%d %w to uint", value, ErrValueOutOfRange)
	}
	return uint(value), nil
}

// UInt8 will convert any integer to uint8, with range checks
func UInt8[V types.Integer](value V) (uint8, error) {
	if value < 0 || uint64(value) > math.MaxUint8 {
		return 0, fmt.Errorf("%d %w to uint8", value, ErrValueOutOfRange)
	}
	return uint8(value), nil
}

// UInt16 will convert any integer to uint16, with range checks
func UInt16[V types.Integer](value V) (uint16, error) {
	if value < 0 || uint64(value) > math.MaxUint16 {
		return 0, fmt.Errorf("%d %w to uint16", value, ErrValueOutOfRange)
	}
	return uint16(value), nil
}

// UInt32 will convert any integer to uint32, with range checks
func UInt32[V types.Integer](value V) (uint32, error) {
	if value < 0 || uint64(value) > math.MaxUint32 {
		return 0, fmt.Errorf("%d %w to uint32", value, ErrValueOutOfRange)
	}
	return uint32(value), nil
}

// UInt64 will convert any integer to uint64, with range checks
func UInt64[V types.Integer](value V) (uint64, error) {
	if value < 0 {
		return 0, fmt.Errorf("%d %w to uint64", value, ErrValueOutOfRange)
	}
	return uint64(value), nil
}

// Float64FromInteger will convert any integer to float64
func Float64FromInteger[V types.Integer](value V) (float64, error) {
	return float64(value), nil
}
