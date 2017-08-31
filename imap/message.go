package imap

import (
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/backend/backendutil"
	"github.com/emersion/go-message"
	"github.com/hectane/hectane/db"
)

// message converts a message from the database to its IMAP equivalent for
// sending down the wire.
func (m *mailbox) message(msg *db.Message, seqNum uint32, items []string) (*imap.Message, error) {
	r, err := m.imap.storage.CreateReader(msg.ID)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	e, err := message.Read(r)
	if err != nil {
		return nil, err
	}
	n := imap.NewMessage(seqNum, items)
	for _, item := range items {
		switch item {
		case imap.BodyMsgAttr, imap.BodyStructureMsgAttr:
			if n.BodyStructure, err = backendutil.FetchBodyStructure(e, item == imap.BodyStructureMsgAttr); err != nil {
				return nil, err
			}
		case imap.EnvelopeMsgAttr:
			if n.Envelope, err = backendutil.FetchEnvelope(e.Header); err != nil {
				return nil, err
			}
		case imap.FlagsMsgAttr:
			if msg.IsUnread {
				n.Flags = []string{imap.MailboxUnseen}
			} else {
				n.Flags = []string{}
			}
		case imap.InternalDateMsgAttr:
			n.InternalDate = msg.Time
		case imap.SizeMsgAttr:
			s, err := m.imap.storage.GetSize(msg.ID)
			if err != nil {
				return nil, err
			}
			n.Size = uint32(s)
		case imap.UidMsgAttr:
			n.Uid = uint32(msg.ID)
		default:
			s, err := imap.ParseBodySectionName(item)
			if err != nil {
				return nil, err
			}
			if n.Body[s], err = backendutil.FetchBodySection(e, s); err != nil {
				return nil, err
			}
		}
	}
	return n, nil
}
