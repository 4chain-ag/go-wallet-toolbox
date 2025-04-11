package models

import "gorm.io/gorm"

type Label struct {
	gorm.Model

	Name   string `gorm:"uniqueIndex:idx_name_userid"`
	UserID int    `gorm:"uniqueIndex:idx_name_userid"`
}
