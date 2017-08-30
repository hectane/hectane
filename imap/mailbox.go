package imap

import (
	"errors"
	"time"

	"github.com/emersion/go-imap"
	"github.com/hectane/hectane/db"
)

var ErrUnimplemented = errors.New("not yet implemented")

// mailbox maintains information about a specific folder.
type mailbox struct {
	folder *db.Folder
}

// Name retrieves the name of the folder.
func (m *mailbox) Name() string {
	return m.folder.Name
}

// Info retrieves information about the folder.
func (m *mailbox) Info() (*imap.MailboxInfo, error) {
	return &imap.MailboxInfo{
		Name: m.folder.Name,
	}, nil
}

func (m *mailbox) Status(items []string) (*imap.MailboxStatus, error) {
	//...
	return nil, ErrUnimplemented
}

func (m *mailbox) SetSubscribed(subscribed bool) error {
	//...
	return ErrUnimplemented
}

func (m *mailbox) Check() error {
	//...
	return ErrUnimplemented
}

func (m *mailbox) ListMessages(uid bool, seqset *imap.SeqSet, items []string, ch chan<- *imap.Message) error {
	//...
	return ErrUnimplemented
}

func (m *mailbox) SearchMessages(uid bool, criteria *imap.SearchCriteria) ([]uint32, error) {
	//...
	return nil, ErrUnimplemented
}

func (m *mailbox) CreateMessage(flags []string, date time.Time, body imap.Literal) error {
	//...
	return ErrUnimplemented
}

func (m *mailbox) UpdateMessagesFlags(uid bool, seqset *imap.SeqSet, operation imap.FlagsOp, flags []string) error {
	//...
	return ErrUnimplemented
}

func (m *mailbox) CopyMessages(uid bool, seqset *imap.SeqSet, dest string) error {
	//...
	return ErrUnimplemented
}

func (m *mailbox) Expunge() error {
	//...
	return ErrUnimplemented
}
