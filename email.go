package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/mail"
	"strings"
)

// Single email message for delivery. This struct is used both for
// unmarshalling requests and storing emails queued for delivery on disk.
type Email struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Text    string `json:"text"`
	Html    string `json:"html"`
}

// Attempt to parse the provided JSON data and populate an Email struct.
func NewEmailFromJson(data []byte) (*Email, error) {
	var email Email
	if err := json.Unmarshal(data, &email); err != nil {
		return nil, err
	} else {
		return &email, nil
	}
}

// Attempt to extract the host from the To email address.
func (e Email) Host() (string, error) {
	if addr, err := mail.ParseAddress(e.To); err != nil {
		return "", err
	} else {
		// If ParseAddress succeeded, Address is guaranteed to contain "@"
		return strings.Split(addr.Address, "@")[1], nil
	}
}

// Write the email to a Writer instance.
func (e Email) Write(w io.Writer) error {
	w.Write([]byte(fmt.Sprintf("From: %s\r\n", e.From)))
	w.Write([]byte(fmt.Sprintf("To: %s\r\n", e.To)))
	w.Write([]byte(fmt.Sprintf("Subject: %s\r\n", e.Subject)))
	w.Write([]byte("Content-Type: text/plain\r\n\r\n"))
	w.Write([]byte(e.Text))
	return nil
}
