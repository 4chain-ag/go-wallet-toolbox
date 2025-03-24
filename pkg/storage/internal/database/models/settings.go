package models

import "time"

// Settings is the database model of the settings
type Settings struct {
	CreatedAt time.Time
	UpdatedAt time.Time

	StorageIdentityKey string `gorm:"primaryKey;type:varchar(130);not null"`
	StorageName        string `gorm:"type:varchar(128);not null"`
	Chain              string `gorm:"type:varchar(10);not null"`
	MaxOutputScript    int    `gorm:"not null"`
}
