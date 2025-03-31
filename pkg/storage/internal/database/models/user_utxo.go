package models

import "time"

// UserUTXO is a table holding user's Unspent Transaction Outputs (UTXOs).
type UserUTXO struct {
	UserID   int    `gorm:"primaryKey;uniqueIndex:idx_window,sort:asc,priority:1"`
	TxID     string `gorm:"primaryKey;uniqueIndex:idx_window,sort:asc,priority:4"`
	Vout     uint32 `gorm:"primaryKey;uniqueIndex:idx_window,sort:asc,priority:5"`
	Satoshis uint64
	// EstimatedInputSize is the estimated size increase when adding and unlocking this UTXO to a transaction.
	EstimatedInputSize uint64
	Basket             string
	CreatedAt          time.Time `gorm:"uniqueIndex:idx_window,sort:asc,priority:3"`
	// TouchedAt is the time when the UTXO was last touched (selected for preparing transaction outline) - used for prioritizing UTXO selection.
	TouchedAt time.Time `gorm:"uniqueIndex:idx_window,sort:asc,priority:2"`
}
