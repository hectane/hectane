package db

import (
	"time"

	"github.com/jinzhu/gorm"
)

// QueueItem represents a message ready for delivery to a specific host. The
// message may need to be delivered to more than one recipient on the host. A
// comma is used to delimit the addresses. When a sender is ready to send the
// item, it sets the lock expiry time (preventing other senders from attempting
// to send it). If the attempt fails, Attempts is incremented and the lock
// expiry time is reset.
type QueueItem struct {
	ID          int64
	Time        time.Time `gorm:"not null"`
	LockExpires time.Time
	NextAttempt time.Time `gorm:"not null"`
	Attempts    int
	Host        string `gorm:"type:varchar(80);not null"`
	From        string `gorm:"type:varchar(200);not null"`
	To          string `gorm:"not null"`
	User        *User  `gorm:"ForeignKey:UserID"`
	UserID      int64
}

// GetQueueItem attempts to retrieve an item from the queue for the specified
// host (if nonzero) or any host for which an unlocked item exists. If there
// are no available items, both return values will be nil.
func GetQueueItem(host string, lockDuration time.Duration) (*QueueItem, error) {
	i := &QueueItem{}
	err := Transaction(C, func(c *gorm.DB) error {
		var (
			now = time.Now()
			d   = c.
				Where("lock_expires <= ? OR lock_expires IS NULL", now).
				Where("next_attempt >= ?", now)
		)
		if len(host) != 0 {
			d = d.Where("host = ?", host)
		}
		if err := c.First(i).Error; err != nil {
			return nil
		}
		i.LockExpires = time.Now().Add(lockDuration)
		return d.Save(i).Error
	})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return i, nil
}

// Finished indicates that a sender is finished with an item. If the maximum
// number of attempts has been exceeded, it is simply deleted. Otherwise, the
// number of attempts is incremented and the lock reset.
func (q *QueueItem) Finished(maxAttempts int) error {
	q.Attempts++
	if q.Attempts == maxAttempts {
		return C.Delete(q).Error
	}
	q.LockExpires = time.Time{}
	q.NextAttempt = time.Now().Add(time.Minute)
	return C.Save(q).Error
}
