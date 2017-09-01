package receiver

import (
	"time"

	"github.com/hectane/go-smtpsrv"
	"github.com/hectane/hectane/db"
)

// parse reads the message body and creates a db.Message that can be persisted
// to the database.
func parse(m *smtpsrv.Message, addr string, userID, folderID int64) (*db.Message, error) {

	// TODO: do parsing here

	return &db.Message{
		Time:     time.Now(),
		From:     m.From,
		To:       addr,
		Subject:  "[Untitled]",
		UserID:   userID,
		FolderID: folderID,
	}, nil
}
