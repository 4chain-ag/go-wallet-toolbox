package defs

import "fmt"

// BSVNetwork represents the Bitcoin SV network type (mainnet or testnet)
type BSVNetwork string

// BSVNetwork constants for the different Bitcoin SV network types
const (
	NetworkMainnet BSVNetwork = "main"
	NetworkTestnet BSVNetwork = "test"
)

// ParseBSVNetworkStr will parse the given string and return the corresponding BSVNetwork type or an error
func ParseBSVNetworkStr(network string) (BSVNetwork, error) {
	switch BSVNetwork(network) {
	case NetworkTestnet:
		return NetworkTestnet, nil
	case NetworkMainnet:
		return NetworkMainnet, nil
	default:
		return "", fmt.Errorf("invalid network: %s", network)
	}
}
