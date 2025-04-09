package defs

type Commission struct {
	Satoshis  uint64 `mapstructure:"satoshis"`
	PubKeyHex string `mapstructure:"pub_key_hex"`
}

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
