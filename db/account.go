package db

import (
	"strconv"
)

// Account represents an individual email account owned by a user.
type Account struct {
	ID       int64  `json:"-"`
	Name     string `json:"name" gorm:"type:varchar(40);not null"`
	User     *User  `gorm:"ForeignKey:UserID"`
	UserID   int64
	Domain   *Domain `gorm:"ForeignKey:DomainID"`
	DomainID int64
}

// GetID retrieves the ID of the account.
func (a *Account) GetID() string {
	return strconv.FormatInt(a.ID, 10)
}

// SetID sets the ID for the account.
func (a *Account) SetID(id string) error {
	a.ID, _ = strconv.ParseInt(id, 10, 64)
	return nil
}
