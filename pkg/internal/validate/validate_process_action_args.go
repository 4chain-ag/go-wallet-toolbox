package validate

import (
	"fmt"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
)

func ProcessActionArgs(args *wdk.ProcessActionArgs) error {
	if args.TxID != nil {
		if err := args.TxID.Validate(); err != nil {
			return fmt.Errorf("invalid txID argument: %w", err)
		}
	}

	return nil
}
