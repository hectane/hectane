package queue

import (
	"code.google.com/p/go-uuid/uuid"

	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"sync"
)

// The item with the specified ID does not exist.
var InvalidID = errors.New("invalid ID")

// File extensions for files on disk
var (
	bodyExtension    = ".body"
	messageExtension = ".message"
)

// Message metadata.
type Message struct {
	Host string
	From string
	To   []string
	Body string
}

// Storage facilitator of email metadata and message content on disk. All
// methods are safe to call from multiple goroutines. None of the file
// operations are atomic.
type Storage struct {
	sync.Mutex
	directory string
	bodies    map[string]int
}

// Determine the absolute filename of the index file.
func (s *Storage) indexFilename() string {
	return path.Join(s.directory, "index.json")
}

// Determine the absolute filename of a body file.
func (s *Storage) bodyFilename(id string) string {
	return path.Join(s.directory, fmt.Sprintf("%s%s", id, bodyExtension))
}

// Determine the absolute filename of a message file.
func (s *Storage) messageFilename(id string) string {
	return path.Join(s.directory, fmt.Sprintf("%s%s", id, messageExtension))
}

// Attempt to load the specified item from the specified file.
func (s *Storage) loadJSON(filename string, v interface{}) error {
	if r, err := os.Open(filename); err == nil {
		if err := json.NewDecoder(r).Decode(v); err == nil {
			return r.Close()
		} else {
			return err
		}
	} else {
		return err
	}
}

// Attempt to save the specified item to the specified file.
func (s *Storage) saveJSON(filename string, v interface{}) error {
	if w, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600); err == nil {
		if err := json.NewEncoder(w).Encode(v); err == nil {
			return w.Close()
		} else {
			return err
		}
	} else {
		return err
	}
}

// Retrieve a message.
func (s *Storage) getMessage(id string) (*Message, error) {
	var m Message
	if err := s.loadJSON(s.messageFilename(id), &m); err == nil {
		return &m, nil
	} else if os.IsNotExist(err) {
		return nil, InvalidID
	} else {
		return nil, err
	}
}

// Create a new storage object using the specified directory. If the directory
// does not exist, an attempt is made to create it.
func NewStorage(directory string) (*Storage, error) {
	s := &Storage{
		directory: directory,
		bodies:    make(map[string]int),
	}
	if _, err := os.Stat(directory); err == nil {
		if err := s.loadJSON(s.indexFilename(), &s.bodies); err == nil || os.IsNotExist(err) {
			return s, nil
		} else {
			return nil, err
		}
	} else {
		if err := os.MkdirAll(directory, 0700); err == nil {
			return s, nil
		} else {
			return nil, err
		}
	}
}

// Load messages from the storage directory.
func (s *Storage) LoadMessages() ([]string, error) {
	if files, err := ioutil.ReadDir(s.directory); err == nil {
		messages := make([]string, 0)
		for _, f := range files {
			if strings.HasSuffix(f.Name(), messageExtension) {
				messages = append(messages, strings.TrimSuffix(f.Name(), messageExtension))
			}
		}
		return messages, nil
	} else {
		return nil, err
	}
}

// Create a new body and return its ID and a writer for its contents. The
// writer must be closed after the body is written.
func (s *Storage) NewBody() (io.WriteCloser, string, error) {
	id := uuid.New()
	if w, err := os.OpenFile(s.bodyFilename(id), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600); err == nil {
		return w, id, nil
	} else {
		return nil, "", err
	}
}

// Create a new message and return its ID.
func (s *Storage) NewMessage(m *Message) (string, error) {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.bodies[m.Body]; ok {
		s.bodies[m.Body]++
	} else {
		s.bodies[m.Body] = 1
	}
	if err := s.saveJSON(s.indexFilename(), s.bodies); err != nil {
		return "", err
	}
	id := uuid.New()
	if err := s.saveJSON(s.messageFilename(id), m); err != nil {
		return "", err
	}
	return id, nil
}

// Attempt to obtain a reader for the specified body.
func (s *Storage) GetBody(id string) (io.ReadCloser, error) {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.bodies[id]; ok {
		if r, err := os.Open(s.bodyFilename(id)); err == nil {
			return r, nil
		} else {
			return nil, err
		}
	} else {
		return nil, InvalidID
	}
}

// Attempt to retrieve the specified message.
func (s *Storage) GetMessage(id string) (*Message, error) {
	s.Lock()
	defer s.Unlock()
	return s.getMessage(id)
}

// Attempt to delete the specified message. If no other messages are
// referencing the message body, an attempt will be made to delete it as well.
func (s *Storage) DeleteMessage(id string) error {
	s.Lock()
	defer s.Unlock()
	if m, err := s.getMessage(id); err == nil {
		if err := os.Remove(s.messageFilename(id)); err != nil {
			return err
		}
		if n, ok := s.bodies[m.Body]; ok {
			if n == 1 {
				delete(s.bodies, m.Body)
				if err := os.Remove(s.bodyFilename(m.Body)); err != nil {
					return err
				}
			} else {
				s.bodies[m.Body]--
			}
			return s.saveJSON(s.indexFilename(), s.bodies)
		} else {
			return InvalidID
		}
	} else {
		return err
	}
}
