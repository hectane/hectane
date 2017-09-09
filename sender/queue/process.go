package queue

import (
	"context"

	"github.com/hectane/hectane/db"
)

// process retrieves the next available host with messages pending and attempts
// to connect to one of its mail servers. Once a connection is established,
// messages are sent until either no more remain or an error occurs.
func (q *Queue) process(ctx context.Context) error {
	h, err := db.GetAvailableHost()
	if err != nil {
		return err
	}
	if h == nil {
		return nil
	}
	defer h.Finished()
	s, err := findMailServers(ctx, h.Name)
	if err != nil {
		h.Failed()
		return err
	}
	c, err := tryMailServers(ctx, s)
	if err != nil {
		h.Failed()
		return err
	}
	defer c.Close()
	if err := initServer(c); err != nil {
		h.Failed()
		return err
	}
	for {
		i, err := db.GetQueueItem(h.ID)
		if err != nil {
			return err
		}
		if i == nil {
			return nil
		}
		if err := q.deliver(c, i); err != nil {
			if isHostError(err) {
				h.Failed()
				return err
			}
		}
		h.Succeeded()
		if err := db.C.Delete(i).Error; err != nil {
			return err
		}
	}
}
