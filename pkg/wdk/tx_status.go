package wdk

// TxStatus Transaction status stored in database
type TxStatus string

// Possible transaction statuses stored in database
const (
	TxStatusCompleted   TxStatus = "completed"
	TxStatusFailed      TxStatus = "failed"
	TxStatusUnprocessed TxStatus = "unprocessed"
	TxStatusSending     TxStatus = "sending"
	TxStatusUnproven    TxStatus = "unproven"
	TxStatusUnsigned    TxStatus = "unsigned"
	TxStatusNoSend      TxStatus = "nosend"
	TxStatusNonFinal    TxStatus = "nonfinal"
	TxStatusUnfail      TxStatus = "unfail"
)

// ProvenTxReqStatus represents the status of a proven transaction in a defined processing state as a string.
type ProvenTxReqStatus string

// Possible proven transaction statuses stored in database
const (
	ProvenTxStatusSending     ProvenTxReqStatus = "sending"
	ProvenTxStatusUnsent      ProvenTxReqStatus = "unsent"
	ProvenTxStatusNoSend      ProvenTxReqStatus = "nosend"
	ProvenTxStatusUnknown     ProvenTxReqStatus = "unknown"
	ProvenTxStatusNonFinal    ProvenTxReqStatus = "nonfinal"
	ProvenTxStatusUnprocessed ProvenTxReqStatus = "unprocessed"
	ProvenTxStatusUnmined     ProvenTxReqStatus = "unmined"
	ProvenTxStatusCallback    ProvenTxReqStatus = "callback"
	ProvenTxStatusUnconfirmed ProvenTxReqStatus = "unconfirmed"
	ProvenTxStatusCompleted   ProvenTxReqStatus = "completed"
	ProvenTxStatusInvalid     ProvenTxReqStatus = "invalid"
	ProvenTxStatusDoubleSpend ProvenTxReqStatus = "doubleSpend"
	ProvenTxStatusUnfail      ProvenTxReqStatus = "unfail"
)

// TxReqBroadcastStatus is a reduced ProvenTxReqStatus, used to decide whether to broadcast a transaction or not.
type TxReqBroadcastStatus string

// Possible transaction request broadcast statuses
const (
	TxReqSimplifiedReadyToSend TxReqBroadcastStatus = "readyToSend"
	TxReqSimplifiedAlreadySent TxReqBroadcastStatus = "alreadySent"
	TxReqSimplifiedError       TxReqBroadcastStatus = "error"
	TxReqSimplifiedUnknown     TxReqBroadcastStatus = "unknown"
)

// BroadcastStatus returns the simplified broadcast status of a transaction request based on its current status.
func (s ProvenTxReqStatus) BroadcastStatus() TxReqBroadcastStatus {
	switch s {
	case ProvenTxStatusUnknown,
		ProvenTxStatusNonFinal,
		ProvenTxStatusInvalid,
		ProvenTxStatusDoubleSpend:
		return TxReqSimplifiedError
	case ProvenTxStatusSending,
		ProvenTxStatusUnsent,
		ProvenTxStatusNoSend,
		ProvenTxStatusUnprocessed:
		return TxReqSimplifiedReadyToSend
	case ProvenTxStatusUnmined,
		ProvenTxStatusCallback,
		ProvenTxStatusUnconfirmed,
		ProvenTxStatusCompleted:
		return TxReqSimplifiedAlreadySent
	case ProvenTxStatusUnfail:
		fallthrough
	default:
		return TxReqSimplifiedUnknown
	}
}
