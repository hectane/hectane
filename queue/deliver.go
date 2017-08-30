package queue

import (
	"github.com/hectane/go-smtpsrv"
	"github.com/hectane/hectane/db"
)

// deliver receives a message, finds the correct account for it, and stores it.
func (q *Queue) deliver(msg *smtpsrv.Message) {
	for _, addr := range msg.To {
		var (
			c   = db.C.Begin()
			err = func() error {
				a, err := lookup(c, addr)
				if err != nil {
					return err
				}
				f, err := getFolder(c, a.UserID, db.FolderInbox)
				if err != nil {
					return err
				}
				m, err := parse(msg, addr, a.UserID, f.ID)
				if err != nil {
					return err
				}
				if err := c.Create(m).Error; err != nil {
					return err
				}
				w, err := q.storage.CreateWriter(Block, m.ID)
				if err != nil {
					return err
				}
				defer w.Close()
				if _, err := w.Write([]byte(msg.Body)); err != nil {
					return err
				}
				return nil
			}
		)
		if err != nil {
			q.log.Error(err)
			c.Rollback()
			continue
		}
		c.Commit()
	}
}
