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

	DerivationPrefix string
	DerivationSuffix string

	BasketID int           `gorm:"index"`
	Basket   *OutputBasket `gorm:"foreignKey:BasketID"`

	Spendable bool
	Change    bool

	Description string
	ProvidedBy  string
	Purpose     string
	Type        string
}
