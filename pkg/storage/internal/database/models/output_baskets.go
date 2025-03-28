package models

import (
	"time"

	"gorm.io/gorm"
)

// OutputBaskets is the database model of the output baskets
type OutputBaskets struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	BasketID                int    `gorm:"primaryKey;not null"`
	Name                    string `gorm:"type:varchar(300);not null;uniqueIndex:idx_name_user_id"`
	NumberOfDesiredUTXOs    int    `gorm:"not null;default:6"`
	MinimumDesiredUTXOValue int    `gorm:"not null;default:10000"`

	UserID int `gorm:"uniqueIndex:idx_name_user_id"`
}
