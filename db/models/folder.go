package models

import (
	"github.com/hectane/hectane/db/sql"
)

const (
	Folders = "folders"

	InboxFolder = "Inbox"
	SentFolder  = "Sent"
)

// Folder provides a means to organize email messages.
type Folder struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	UserID int    `json:"user_id"`
}

func MigrateFolders(t *sql.Token) error {
	_, err := t.Exec(
		`
CREATE TABLE IF NOT EXISTS Folder (
	ID     SERIAL PRIMARY KEY,
	Name   VARCHAR(40) NOT NULL,
	UserID INTEGER REFERENCES Users (ID) ON DELETE CASCADE,
	UNIQUE (Name, UserID)
)
		`,
	)
	return err
}
