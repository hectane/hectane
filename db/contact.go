package db

import (
	"github.com/jinzhu/gorm"
)

// Contact represents an email address stored by a user.
type Contact struct {
	gorm.Model
	Name   string `gorm:"type:varchar(200)"`
	Email  string `gorm:"type:varchar(200);not null"`
	User   *User  `gorm:"ForeignKey:UserID"`
	UserID uint
}
