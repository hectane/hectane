package receiver

import (
	"github.com/hectane/go-smtpsrv"
	"github.com/hectane/hectane/db"
)

// deliver receives a message, finds the correct account for it, and stores it.
func (r *Receiver) deliver(msg *smtpsrv.Message) {
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
				w, err := r.storage.CreateWriter(m.ID)
				if err != nil {
					return err
				}
				defer w.Close()
				if _, err := w.Write([]byte(msg.Body)); err != nil {
					return err
				}
				return nil
			}()
		)
		if err != nil {
			r.log.Error(err.Error())
			c.Rollback()
			continue
		}
		c.Commit()
	}
}
