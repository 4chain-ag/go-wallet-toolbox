package models

import (
	"gorm.io/gorm"
)

type Output struct {
	gorm.Model

	UserID        int    `gorm:"index"`
	TransactionID uint   `gorm:"index"`
	Vout          uint32 `gorm:"index"`
	Satoshis      int64

	LockingScript      *string `gorm:"type:string"`
	CustomInstructions *string `gorm:"type:string"`

	DerivationPrefix *string
	DerivationSuffix *string

	BasketID *int
	Basket   *OutputBasket

	Spendable bool
	Change    bool

	Description string `gorm:"type:string"`
	ProvidedBy  string
	Purpose     string
	Type        string

	Transaction *Transaction `gorm:"foreignKey:TransactionID;references:ID"`
}
