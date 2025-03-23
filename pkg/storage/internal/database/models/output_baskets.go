package models

import "time"

// OutputBaskets is the database model of the output baskets
type OutputBaskets struct {
	CreatedAt time.Time
	UpdatedAt time.Time

	BasketID                int    `gorm:"primaryKey;not null"`
	Name                    string `gorm:"type:varchar(300);not null"`
	NumberOfDesiredUTXOs    int    `gorm:"not null;default:6"`
	MinimumDesiredUTXOValue int    `gorm:"not null;default:10000"`
	IsDeleted               bool   `gorm:"not null;default:false"` // from-dz: probably would be better to let gorm handle if it is deleted?

	// NO MATTER THE DB IT DOESN'T CREATE FOREIGN KEY
	UserID int  `gorm:"not null"`
	User   User `gorm:"foreignKey:UserID"`
}
