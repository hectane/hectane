package storage

import (
	"io"
	"os"
	"path"
	"strconv"
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
	return path.Join(s.directory, block, strconv.FormatInt(id, 10))
}

// CreateReader attempts to open a reader for the item with the specified ID in
// the specified block.
func (s *Storage) CreateReader(block string, id int64) (io.ReadCloser, error) {
	return os.Open(s.filename(block, id))
}

// CreateWriter attempts to open a write for the item with the specified ID in
// the specified block.
func (s *Storage) CreateWriter(block string, id int64) (io.WriteCloser, error) {
	if err := os.MkdirAll(path.Join(s.directory, block), 0700); err != nil {
		return nil, err
	}
	return os.Create(s.filename(block, id))
}
