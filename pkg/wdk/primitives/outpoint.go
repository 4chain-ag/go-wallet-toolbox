package primitives

import (
	"fmt"
	"strings"

	"github.com/go-softwarelab/common/pkg/to"
)

// OutpointString represents a transaction ID and output index pair.
// The TXID is given as a hex string followed by a period "." and then the output index is given as a decimal integer.
type OutpointString string

// Validate checks if the string is proper outpoint string and contains outpoint index after "."
func (s OutpointString) Validate() error {
	split := strings.Split(string(s), ".")

	if len(split) != 2 {
		return fmt.Errorf("txid as hexstring and numeric output index joined with '.'")
	}

	// check if after decimal point there is an outpoint index
	_, err := to.UInt64FromString(split[1])
	if err != nil {
		return fmt.Errorf("txid as hexstring and numeric output index joined with '.'")
	}

	return nil
}
