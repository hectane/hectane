package queue

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
)

// Metadata for the message body.
type bodyMetadata struct {
	RefCount int
}

// Message body. It is preserved on disk as long as one or more messages
// reference it. Once no more messages reference the body, it is removed.
type Body struct {
	directory string
	id        string
	m         bodyMetadata
}

// Determine the name of the file where the message body is stored.
func (b *Body) messageBodyFilename() string {
	return path.Join(b.directory, b.id)
}

// Determine the name of the file where metadata is stored.
func (b *Body) metadataFilename() string {
	return fmt.Sprintf("%s.json", b.messageBodyFilename())
}

// Update the metadata on disk.
func (b *Body) updateMetadata() error {
	f, err := os.OpenFile(b.metadataFilename(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	return json.NewEncoder(f).Encode(&b.m)
}

// Create a new message body from the specified reader.
func NewBody(directory, id string) (*Body, io.WriteCloser, error) {
	b := &Body{
		id:        id,
		directory: directory,
	}
	if f, err := os.OpenFile(b.messageBodyFilename(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600); err != nil {
		return nil, nil, err
	} else {
		return b, f, nil
	}
}

// Attempt to load the specified message body.
func LoadBody(directory, id string) (*Body, error) {
	b := &Body{
		directory: directory,
		id:        id,
	}
	if f, err := os.Open(b.metadataFilename()); err == nil {
		if err = json.NewDecoder(f).Decode(&b.m); err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}
	return b, nil
}

// Increase the reference count for the body.
func (b *Body) Add() error {
	b.m.RefCount++
	return b.updateMetadata()
}

// Decrease the reference count for the body.
func (b *Body) Release() error {
	b.m.RefCount--
	if b.m.RefCount == 0 {
		if err := os.Remove(b.metadataFilename()); err != nil {
			return err
		}
		return os.Remove(b.messageBodyFilename())
	} else {
		return b.updateMetadata()
	}
}

// Obtain an io.Reader for the message body.
func (b *Body) Reader() (io.Reader, error) {
	return os.Open(b.messageBodyFilename())
}
