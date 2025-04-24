package wdk

// StorageCreateTransactionSdkInput represents the input for SDK transaction creation
type StorageCreateTransactionSdkInput struct {
	Vin                   int
	SourceTxid            string
	SourceVout            uint32
	SourceSatoshis        int64
	SourceLockingScript   string
	SourceTransaction     []byte
	UnlockingScriptLength int
	ProvidedBy            ProvidedBy
	Type                  string
	SpendingDescription   *string
	DerivationPrefix      *string
	DerivationSuffix      *string
	SenderIdentityKey     *string
}

// StorageCreateTransactionSdkOutput represents the output for SDK transaction creation
type StorageCreateTransactionSdkOutput struct {
	ValidCreateActionOutput
	// Additional fields
	Vout             uint32
	ProvidedBy       ProvidedBy
	Purpose          string
	DerivationSuffix *string
}

// StorageCreateActionResult represents the result of creating a transaction action
type StorageCreateActionResult struct {
	// InputBeef contains the raw binary data of the input BEEF, if any
	InputBeef *[]byte
	// Inputs is a list of transaction inputs
	Inputs []StorageCreateTransactionSdkInput
	// Outputs is a list of transaction outputs
	Outputs []StorageCreateTransactionSdkOutput
	// NoSendChangeOutputVouts contains indices of outputs that should not be sent as change
	NoSendChangeOutputVouts *[]int
	// DerivationPrefix is the prefix used for key derivation
	DerivationPrefix string
	// Version is the transaction version
	Version int
	// LockTime is the transaction lock time
	LockTime int
	// Reference is a unique identifier for this transaction
	Reference string
}
