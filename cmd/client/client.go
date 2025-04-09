package main

import (
	"fmt"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/services"
	"net/http"
)

func main() {

	service := services.New(&http.Client{}, "test")

	r, _ := service.BsvExchangeRate()
	fmt.Println()
	fmt.Println(r)
}
