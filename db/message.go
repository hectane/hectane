package db

import (
	"time"

	"github.com/hectane/hectane/db/util"
)

// Message represents a single email message.
type Message struct {
	ID             int       `json:"-"`
	Time           time.Time `json:"time"`
	From           string    `json:"from"`
	To             string    `json:"to"`
	Subject        string    `json:"subject"`
	IsUnread       bool      `json:"is_unread"`
	HasAttachments bool      `json:"has_attachments"`
	UserID         int       `json:"user_id"`
	FolderID       int       `json:"folder_id"`
}

func migrateMessageTable(t *util.Token) error {
	_, err := t.Exec(
		`
CREATE TABLE IF NOT EXISTS Message (
	ID             SERIAL PRIMARY KEY,
	Time           TIMESTAMPTZ NOT NULL,
	From_          VARCHAR(200),
	To_            VARCHAR(200),
	Subject        VARCHAR(200),
	IsUnread       BOOLEAN,
	HasAttachments BOOLEAN,
	UserID         INTEGER REFERENCES Users (ID) ON DELETE CASCADE,
	FolderID       INTEGER REFERENCES Folders (ID) ON DELETE CASCADE
)
		`,
	)
	return err
}
