package wdk

//go:generate go run -tags gen ../../tools/client-gen/main.go -out client_gen.go
//go:generate go tool mockgen -destination=../internal/mocks/mock_wallet_storage_writer.go -package=mocks github.com/4chain-ag/go-wallet-toolbox/pkg/wdk WalletStorageWriter

// WalletStorageWriter is an interface for writing to the wallet storage
type WalletStorageWriter interface {
	Migrate(storageName, storageIdentityKey string) (string, error)
	MakeAvailable() (*TableSettings, error)
	FindOrInsertUser(identityKey string) (*FindOrInsertUserResponse, error)
	CreateAction(auth AuthID, args ValidCreateActionArgs) (*StorageCreateActionResult, error)

	InsertCertificateAuth(auth AuthID, certificate *TableCertificateX) (uint, error)
	RelinquishCertificate(auth AuthID, args RelinquishCertificateArgs) error
	ListCertificates(auth AuthID, args ListCertificatesArgs) (*ListCertificatesResult, error)
}
