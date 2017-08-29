package db

import (
	"strconv"
	"time"
)

// Message represents a single email message.
type Message struct {
	ID             int64     `json:"-"`
	Time           time.Time `json:"time"`
	From           string    `json:"from" gorm:"type:varchar(200)"`
	To             string    `json:"to" gorm:"type:varchar(200)"`
	Subject        string    `json:"subject" gorm:"type:varchar(200)"`
	IsUnread       bool      `json:"is-unread"`
	HasAttachments bool      `json:"has-attachments"`
	User           *User     `json:"-" gorm:"ForeignKey:UserID"`
	UserID         int64     `json:"-"`
	Folder         *Folder   `json:"-" gorm:"ForeignKey:FolderID"`
	FolderID       int64     `json:"-"`
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
