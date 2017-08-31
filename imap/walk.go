package imap

import (
	"github.com/emersion/go-imap"
	"github.com/hectane/hectane/db"
)

type walkFn func(uint32, *db.Message) error

// walk steps through each of the messages in the sequence set and invokes the
// specified callback.
func (m *mailbox) walk(uid bool, seqset *imap.SeqSet, fn walkFn) error {
	for _, seq := range seqset.Set {
		var (
			messages        = []*db.Message{}
			offset   uint32 = 0
			c               = db.C.
					Where("folder_id = ?", m.folder.ID).
					Order("time desc")
		)
		if uid {
			if seq.Start != 0 {
				if err := c.
					Model(&db.Message{}).
					Where("id < ?", seq.Start).
					Count(&offset).
					Error; err != nil {
					return err
				}
				c = c.Where("id >= ?", seq.Start)
			}
			if seq.Stop != 0 {
				c = c.Where("id <= ?", seq.Stop)
			}
		} else {
			if seq.Start != 0 {
				offset = seq.Start - 1
				c = c.Offset(offset)
			}
			if seq.Stop != 0 {
				c = c.Limit(seq.Stop - seq.Start + 1)
			}
		}
		if err := c.Find(&messages).Error; err != nil {
			return err
		}
		for i, msg := range messages {
			if err := fn(uint32(i)+offset+1, msg); err != nil {
				return err
			}
		}
	}
	return nil
}
