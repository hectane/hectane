package db

import (
	"time"

	"github.com/jinzhu/gorm"
)

// QueueItem represents a message ready for delivery to a specific host.
type QueueItem struct {
	ID     int64
	Time   time.Time `gorm:"not null"`
	From   string    `gorm:"type:varchar(200);not null"`
	To     string    `gorm:"not null"`
	Host   *Host     `gorm:"Foreignkey:HostID"`
	HostID int64
	User   *User `gorm:"ForeignKey:UserID"`
	UserID int64
}

// GetQueueItem attempts to retrieve an available item from the queue for the
// specified host. Both return values are nil if no messages remain.
func GetQueueItem(hostID int64) (*QueueItem, error) {
	i := &QueueItem{}
	if err := C.Where("host_id = ?", hostID).First(i).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return i, nil
}
