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
