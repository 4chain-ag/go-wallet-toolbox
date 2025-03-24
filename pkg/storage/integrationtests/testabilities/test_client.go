package testabilities

import "github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"

type TestClient struct {
	MakeAvailable func() (*wdk.TableSettings, error)
}
