package wdk

// WalletStorageWriter is an interface for writing to the wallet storage
type WalletStorageWriter interface {
	Migrate(storageName, storageIdentityKey string) (string, error)
	MakeAvailable() (*TableSettings, error)
	FindOrInsertUser(identityKey string) (*TableUser, error)
	CreateAction(auth AuthID, args ValidCreateActionArgs)
}
