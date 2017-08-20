package db

import (
	"github.com/hectane/hectane/db/util"
)

const (
	InboxFolder = "Inbox"
	SentFolder  = "Sent"
)

// Folder provides a means to organize email messages.
type Folder struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	UserID int    `json:"user_id"`
}

func migrateFolderTable(t *util.Token) error {
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
