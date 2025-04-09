package main

import (
	"fmt"
	"net/http"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/services"
)

func main() {

	service := services.New(&http.Client{}, "test")

	r, _ := service.BsvExchangeRate()
	fmt.Println()
	fmt.Println(r)
}
