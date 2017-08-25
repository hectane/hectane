package db

import (
	"github.com/jinzhu/gorm"
)

// Account represents an individual email account owned by a user.
type Account struct {
	gorm.Model
	Name     string `gorm:"type:varchar(40);not null"`
	User     *User  `gorm:"ForeignKey:UserID"`
	UserID   uint
	Domain   *Domain `gorm:"ForeignKey:DomainID"`
	DomainID uint
}
