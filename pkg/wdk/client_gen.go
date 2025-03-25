// Code generated by client-gen; DO NOT EDIT.

package wdk

type WalletStorageWriterClient struct {
	client *rpcWalletStorageWriter
}

func (c *WalletStorageWriterClient) Migrate(storageName string, storageIdentityKey string) (string, error) {
	return c.client.Migrate(storageName, storageIdentityKey)
}

func (c *WalletStorageWriterClient) MakeAvailable() (*TableSettings, error) {
	return c.client.MakeAvailable()
}

func (c *WalletStorageWriterClient) FindOrInsertUser(identityKey string) (*TableUser, error) {
	return c.client.FindOrInsertUser(identityKey)
}

func (c *WalletStorageWriterClient) CreateAction(auth AuthID, args ValidCreateActionArgs) {
	c.client.CreateAction(auth, args)
}

type rpcWalletStorageWriter struct {
	Migrate func(string, string) (string, error)

	MakeAvailable func() (*TableSettings, error)

	FindOrInsertUser func(string) (*TableUser, error)

	CreateAction func(AuthID, ValidCreateActionArgs)
}

// ===== TODO: REMOVE BELOW LINES FROM TEMPLATE =====
// ===== THIS IS JUST FOR SHOWING THE MODEL =====

// github.com/4chain-ag/go-wallet-toolbox/pkg/wdk
// Is in the same package: true

type InterfaceInfo struct {
	Name    string
	Methods []MethodInfo
}

type MethodInfo struct {
	Name      string
	Arguments []ParamInfo
	Results   []TypeInfo
}

type ParamInfo struct {
	Name string
	Type string
}

type TypeInfo struct {
	Type string
}

var Interfaces = []InterfaceInfo{
	{
		Name: "WalletStorageWriter",
		Methods: []MethodInfo{
			{
				Name: "Migrate",
				Arguments: []ParamInfo{
					{Name: "storageName", Type: "string"},
					{Name: "storageIdentityKey", Type: "string"},
				},
				Results: []TypeInfo{
					{Type: "string"},
					{Type: "error"},
				},
			},
			{
				Name:      "MakeAvailable",
				Arguments: []ParamInfo{},
				Results: []TypeInfo{
					{Type: "*TableSettings"},
					{Type: "error"},
				},
			},
			{
				Name: "FindOrInsertUser",
				Arguments: []ParamInfo{
					{Name: "identityKey", Type: "string"},
				},
				Results: []TypeInfo{
					{Type: "*TableUser"},
					{Type: "error"},
				},
			},
			{
				Name: "CreateAction",
				Arguments: []ParamInfo{
					{Name: "auth", Type: "AuthID"},
					{Name: "args", Type: "ValidCreateActionArgs"},
				},
				Results: []TypeInfo{},
			},
		},
	},
}
