package main

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/mail"
	"net/textproto"
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

// Add the specified content to a multipart/alternative writer.
func addPart(w *multipart.Writer, contentType string, content []byte) error {
	header := make(textproto.MIMEHeader)
	header.Add("Content-Type", contentType)
	writer, err := w.CreatePart(header)
	if err != nil {
		return err
	}
	_, err = writer.Write(content)
	return err
}

// Write the email to a Writer instance.
func (e Email) Write(w io.Writer) error {
	writer := multipart.NewWriter(w)

	// Write the headers, including the content-type and boundary
	w.Write([]byte(fmt.Sprintf("From: %s\r\n", e.From)))
	w.Write([]byte(fmt.Sprintf("To: %s\r\n", e.To)))
	w.Write([]byte(fmt.Sprintf("Subject: %s\r\n", e.Subject)))
	w.Write([]byte(fmt.Sprintf("Content-Type: multipart/alternative; boundary=\"%s\"", writer.Boundary())))
	w.Write([]byte("\r\n"))

	// Write the parts of the message
	if err := addPart(writer, "text/plain", []byte(e.Text)); err != nil {
		return err
	}
	if err := addPart(writer, "text/html", []byte(e.Html)); err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}

	return nil
}
