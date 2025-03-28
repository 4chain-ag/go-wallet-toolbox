package wdk

// ValidCreateActionInput represents the input for a transaction action
type ValidCreateActionInput struct {
	Outpoint              OutPoint                      `json:"outpoint,omitempty"`
	InputDescription      DescriptionString5to2000Bytes `json:"inputDescription,omitempty"`
	SequenceNumber        PositiveIntegerOrZero         `json:"sequenceNumber,omitempty"`
	UnlockingScript       *HexString                    `json:"unlockingScript,omitempty"`
	UnlockingScriptLength *PositiveInteger              `json:"unlockingScriptLength,omitempty"`
}

// ValidCreateActionOutput represents the output for a transaction action
type ValidCreateActionOutput struct {
	LockingScript      HexString                     `json:"lockingScript,omitempty"`
	Satoshis           SatoshiValue                  `json:"satoshis,omitempty"`
	OutputDescription  DescriptionString5to2000Bytes `json:"outputDescription,omitempty"`
	Basket             *BasketStringUnder300Bytes    `json:"basket,omitempty"`
	CustomInstructions *string                       `json:"customInstructions,omitempty"`
	Tags               []BasketStringUnder300Bytes   `json:"tags,omitempty"`
}

// ValidProcessActionOptions represents options for processing an action
type ValidProcessActionOptions struct {
	AcceptDelayedBroadcast BooleanDefaultTrue  `json:"acceptDelayedBroadcast,omitempty"`
	ReturnTXIDOnly         BooleanDefaultFalse `json:"returnTXIDOnly,omitempty"`
	NoSend                 BooleanDefaultFalse `json:"noSend,omitempty"`
	SendWith               []TXIDHexString     `json:"sendWith,omitempty"`
}

// ValidCreateActionOptions extends ValidProcessActionOptions with additional options
type ValidCreateActionOptions struct {
	ValidProcessActionOptions `json:",inline"`
	SignAndProcess            bool            `json:"signAndProcess,omitempty"`
	TrustSelf                 *string         `json:"trustSelf,omitempty"`
	KnownTxids                []TXIDHexString `json:"knownTxids,omitempty"`
	NoSendChange              []OutPoint      `json:"noSendChange,omitempty"`
	RandomizeOutputs          bool            `json:"randomizeOutputs,omitempty"`
}

// ValidProcessActionArgs represents arguments for processing an action.
// It contains the core parameters needed to process a transaction.
type ValidProcessActionArgs struct {
	// Options contains configuration settings for how the action should be processed
	Options ValidProcessActionOptions `json:"options,omitempty"`
	// IsSendWith is true if a batch of transactions is included for processing
	IsSendWith bool `json:"isSendWith,omitempty"`
	// IsNewTx is true if there is a new transaction (not no inputs and no outputs)
	IsNewTx bool `json:"isNewTx,omitempty"`
	// IsRemixChange is true if this is a request to remix change
	// When true, IsNewTx will also be true and IsSendWith must be false
	IsRemixChange bool `json:"isRemixChange,omitempty"`
	// IsNoSend is true if any new transaction should NOT be sent to the network
	IsNoSend bool `json:"isNoSend,omitempty"`
	// IsDelayed is true if options.AcceptDelayedBroadcast is true
	IsDelayed bool `json:"isDelayed,omitempty"`
}

// ValidCreateActionArgs represents the arguments for creating a transaction action
type ValidCreateActionArgs struct {
	Description                  DescriptionString5to2000Bytes `json:"description,omitempty"`
	InputBEEF                    BEEF                          `json:"inputBEEF,omitempty"`
	Inputs                       []ValidCreateActionInput      `json:"inputs,omitempty"`
	Outputs                      []ValidCreateActionOutput     `json:"outputs,omitempty"`
	LockTime                     int                           `json:"lockTime,omitempty"`
	Version                      int                           `json:"version,omitempty"`
	Labels                       []string                      `json:"labels,omitempty"`
	IsSignAction                 bool                          `json:"isSignAction,omitempty"`
	RandomVals                   *[]int                        `json:"randomVals,omitempty"`
	IncludeAllSourceTransactions bool                          `json:"includeAllSourceTransactions,omitempty"`

	Options ValidCreateActionOptions `json:"options,omitempty"`

	// Below are args from ValidProcessActionArgs

	// IsSendWith is true if a batch of transactions is included for processing
	IsSendWith bool `json:"isSendWith,omitempty"`
	// IsNewTx is true if there is a new transaction (not no inputs and no outputs)
	IsNewTx bool `json:"isNewTx,omitempty"`
	// IsRemixChange is true if this is a request to remix change
	// When true, IsNewTx will also be true and IsSendWith must be false
	IsRemixChange bool `json:"isRemixChange,omitempty"`
	// IsNoSend is true if any new transaction should NOT be sent to the network
	IsNoSend bool `json:"isNoSend,omitempty"`
	// IsDelayed is true if options.AcceptDelayedBroadcast is true
	IsDelayed bool `json:"isDelayed,omitempty"`
}
