package db

import (
	"strconv"
)

// Message represents a single email message.
type Message struct {
	ID             int64  `json:"-"`
	From           string `json:"from" gorm:"type:varchar(200)"`
	To             string `json:"to" gorm:"type:varchar(200)"`
	Subject        string `json:"subject" gorm:"type:varchar(200)"`
	IsUnread       bool   `json:"is_unread"`
	HasAttachments bool   `json:"has_attachments"`
	User           *User  `gorm:"ForeignKey:UserID"`
	UserID         int64
	Folder         *Folder `gorm:"ForeignKey:FolderID"`
	FolderID       int64
}

// GetID retrieves the ID of the message.
func (m *Message) GetID() string {
	return strconv.FormatInt(m.ID, 10)
}

// SetID sets the ID for the message.
func (m *Message) SetID(id string) error {
	m.ID, _ = strconv.ParseInt(id, 10, 64)
	return nil
}
