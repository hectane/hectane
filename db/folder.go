package db

import (
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
)

const (
	FolderInbox = "Inbox"
	FolderSent  = "Sent"
)

// Folder provides a means to organize email messages.
type Folder struct {
	ID     int64  `json:"-"`
	Name   string `json:"name" gorm:"type:varchar(40);not null"`
	User   *User  `json:"-" gorm:"ForeignKey:UserID"`
	UserID int64  `json:"-"`
}

// GetFolder attempts to retrieve a folder by name. If created is set to true,
// the folder will be created before being returned. The names "INBOX" and
// "SENT" will be normalized.
func GetFolder(c *gorm.DB, userID int64, name string, create bool) (*Folder, error) {
	if strings.ToLower(name) == strings.ToLower(FolderInbox) {
		name = FolderInbox
	}
	var (
		f = &Folder{
			Name:   name,
			UserID: userID,
		}
		err = c.Where(f).First(f).Error
	)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if create {
		f.Name = name
		err = c.Save(f).Error
	}
	if err != nil {
		return nil, err
	}
	return f, nil
}

// GetFolderTransaction wraps GetFolder in a transaction.
func GetFolderTransaction(userID int64, name string, create bool) (*Folder, error) {
	var f *Folder
	if err := Transaction(C, func(c *gorm.DB) error {
		var err error
		f, err = GetFolder(c, userID, name, create)
		return err
	}); err != nil {
		return nil, err
	}
	return f, nil
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
