package queue

import (
	"code.google.com/p/go-uuid/uuid"

	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/smtp"
	"os"
	"path"
	"strings"
)

// Message metadata.
type messageMetadata struct {
	Host  string
	From  string
	To    []string
	Tries int
	Body  string
}

// Message prepared for delivery to a specific host.
type Message struct {
	directory string
	filename  string
	m         messageMetadata
}

// Determine the name of the file where metadata is stored.
func (m *Message) metadataFilename() string {
	return path.Join(m.directory, m.filename)
}

// Update the metadata on disk.
func (m *Message) updateMetadata() error {

	f, err := os.OpenFile(m.metadataFilename(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	return json.NewEncoder(f).Encode(&m.m)
}

// Create a new message with the specified information.
func NewMessage(directory, host, from string, to []string, body string) (*Message, error) {
	m := &Message{
		directory: directory,
		filename:  fmt.Sprintf("%s.msg", uuid.New()),
		m: messageMetadata{
			Host: host,
			From: from,
			To:   to,
			Body: body,
		},
	}
	if b, err := LoadBody(directory, body); err == nil {
		if err := b.Add(); err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}
	if err := m.updateMetadata(); err != nil {
		return nil, err
	}
	return m, nil
}

// Load an existing message from the specified location.
func LoadMessage(directory, filename string) (*Message, error) {
	m := &Message{
		directory: directory,
		filename:  filename,
	}
	if f, err := os.Open(m.metadataFilename()); err == nil {
		if err = json.NewDecoder(f).Decode(&m.m); err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}
	return m, nil
}

// Load all messages from the specified directory. If the metadata of a message
// could not be loaded, it will be skipped.
func LoadMessages(directory string) ([]*Message, error) {
	if files, err := ioutil.ReadDir(directory); err == nil {
		messages := make([]*Message, 0)
		for _, f := range files {
			if strings.HasSuffix(f.Name(), ".msg") {
				if m, err := LoadMessage(directory, f.Name()); err == nil {
					messages = append(messages, m)
				}
			}
		}
		return messages, nil
	} else {
		return nil, err
	}
}

// Attempt to send the specified message using the specified client.
func (m *Message) Send(c *smtp.Client) error {
	if err := c.Mail(m.m.From); err != nil {
		return err
	}
	for _, t := range m.m.To {
		if err := c.Rcpt(t); err != nil {
			return err
		}
	}
	if w, err := c.Data(); err == nil {
		if b, err := LoadBody(m.directory, m.m.Body); err == nil {
			if r, err := b.reader(); err == nil {
				if _, err := io.Copy(w, r); err != nil {
					return err
				}
			} else {
				return err
			}
		} else {
			return err
		}
		if err := w.Close(); err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}

// Delete the message on disk. Note that this will decrease the reference count
// of the message body, causing it to be deleted if the reference count hits 0.
func (m *Message) Delete() error {
	if err := os.Remove(m.metadataFilename()); err != nil {
		return err
	}
	if b, err := LoadBody(m.directory, m.m.Body); err == nil {
		b.Release()
		return nil
	} else {
		return err
	}
}
