package db

import (
	"strconv"
	"time"

	"github.com/emersion/go-imap"
)

// Message represents a single email message.
type Message struct {
	ID             int64     `json:"-"`
	Time           time.Time `json:"time"`
	From           string    `json:"from" gorm:"type:varchar(200)"`
	To             string    `json:"to" gorm:"type:varchar(200)"`
	Subject        string    `json:"subject" gorm:"type:varchar(200)"`
	IsSeen         bool      `json:"is-seen"`
	IsAnswered     bool      `json:"is-answered"`
	IsFlagged      bool      `json:"is-flagged"`
	IsDeleted      bool      `json:"is-deleted"`
	IsDraft        bool      `json:"is-draft"`
	IsRecent       bool      `json:"is-recent"`
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

// GetFlags retrieves the list of flags for the message.
func (m *Message) GetFlags() []string {
	flags := []string{}
	if m.IsSeen {
		flags = append(flags, imap.SeenFlag)
	}
	if m.IsAnswered {
		flags = append(flags, imap.AnsweredFlag)
	}
	if m.IsFlagged {
		flags = append(flags, imap.FlaggedFlag)
	}
	if m.IsDeleted {
		flags = append(flags, imap.DeletedFlag)
	}
	if m.IsDraft {
		flags = append(flags, imap.DraftFlag)
	}
	if m.IsRecent {
		flags = append(flags, imap.RecentFlag)
	}
	return flags
}

// SetFlags sets the list of flags for the message.
func (m *Message) SetFlags(flags []string) {
	m.IsSeen = false
	m.IsAnswered = false
	m.IsFlagged = false
	m.IsDeleted = false
	m.IsDraft = false
	m.IsRecent = false
	for _, f := range flags {
		switch f {
		case imap.SeenFlag:
			m.IsSeen = true
		case imap.AnsweredFlag:
			m.IsAnswered = true
		case imap.FlaggedFlag:
			m.IsFlagged = true
		case imap.DeletedFlag:
			m.IsDeleted = true
		case imap.DraftFlag:
			m.IsDraft = true
		case imap.RecentFlag:
			m.IsRecent = true
		}
	}
}
