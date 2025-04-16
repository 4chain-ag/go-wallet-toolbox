package satoshi

import (
	"fmt"
	"reflect"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk/primitives"
	"github.com/go-softwarelab/common/types"
)

type Value int64

func Add[A types.Integer, B types.Integer](a A, b B) (Value, error) {
	satsA, err := From(a)
	if err != nil {
		return 0, err
	}
	satsB, err := From(b)
	if err != nil {
		return 0, err
	}
	c := satsA + satsB
	if err = Validate(c); err != nil {
		return 0, err
	}
	return c, nil
}

func Subtract[A types.Integer, B types.Integer](a A, b B) (Value, error) {
	satsA, err := From(a)
	if err != nil {
		return 0, err
	}
	satsB, err := From(b)
	if err != nil {
		return 0, err
	}
	c := satsA - satsB
	if err = Validate(c); err != nil {
		return 0, err
	}
	return c, nil
}

func From[T types.Integer](value T) (Value, error) {
	if err := Validate(value); err != nil {
		return 0, err
	}
	return Value(value), nil
}

func Validate[T types.Integer](value T) error {
	var anyType any = value
	switch typed := anyType.(type) {
	case int:
		return validateInt(typed)
	case int64:
		return validateInt(typed)
	case uint:
		return validateUint(typed)
	case uint64:
		return validateUint(typed)
	case Value:
		return validateInt(typed)
	default:
		return validateGeneric(typed)
	}
}

func validateInt[T ~int | ~int64](value T) error {
	if value > primitives.MaxSatoshis {
		return fmt.Errorf("satoshi value exceeded max value")
	}
	if value < -primitives.MaxSatoshis {
		return fmt.Errorf("negative satoshi value exceeded max value")
	}
	return nil
}

func validateUint[T ~uint | ~uint64](value T) error {
	if value > primitives.MaxSatoshis {
		return fmt.Errorf("satoshi value exceeded max value")
	}
	return nil
}

func validateGeneric(value any) error {
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return validateInt(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return validateUint(v.Uint())
	default:
		return fmt.Errorf("unsupported type in validateGeneric")
	}
}
