package wdk

//go:generate go tool mockgen -destination=../internal/mocks/mock_wallet_storage_writer.go -package=mocks github.com/4chain-ag/go-wallet-toolbox/pkg/wdk WalletStorageWriter

// WalletStorageWriter is an interface for writing to the wallet storage
type WalletStorageWriter interface {
	Migrate(storageName, storageIdentityKey string) (string, error)
	MakeAvailable() (*TableSettings, error)
	FindOrInsertUser(identityKey string) (*FindOrInsertUserResponse, error)
	CreateAction(auth AuthID, args ValidCreateActionArgs) (*StorageCreateActionResult, error)
}
