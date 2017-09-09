package queue

import (
	"io"
	"net/smtp"
	"strings"

	"github.com/hectane/hectane/db"
	"github.com/hectane/hectane/storage"
)

// deliver attempts to deliver the specified item to the server. If an error
// occurs, its nature is determined.
func (q *Queue) deliver(c *smtp.Client, i *db.QueueItem) error {
	r, err := q.storage.CreateReader(storage.BlockQueue, i.ID)
	if err != nil {
		return err
	}
	defer r.Close()
	if err := c.Mail(i.From); err != nil {
		return err
	}
	for _, t := range strings.Split(i.To, ",") {
		if err := c.Rcpt(t); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	defer w.Close()
	if _, err := io.Copy(w, r); err != nil {
		return err
	}
	return nil
}
