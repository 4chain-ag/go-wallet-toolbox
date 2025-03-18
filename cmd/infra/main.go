package main

import (
	"fmt"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/infra"
)

func main() {
	server, err := infra.NewServer(
		infra.WithConfigFile("infra-config.yaml"),
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Listening on :%d", server.Config.HTTPConfig.Port)

	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
