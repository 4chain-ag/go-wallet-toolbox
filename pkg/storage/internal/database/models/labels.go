package models

import (
	"gorm.io/gorm"
	"time"
)

type Label struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Name   string `gorm:"primarykey"`
	UserID int    `gorm:"primarykey"`
}
