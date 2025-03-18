package defs

import "fmt"

type BSVNetwork string

const (
	NetworkMainnet BSVNetwork = "main"
	NetworkTestnet BSVNetwork = "test"
)

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
