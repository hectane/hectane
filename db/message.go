package db

import (
	"github.com/jinzhu/gorm"
)

// Message represents a single email message.
type Message struct {
	gorm.Model
	From           string `gorm:"type:varchar(200)"`
	To             string `gorm:"type:varchar(200)"`
	Subject        string `gorm:"type:varchar(200)"`
	IsUnread       bool
	HasAttachments bool
	User           *User `gorm:"ForeignKey:UserID"`
	UserID         uint
	Folder         *Folder `gorm:"ForeignKey:FolderID"`
	FolderID       uint
}
