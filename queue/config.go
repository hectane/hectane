package queue

import (
	"github.com/hectane/hectane/storage"
)

type Config struct {
	Addr    string
	Storage *storage.Storage
}
