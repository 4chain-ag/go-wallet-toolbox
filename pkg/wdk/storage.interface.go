package wdk

//go:generate go run -tags gen ../../tools/client-gen/main.go -out client_gen.go

// ValidCreateActionInput represents the input for a transaction action
type ValidCreateActionInput struct {
	Outpoint              OutPoint
	InputDescription      DescriptionString5to2000Bytes
	SequenceNumber        PositiveIntegerOrZero
	UnlockingScript       *HexString
	UnlockingScriptLength PositiveInteger
}

// ValidCreateActionOutput represents the output for a transaction action
type ValidCreateActionOutput struct {
	LockingScript      HexString
	Satoshis           SatoshiValue
	OutputDescription  DescriptionString5to2000Bytes
	Basket             *BasketStringUnder300Bytes
	CustomInstructions *string
	Tags               []BasketStringUnder300Bytes
}

// ValidProcessActionOptions represents options for processing an action
type ValidProcessActionOptions struct {
	AcceptDelayedBroadcast BooleanDefaultTrue
	ReturnTXIDOnly         BooleanDefaultFalse
	NoSend                 BooleanDefaultFalse
	SendWith               []TXIDHexString
}

// ValidCreateActionOptions extends ValidProcessActionOptions with additional options
type ValidCreateActionOptions struct {
	ValidProcessActionOptions
	SignAndProcess   bool
	TrustSelf        *string
	KnownTxids       []TXIDHexString
	NoSendChange     []OutPoint
	RandomizeOutputs bool
}

// ValidProcessActionArgs represents arguments for processing an action.
// It contains the core parameters needed to process a transaction.
type ValidProcessActionArgs struct {
	// Options contains configuration settings for how the action should be processed
	Options ValidProcessActionOptions
	// IsSendWith is true if a batch of transactions is included for processing
	IsSendWith bool
	// IsNewTx is true if there is a new transaction (not no inputs and no outputs)
	IsNewTx bool
	// IsRemixChange is true if this is a request to remix change
	// When true, IsNewTx will also be true and IsSendWith must be false
	IsRemixChange bool
	// IsNoSend is true if any new transaction should NOT be sent to the network
	IsNoSend bool
	// IsDelayed is true if options.AcceptDelayedBroadcast is true
	IsDelayed bool
}

// ValidCreateActionArgs represents the arguments for creating a transaction action
type ValidCreateActionArgs struct {
	Description                  DescriptionString5to2000Bytes `json:"description,omitempty"`
	InputBEEF                    BEEF                          `json:"input_beef,omitempty"`
	Inputs                       []ValidCreateActionInput      `json:"inputs,omitempty"`
	Outputs                      []ValidCreateActionOutput     `json:"outputs,omitempty"`
	LockTime                     int                           `json:"lock_time,omitempty"`
	Version                      int                           `json:"version,omitempty"`
	Labels                       []string                      `json:"labels,omitempty"`
	IsSignAction                 bool                          `json:"is_sign_action,omitempty"`
	RandomVals                   *[]int                        `json:"random_vals,omitempty"`
	IncludeAllSourceTransactions bool                          `json:"include_all_source_transactions,omitempty"`

	// Below are args from ValidProcessActionArgs

	// IsSendWith is true if a batch of transactions is included for processing
	IsSendWith bool `json:"is_send_with,omitempty"`
	// IsNewTx is true if there is a new transaction (not no inputs and no outputs)
	IsNewTx bool `json:"is_new_tx,omitempty"`
	// IsRemixChange is true if this is a request to remix change
	// When true, IsNewTx will also be true and IsSendWith must be false
	IsRemixChange bool `json:"is_remix_change,omitempty"`
	// IsNoSend is true if any new transaction should NOT be sent to the network
	IsNoSend bool `json:"is_no_send,omitempty"`
	// IsDelayed is true if options.AcceptDelayedBroadcast is true
	IsDelayed bool `json:"is_delayed,omitempty"`
}

// WalletStorageWriter is an interface for writing to the wallet storage
type WalletStorageWriter interface {
	CreateAction(auth AuthID, args ValidCreateActionArgs)
}
