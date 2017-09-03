package storage

import (
	"io"
	"os"
	"path"
	"strconv"
)

const (
	BlockMailbox = "mailbox"
	BlockQueue   = "queue"
)

// Storage brokers access to the on-disk storage for email message bodies.
type Storage struct {
	directory string
}

// New creates a new storage backend using the specified configuration.
func New(cfg *Config) *Storage {
	return &Storage{
		directory: cfg.Directory,
	}
}

func (s *Storage) filename(block string, id int64) string {
	return path.Join(s.directory, strconv.FormatInt(id, 10))
}

// CreateReader attempts to open a reader for the item with the specified ID.
func (s *Storage) CreateReader(block string, id int64) (io.ReadCloser, error) {
	return os.Open(s.filename(block, id))
}

// CreateWriter attempts to open a write for the item with the specified ID.
func (s *Storage) CreateWriter(block string, id int64) (io.WriteCloser, error) {
	if err := os.MkdirAll(path.Join(s.directory, block), 0700); err != nil {
		return nil, err
	}
	return os.Create(s.filename(block, id))
}

// GetSize attempts to retrieve the size of the item.
func (s *Storage) GetSize(block string, id int64) (int64, error) {
	i, err := os.Stat(s.filename(block, id))
	if err != nil {
		return 0, err
	}
	return i.Size(), nil
}
