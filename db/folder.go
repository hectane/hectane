package db

import (
	"database/sql"
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

func migrateFoldersTable(t *Token) error {
	_, err := t.exec(
		`
CREATE TABLE IF NOT EXISTS Folders (
	ID     SERIAL PRIMARY KEY,
	Name   VARCHAR(40) NOT NULL,
	UserID INTEGER REFERENCES Users (ID) ON DELETE CASCADE,
    UNIQUE (Name, UserID)
)
		`,
	)
	return err
}

func rowsToFolders(r *sql.Rows) ([]*Folder, error) {
	folders := make([]*Folder, 0, 1)
	for r.Next() {
		f := &Folder{}
		if err := r.Scan(&f.ID, &f.Name, &f.UserID); err != nil {
			return nil, err
		}
		folders = append(folders, f)
	}
	return folders, nil
}

// Folders retrieves all folders for the specified user.
func Folders(t *Token, userID int) ([]*Folder, error) {
	r, err := t.query(
		`
SELECT ID, Name, UserID
FROM Folders WHERE UserID = $1 ORDER BY Name
		`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	return rowsToFolders(r)
}

// FindAccount attempts to find a folder for the specified user matching the
// specified ID.
func FindFolder(t *Token, id int, userID int) (*Folder, error) {
	r, err := t.query(
		`
SELECT ID, Name, UserID
FROM Folders WHERE ID = $1 AND UserID = $2
		`,
		id,
		userID,
	)
	if err != nil {
		return nil, err
	}
	f, err := rowsToFolders(r)
	if err != nil {
		return nil, err
	}
	if len(f) != 1 {
		return nil, ErrRowCount
	}
	return f[0], nil
}

// Save persists changes to the folder. If ID is set to zero, a new folder is
// created and its ID updated.
func (f *Folder) Save(t *Token) error {
	if f.ID == 0 {
		err := t.queryRow(
			`
INSERT INTO Folders (Name, UserID)
VALUES ($1, $2) RETURNING ID
			`,
			f.Name,
			f.UserID,
		).Scan(&f.ID)
		if err != nil {
			return err
		}
		return nil
	} else {
		_, err := t.exec(
			`
UPDATE Folders SET Name=$1, UserID=$2
WHERE ID = $3
			`,
			f.Name,
			f.UserID,
			f.ID,
		)
		return err
	}
}

// Delete the folder from the database.
func (f *Folder) Delete(t *Token) error {
	_, err := t.exec(
		`
DELETE FROM Folders WHERE ID = $1
        `,
		f.ID,
	)
	return err
}
