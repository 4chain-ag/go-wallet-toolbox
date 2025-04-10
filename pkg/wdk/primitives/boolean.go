package primitives

// BooleanDefaultFalse is a boolean with a default value of false
type BooleanDefaultFalse bool

// Value returns the boolean value with a default when nil
func (b *BooleanDefaultFalse) Value() bool {
	if b == nil {
		return false
	}
	return bool(*b)
}

// BooleanDefaultTrue is a boolean with a default value of true
type BooleanDefaultTrue bool

// Value returns the boolean value with a default when nil
func (b *BooleanDefaultTrue) Value() bool {
	if b == nil {
		return true
	}
	return bool(*b)
}
