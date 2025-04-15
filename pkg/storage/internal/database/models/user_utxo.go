package models

import "time"

// UserUTXO is a table holding user's Unspent Transaction Outputs (UTXOs).
type UserUTXO struct {
	UserID   int           `gorm:"primaryKey"`
	TxID     string        `gorm:"primaryKey"`
	Vout     uint32        `gorm:"primaryKey"`
	BasketID int           `gorm:"not null"`
	Basket   *OutputBasket `gorm:"foreignKey:BasketID"`
	Satoshis uint64
	// EstimatedInputSize is the estimated size increase when adding and unlocking this UTXO to a transaction.
	EstimatedInputSize uint64
	CreatedAt          time.Time
}
