package testabilities

import "github.com/4chain-ag/go-wallet-toolbox/pkg/wdk"

type TestClient struct {
	Migrate          func(storageName, storageIdentityKey string) (string, error)
	MakeAvailable    func() (*wdk.TableSettings, error)
	FindOrInsertUser func(identityKey string) (*wdk.FindOrInsertUserResponse, error)
}
