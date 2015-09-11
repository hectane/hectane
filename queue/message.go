package queue

import (
	"code.google.com/p/go-uuid/uuid"

	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

// Message prepared for delivery to a specific host.
type Message struct {
	Id      string
	Host    string
	From    string
	To      []string
	Message []byte
}

// Create a message with a freshly generated ID.
func NewMessage() *Message {
	return &Message{
		Id: uuid.New(),
	}
}

// Load the message with the specified filename from disk.
func LoadMessage(filename string) (*Message, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var m Message
	if err = json.NewDecoder(f).Decode(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

// Load all messages from the specified directory. Any files that fail to
// load as messages are ignored.
func LoadMessages(directory string) ([]*Message, error) {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, err
	}
	messages := make([]*Message, 0)
	for _, f := range files {
		if m, err := LoadMessage(f.Name()); err == nil {
			messages = append(messages, m)
		}
	}
	return messages, nil
}

// Write the message to the specified directory. Permissions are carefully set
// to prevent other users from being able to read the contents.
func (m *Message) Save(directory string) error {
	f, err := os.OpenFile(path.Join(directory, m.Id), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	return json.NewEncoder(f).Encode(m)
}

// Delete the message from the specified directory.
func (m *Message) Delete(directory string) error {
	return os.Remove(path.Join(directory, m.Id))
}
