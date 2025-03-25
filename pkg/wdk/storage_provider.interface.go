package wdk

// StorageProvider is an interface that defines the methods that a storage provider should implement
type StorageProvider interface {
	Migrate(storageName, storageIdentityKey string) (string, error)
	MakeAvailable() (*TableSettings, error)
}
