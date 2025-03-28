package defs

import "fmt"

// FeeModelType represents different fee models which can be configured.
type FeeModelType string

// SatPerKB - Currently only "sat/kb" is supported.
const (
	SatPerKB FeeModelType = "sat/kb"
)

// ParseFeeModelType parses a string into a FeeModelType (case-insensitive).
func ParseFeeModelType(str string) (FeeModelType, error) {
	return parseEnumCaseInsensitive(str, SatPerKB)
}

// FeeModel represents a fee model with its value.
type FeeModel struct {
	Model FeeModelType `mapstructure:"model"`
	Value int64        `mapstructure:"value"`
}

// Validate double checks if under the Model is a valid enum, and checks if the value is valid.
func (f *FeeModel) Validate() error {
	var err error
	if f.Model, err = ParseFeeModelType(string(f.Model)); err != nil {
		return fmt.Errorf("invalid fee model: %s", f.Model)
	}
	if f.Model == SatPerKB && f.Value < 0 {
		return fmt.Errorf("invalid fee value: %d", f.Value)
	}
	return nil
}

// DefaultFeeModel returns minimal fee model.
func DefaultFeeModel() FeeModel {
	return FeeModel{
		Model: SatPerKB,
		Value: 1,
	}
}
