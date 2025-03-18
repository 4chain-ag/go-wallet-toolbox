package models

import "time"

type Settings struct {
	CreatedAt time.Time
	UpdatedAt time.Time

	StorageIdentityKey string `gorm:"type:char(130);not null"`
	StorageName        string `gorm:"type:char(128);not null"`
	Chain              string `gorm:"type:char(10);not null"`
	DBType             string `gorm:"type:char(10);not null"`
	MaxOutputs         int32  `gorm:"not null"`
}
