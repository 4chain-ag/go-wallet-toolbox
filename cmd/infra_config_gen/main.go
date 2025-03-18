package main

import (
	"github.com/4chain-ag/go-wallet-toolbox/pkg/infra"
)

func main() {
	filename := "infra-config.yaml"

	cfg := infra.Defaults()
	err := cfg.ToYAMLFile(filename)
	if err != nil {
		panic(err)
	}
}
