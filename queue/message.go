package queue

import (
	"code.google.com/p/go-uuid/uuid"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// Message prepared for delivery to a specific host. When persisted to disk,
// the message metadata is stored separately from the message body.
type Message struct {
	Id        string `json:"-"`
	Host      string
	From      string
	To        []string
	Tries     int
	Message   []byte `json:"-"`
	Directory string `json:"-"`
}

// Create a new message with a freshly generated ID.
func NewMessage(host, from string, to []string, message []byte, directory string) (*Message, error) {
	m := &Message{
		Id:        uuid.New(),
		Host:      host,
		From:      from,
		To:        to,
		Message:   message,
		Directory: directory,
	}
	if err := ioutil.WriteFile(m.bodyFilename(), message, 0600); err != nil {
		return nil, err
	}
	if err := m.Update(); err != nil {
		return nil, err
	}
	return m, nil
}

// Load the specified message from the directory with the specified ID.
func LoadMessage(directory, id string) (*Message, error) {
	m := &Message{
		Id:        id,
		Directory: directory,
	}
	if f, err := os.Open(m.metadataFilename()); err == nil {
		if err = json.NewDecoder(f).Decode(m); err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}
	if message, err := ioutil.ReadFile(m.bodyFilename()); err != nil {
		return nil, err
	} else {
		m.Message = message
	}
	return m, nil
}

// Load all messages from the specified directory. If either the metadata of a
// message or the message body could not be loaded, it will be skipped.
func LoadMessages(directory string) ([]*Message, error) {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, err
	}
	messages := make([]*Message, 0)
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".json") {
			if m, err := LoadMessage(directory, f.Name()); err == nil {
				messages = append(messages, m)
			}
		}
	}
	return messages, nil
}

// Determine the name of the file used for storing message metadata.
func (m *Message) metadataFilename() string {
	return fmt.Sprintf("%s.json", m.bodyFilename())
}

// Determine the name of the file containing the message body.
func (m *Message) bodyFilename() string {
	return path.Join(m.Directory, m.Id)
}

// Attempt to update the message metadata on disk.
func (m *Message) Update() error {
	f, err := os.OpenFile(m.metadataFilename(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	return json.NewEncoder(f).Encode(m)
}

// Delete the message on disk - both the metadata and the message body.
func (m *Message) Delete() error {
	if err := os.Remove(m.metadataFilename()); err != nil {
		return err
	}
	return os.Remove(m.bodyFilename())
}
