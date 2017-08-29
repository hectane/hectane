package db

import (
	"strconv"
)

// Folder provides a means to organize email messages.
type Folder struct {
	ID     int64  `json:"-"`
	Name   string `json:"name" gorm:"type:varchar(40);not null"`
	User   *User  `gorm:"ForeignKey:UserID"`
	UserID int64
}

// GetID retrieves the ID of the folder.
func (f *Folder) GetID() string {
	return strconv.FormatInt(f.ID, 10)
}

// SetID sets the ID for the folder.
func (f *Folder) SetID(id string) error {
	f.ID, _ = strconv.ParseInt(id, 10, 64)
	return nil
}
