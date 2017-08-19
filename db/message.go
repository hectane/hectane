package db

import (
	"database/sql"
	"time"
)

// Message represents a single email message.
type Message struct {
	ID             int       `json:"id"`
	Time           time.Time `json:"time"`
	From           string    `json:"from"`
	Subject        string    `json:"subject"`
	IsUnread       bool      `json:"is_unread"`
	HasAttachments bool      `json:"has_attachments"`
	UserID         int       `json:"user_id"`
	FolderID       int       `json:"folder_id"`
}

func migrateMessagesTable(t *Token) error {
	_, err := t.exec(
		`
CREATE TABLE IF NOT EXISTS Messages (
	ID             SERIAL PRIMARY KEY,
	Time           TIMESTAMPTZ NOT NULL,
	From           VARCHAR(200),
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

func rowsToMessages(r *sql.Rows) ([]*Message, error) {
	messages := make([]*Message, 0, 1)
	for r.Next() {
		m := &Message{}
		if err := r.Scan(&m.ID, &m.Time, &m.From, &m.Subject, &m.IsUnread, &m.HasAttachments, &m.UserID, &m.FolderID); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, nil
}

// Messages retrieves all messages for the specified user in the specified
// folder.
func Messages(t *Token, userID, folderID int) ([]*Message, error) {
	r, err := t.query(
		`
SELECT ID, Time, From, Subject, IsUnread, HasAttachment, UserID, FolderID
FROM Messages WHERE UserID = $1 AND FolderID = $2 ORDER BY Time DESC
		`,
	)
	if err != nil {
		return nil, err
	}
	return rowsToMessages(r)
}

// Save persists changes to the message. If ID is set to zero, a new message is
// created and its ID updated.
func (m *Message) Save(t *Token) error {
	if m.ID == 0 {
		return t.queryRow(
			`
INSERT INTO Messages (Time, From, Subject, IsUnread, HasAttachment, UserID, FolderID)
VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING ID
			`,
			m.Time,
			m.From,
			m.Subject,
			m.IsUnread,
			m.HasAttachments,
			m.UserID,
			m.FolderID,
		).Scan(&m.ID)
	} else {
		_, err := t.exec(
			`
UPDATE Messages SET Time=$1, From=$2, Subject=$3, IsUnread=$4, HasAttachment=$5, UserID=$6, FolderID=$7
WHERE ID = $8
			`,
			m.Time,
			m.From,
			m.Subject,
			m.IsUnread,
			m.HasAttachments,
			m.UserID,
			m.FolderID,
			m.ID,
		)
		return err
	}
}

// Delete the message from the database.
func (m *Message) Delete(t *Token) error {
	_, err := t.exec(
		`
DELETE FROM Messages WHERE ID = $1
        `,
		m.ID,
	)
	return err
}
