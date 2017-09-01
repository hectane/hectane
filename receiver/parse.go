package receiver

import (
	"strings"
	"time"

	"github.com/emersion/go-message"
	"github.com/hectane/go-smtpsrv"
	"github.com/hectane/hectane/db"
)

// TODO: check for attachments

// parse reads the message body and creates a db.Message that can be persisted
// to the database.
func parse(m *smtpsrv.Message, userID, folderID int64) (*db.Message, error) {
	r := strings.NewReader(m.Body)
	e, err := message.Read(r)
	if err != nil {
		return nil, err
	}
	return &db.Message{
		Time:     time.Now(),
		From:     m.From,
		To:       strings.Join(m.To, ", "),
		Subject:  e.Header.Get("subject"),
		UserID:   userID,
		FolderID: folderID,
	}, nil
}
