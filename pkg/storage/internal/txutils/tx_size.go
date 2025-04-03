package txutils

import (
	"fmt"
	"iter"
)

// VarUintSize returns the byte size required to encode number as Bitcoin VarUint
func VarUintSize(val uint64) uint64 {
	switch {
	case val <= 0xfc:
		return 1
	case val <= 0xffff:
		return 3
	case val <= 0xffffffff:
		return 5
	default:
		return 9
	}
}

// TransactionInputSize calculates the size in bytes of a transaction input
// with the given script size
func TransactionInputSize(scriptSize uint64) uint64 {
	return 32 + // txid
		4 + // vout
		VarUintSize(scriptSize) + // script size in bytes
		scriptSize + // script
		4 // sequence number
}

// TransactionOutputSize calculates the serialized byte length of a transaction output
// with the given script size in bytes
func TransactionOutputSize(scriptSize uint64) uint64 {
	return VarUintSize(scriptSize) + // output script length
		scriptSize + // output script
		8
}

// TransactionSize calculates the total size of a transaction in bytes
// inputs is a slice of script sizes for each input
// outputs is a slice of script sizes for each output
func TransactionSize(inputSizes iter.Seq2[uint64, error], outputSizes iter.Seq2[uint64, error]) (uint64, error) {
	var inputsCount uint64
	var inputsSize uint64
	for scriptSize, err := range inputSizes {
		if err != nil {
			return 0, fmt.Errorf("failed to calculate unlocking script size: %w", err)
		}
		inputsCount++
		inputsSize += TransactionInputSize(scriptSize)
	}

	var outputsCount uint64
	var outputsSize uint64
	for scriptSize, err := range outputSizes {
		if err != nil {
			return 0, fmt.Errorf("failed to calculate locking script size: %w", err)
		}
		outputsCount++
		outputsSize += TransactionOutputSize(scriptSize)
	}

	return 4 + // Version
			VarUintSize(inputsCount) + // Number of inputs
			inputsSize + // All inputs accumulated size
			VarUintSize(outputsCount) + // Number of outputs
			outputsSize + // All outputs accumulated size
			4, // lock time
		nil
}
