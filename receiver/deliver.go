package receiver

import (
	"github.com/hectane/go-smtpsrv"
	"github.com/hectane/hectane/db"
	"github.com/jinzhu/gorm"
)

// deliver receives a message, finds the correct account for it, and stores it.
func (r *Receiver) deliver(msg *smtpsrv.Message) {
	for _, addr := range msg.To {
		if err := db.Transaction(db.C, func(c *gorm.DB) error {
			a, err := lookup(c, addr)
			if err != nil {
				return err
			}
			f, err := db.GetFolder(c, a.UserID, db.FolderInbox, true)
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
		}); err != nil {
			r.log.Warn(err.Error())
		}
	}
}
