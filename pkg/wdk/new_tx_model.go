package wdk

type NewTxModel struct {
	UserID int

	Version     int
	LockTime    int
	Status      TxStatus
	Reference   string
	Satoshis    uint64
	IsOutgoing  bool
	InputBeef   []byte
	Description string

	Labels []string
}
