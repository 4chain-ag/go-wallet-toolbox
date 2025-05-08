package wdk

import "github.com/4chain-ag/go-wallet-toolbox/pkg/wdk/primitives"

// ProcessActionArgs defines the arguments required to process an action with transaction and sending options.
type ProcessActionArgs struct {
	IsNewTx    bool                         `json:"isNewTx"`
	IsSendWith bool                         `json:"isSendWith"`
	IsNoSend   bool                         `json:"isNoSend"`
	IsDelayed  bool                         `json:"isDelayed"`
	Reference  *string                      `json:"reference,omitempty"`
	TxID       *primitives.TXIDHexString    `json:"txid,omitempty"`
	RawTx      primitives.ExplicitByteArray `json:"rawTx,omitempty"`
	SendWith   []string                     `json:"sendWith"`
}
