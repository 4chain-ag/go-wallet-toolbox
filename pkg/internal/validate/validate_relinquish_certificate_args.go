package validate

import (
	"fmt"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
)

func ValidateRelinquishCertificateArgs(args *wdk.RelinquishCertificateArgs) error {
	err := args.Type.Validate()
	if err != nil {
		return fmt.Errorf("invalid type argument: %w", err)
	}

	err = args.SerialNumber.Validate()
	if err != nil {
		return fmt.Errorf("invalid serialNumber argument: %w", err)
	}

	err = args.Certifier.Validate()
	if err != nil {
		return fmt.Errorf("invalid certifier argument: %w", err)
	}
	return nil
}
