package defs

// Commission represents the commission configuration for a storage provider.
// If satoshis is greater than 0, it means that the commission is enabled.
type Commission struct {
	Satoshis  uint64 `mapstructure:"satoshis"`
	PubKeyHex string `mapstructure:"pub_key_hex"`
}

// Enabled checks if the commission is enabled.
func (c *Commission) Enabled() bool {
	return c.Satoshis > 0
}

// Validate double checks if under the Type is a valid enum, and checks if the value is valid.
func (c *Commission) Validate() error {
	if !c.Enabled() {
		return nil
	}

	// TODO: Validate PubKeyHex

	return nil
}

// DefaultCommission returns a default commission configuration - disabled.
func DefaultCommission() Commission {
	return Commission{
		Satoshis:  0,
		PubKeyHex: "",
	}
}
