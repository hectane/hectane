package queue

import (
	"github.com/hectane/hectane/db"
	"github.com/hectane/hectane/storage"
)

type Config struct {
	Host            *db.Host
	Storage         *storage.Storage
	QueueFinishedCh chan string
}
