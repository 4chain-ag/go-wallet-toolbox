package wdk

//go:generate go run -tags gen ../../tools/client-gen/main.go -out client_gen.go

// WalletStorageWriter is an interface for writing to the wallet storage
type WalletStorageWriter interface {
	Migrate(storageName, storageIdentityKey string) (string, error)
	MakeAvailable() (*TableSettings, error)
	FindOrInsertUser(identityKey string) (*TableUser, error)
	CreateAction(auth AuthID, args ValidCreateActionArgs) (*StorageCreateActionResult, error)
}
