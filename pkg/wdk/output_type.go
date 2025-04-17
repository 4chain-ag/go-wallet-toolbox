package wdk

// OutputType represents the type of output in a transaction.
type OutputType string

// Possible output types.
const (
	OutputTypeP2PKH  OutputType = "P2PKH"
	OutputTypeCustom OutputType = "custom"
)
