package validate

import (
	"fmt"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"
)

func ValidCreateActionArgs(args *wdk.ValidCreateActionArgs) error {
	deducedIsSendWith := len(args.Options.SendWith) > 0
	if args.IsSendWith != deducedIsSendWith {
		return fmt.Errorf("inconsistent IsSendWith with Options.SendWith")
	}

	deducedIsRemixChange := !args.IsSendWith && len(args.Inputs) == 0 && len(args.Outputs) == 0
	if args.IsRemixChange != deducedIsRemixChange {
		return fmt.Errorf("inconsistent IsRemixChange with IsSendWith and Inputs and Outputs")
	}

	deducedIsNewTx := args.IsRemixChange || len(args.Inputs) > 0 || len(args.Outputs) > 0
	if args.IsNewTx != deducedIsNewTx {
		return fmt.Errorf("inconsistent IsNewTx with IsRemixChange and Inputs and Outputs")
	}

	if !args.IsNewTx {
		return fmt.Errorf("create action is meant to create a new transaction")
	}

	deducedIsSignAction := args.IsNewTx && !args.Options.SignAndProcess && containsNilUnlockingScript(args.Inputs)
	if args.IsSignAction != deducedIsSignAction {
		return fmt.Errorf("inconsistent IsSignAction with IsNewTx and Options.SignAndProcess and Inputs.UnlockingScript")
	}

	deducedIsDelayed := bool(args.Options.AcceptDelayedBroadcast)
	if args.IsDelayed != deducedIsDelayed {
		return fmt.Errorf("inconsistent IsDelayed with Options.AcceptDelayedBroadcast")
	}

	deducedIsNoSend := bool(args.Options.NoSend)
	if args.IsNoSend != deducedIsNoSend {
		return fmt.Errorf("inconsistent IsNoSend with Options.NoSend")
	}

	if err := args.Description.Validate(); err != nil {
		return fmt.Errorf("the description parameter must be %w", err)
	}

	for i, label := range args.Labels {
		if err := label.Validate(); err != nil {
			return fmt.Errorf("label as %d must be %w", i, err)
		}
	}

	for i, input := range args.Inputs {
		if err := validateCreateActionInput(&input); err != nil {
			return fmt.Errorf("invalid input as %d: %w", i, err)
		}
	}

	for i, output := range args.Outputs {
		if err := validateCreateActionOutput(&output); err != nil {
			return fmt.Errorf("invalid output as %d: %w", i, err)
		}
	}

	return nil
}

func validateCreateActionInput(input *wdk.ValidCreateActionInput) error {
	if input.UnlockingScript == nil && input.UnlockingScriptLength == nil {
		return fmt.Errorf("at least one of unlockingScript, unlockingScriptLength must be set")
	}

	if input.UnlockingScript != nil {
		if err := input.UnlockingScript.Validate(); err != nil {
			return fmt.Errorf("unlockingScript must be %w", err)
		}

		if input.UnlockingScriptLength != nil && uint(len(*input.UnlockingScript)) != uint(*input.UnlockingScriptLength) {
			return fmt.Errorf("unlockingScriptLength must match provided unlockingScript length")
		}
	}

	if err := input.InputDescription.Validate(); err != nil {
		return fmt.Errorf("inputDescription must be %w", err)
	}

	return nil
}

func validateCreateActionOutput(output *wdk.ValidCreateActionOutput) error {
	if err := output.LockingScript.Validate(); err != nil {
		return fmt.Errorf("lockingScript must be %w", err)
	}

	if err := output.Satoshis.Validate(); err != nil {
		return fmt.Errorf("satoshis must be %w", err)
	}

	if err := output.OutputDescription.Validate(); err != nil {
		return fmt.Errorf("outputDescription must be %w", err)
	}

	if output.Basket != nil {
		if err := output.Basket.Validate(); err != nil {
			return fmt.Errorf("basket must be %w", err)
		}
	}

	for _, tag := range output.Tags {
		if err := tag.Validate(); err != nil {
			return fmt.Errorf("tag must be %w", err)
		}
	}

	return nil
}

func containsNilUnlockingScript(inputs []wdk.ValidCreateActionInput) bool {
	for _, input := range inputs {
		if input.UnlockingScript == nil {
			return true
		}
	}
	return false
}
