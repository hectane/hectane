package imap

import (
	"errors"
	"io"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/backend/backendutil"
	"github.com/hectane/hectane/db"
	"github.com/jinzhu/gorm"
)

var ErrUnimplemented = errors.New("not yet implemented")

// mailbox maintains information about a specific folder.
type mailbox struct {
	imap   *IMAP
	folder *db.Folder
}

func (m *mailbox) count(unseen bool) (uint32, error) {
	var (
		count uint32
		c     = db.C.
			Model(&db.Message{}).
			Where("folder_id = ?", m.folder.ID).
			Count(&count)
	)
	if unseen {
		c = c.Where("is_seen = ?", false)
	}
	if err := c.Error; err != nil {
		return 0, err
	}
	return count, nil
}

// Name retrieves the name of the folder.
func (m *mailbox) Name() string {
	return m.folder.Name
}

// Info retrieves information about the folder.
func (m *mailbox) Info() (*imap.MailboxInfo, error) {
	return &imap.MailboxInfo{
		Attributes: []string{},
		Name:       m.folder.Name,
	}, nil
}

// TODO

// Status retrieves information about the messages in the folder.
func (m *mailbox) Status(items []string) (*imap.MailboxStatus, error) {
	s := imap.NewMailboxStatus(m.folder.Name, items)
	s.Flags = []string{}
	s.PermanentFlags = []string{"\\*"}
	for _, item := range items {
		switch item {
		case imap.MailboxMessages, imap.MailboxUnseen:
			c, err := m.count(item == imap.MailboxUnseen)
			if err != nil {
				return nil, err
			}
			s.Messages = c
		case imap.MailboxRecent:
			s.Recent = 0
		case imap.MailboxUidNext:
			s.UidNext = 1
		case imap.MailboxUidValidity:
			s.UidValidity = 1
		}
	}
	return s, nil
}

// SetSubscribed is unimplemented.
func (m *mailbox) SetSubscribed(subscribed bool) error {
	return ErrUnimplemented
}

// Check doesn't do anything.
func (m *mailbox) Check() error {
	return nil
}

// List messages retrieves all of the requested messages in the folder.
func (m *mailbox) ListMessages(uid bool, seqset *imap.SeqSet, items []string, ch chan<- *imap.Message) error {
	defer close(ch)
	return m.walk(uid, seqset, func(seqNum uint32, msg *db.Message) error {
		n, err := m.message(msg, seqNum, items)
		if err != nil {
			return err
		}
		ch <- n
		return nil
	})
}

// SearchMessages is unimplemented.
func (m *mailbox) SearchMessages(uid bool, criteria *imap.SearchCriteria) ([]uint32, error) {
	return nil, ErrUnimplemented
}

// CreateMessage is unimplemented.
func (m *mailbox) CreateMessage(flags []string, date time.Time, body imap.Literal) error {
	return ErrUnimplemented
}

// UpdateMessagesFlags updates the flags for the message.
func (m *mailbox) UpdateMessagesFlags(uid bool, seqset *imap.SeqSet, op imap.FlagsOp, flags []string) error {
	return m.walk(uid, seqset, func(seqNum uint32, msg *db.Message) error {
		msg.SetFlags(backendutil.UpdateFlags(msg.GetFlags(), op, flags))
		return db.C.Save(msg).Error
	})
}

// CopyMessages copies the specified messages to a new folder.
func (m *mailbox) CopyMessages(uid bool, seqset *imap.SeqSet, dest string) error {
	return db.Transaction(db.C, func(c *gorm.DB) error {
		f, err := db.GetFolder(c, m.folder.UserID, dest, false)
		if err != nil {
			return err
		}
		return m.walk(uid, seqset, func(seqNum uint32, msg *db.Message) error {
			r, err := m.imap.storage.CreateReader(msg.ID)
			if err != nil {
				return err
			}
			defer r.Close()
			msg.ID = 0
			msg.FolderID = f.ID
			if err := c.Save(msg).Error; err != nil {
				return err
			}
			w, err := m.imap.storage.CreateWriter(msg.ID)
			if err != nil {
				return err
			}
			defer w.Close()
			_, err = io.Copy(w, r)
			return err
		})
	})
}

// Expunge is unimplemented.
func (m *mailbox) Expunge() error {
	return ErrUnimplemented
}
