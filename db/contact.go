package db

import (
	"strconv"
)

// Contact represents an email address stored by a user.
type Contact struct {
	ID     int64  `json:"-"`
	Name   string `json:"name" gorm:"type:varchar(200)"`
	Email  string `json:"email" gorm:"type:varchar(200);not null"`
	User   *User  `gorm:"ForeignKey:UserID"`
	UserID int64
}

// GetID retrieves the ID of the contact.
func (c *Contact) GetID() string {
	return strconv.FormatInt(c.ID, 10)
}

// SetID sets the ID for the contact.
func (c *Contact) SetID(id string) error {
	c.ID, _ = strconv.ParseInt(id, 10, 64)
	return nil
}
