package db

import (
	"time"

	"github.com/jinzhu/gorm"
)

const (
	lockDuration = 2 * time.Minute
	maxAttempts  = 8
)

// Host represents an individual host for message delivery. A queue can "lock"
// a host for a period of time, ensuring no other queues will attempt to send
// messages to the host during that period of time. The Attempts field
// indicates how many consecutive failed attempts have been made to deliver a
// message to the host. If this number exceeds a certain threshold, all
// messages being sent to the host are deleted.
type Host struct {
	ID          int64
	Name        string `gorm:"not null"`
	LockExpires time.Time
	NextAttempt time.Time
	Attempts    int
}

// GetHost retrieves the host by name, creating it if it does not exist.
func GetHost(host string) (*Host, error) {
	h := &Host{
		Name: host,
	}
	if err := C.FirstOrCreate(h, h).Error; err != nil {
		return nil, err
	}
	return h, nil
}

// GetAvailableHost attempts to retrieve an available host, locking and
// returning it. Both return values are nil if no valid hosts are present. Use
// a defer statement with the Finished method to ensure the host is properly
// cleaned up.
func GetAvailableHost() (*Host, error) {
	var (
		h   = &Host{}
		err = Transaction(C, func(c *gorm.DB) error {
			var (
				now = time.Now()
				q   = c.
					Where("lock_expires <= ? OR lock_expires IS NULL", now).
					Where("next_attempt <= ? OR next_attempt IS NULL", now)
			)
			if err := q.First(h).Error; err != nil {
				return err
			}
			h.LockExpires = now.Add(lockDuration)
			return c.Save(h).Error
		})
	)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return h, nil
}

// Succeeded resets the failure columns for the host. The changes are not saved
// to the database.
func (h *Host) Succeeded() {
	h.NextAttempt = time.Time{}
	h.Attempts = 0
}

// Failed increments the number of attempts and calculates when the next
// attempt should be made. The changes are not saved to the database.
func (h *Host) Failed() {
	h.NextAttempt = time.Now().Add(time.Duration(2^h.Attempts) * time.Minute)
	h.Attempts++
}

// Renew obtains a new lock for the host.
func (h *Host) Renew() error {
	h.LockExpires = time.Now().Add(lockDuration)
	return C.Save(h).Error
}

// Finished indicates that a queue is done with the host. If the maximum number
// of attempts have been reached.
func (h *Host) Finished() error {
	return Transaction(C, func(c *gorm.DB) error {
		if h.Attempts > maxAttempts {
			if err := c.Delete(&Message{}, "host_id = ?", h.ID).Error; err != nil {
				return err
			}
		} else {
			var (
				count int
				q     = c.
					Table("messages").
					Where("host_id = ?", h.ID).
					Count(&count)
			)
			if err := q.Error; err != nil {
				return err
			}
			if count != 0 {
				h.LockExpires = time.Time{}
				return c.Save(h).Error
			}
		}
		return c.Delete(h).Error
	})
}
