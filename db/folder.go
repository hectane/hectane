package db

import (
	"github.com/jinzhu/gorm"
)

// Folder provides a means to organize email messages.
type Folder struct {
	gorm.Model
	Name   string `gorm:"type:varchar(40);not null"`
	User   *User  `gorm:"ForeignKey:UserID"`
	UserID uint
}
