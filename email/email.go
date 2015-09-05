package email

import (
	"code.google.com/p/go-uuid/uuid"

	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/mail"
	"net/textproto"
	"os"
	"path"
	"strings"
	"time"
)

// Individual message for sending to a host.
type Email struct {
	Id      string
	Host    string
	From    string
	To      []string
	Message []byte
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

// Create a multipart/alternative body with the specified data. Both the body
// and the separator are returned. Because data is written to a bytes.Buffer,
// it is probably safe to ignore any errors.
func createMultipartBody(text, html string) ([]byte, string) {
	var (
		buff   = &bytes.Buffer{}
		writer = multipart.NewWriter(buff)
	)
	addPart(writer, "text/plain", []byte(text))
	addPart(writer, "text/html", []byte(html))
	writer.Close()
	return buff.Bytes(), writer.Boundary()
}

// Attempt to extract the host from the provided email address.
func hostFromAddress(data string) (string, error) {
	if addr, err := mail.ParseAddress(data); err != nil {
		return "", err
	} else {
		return strings.Split(addr.Address, "@")[1], nil
	}
}

// Group a list of email addresses by their host.
func groupAddressesByHost(addrs []string) (map[string][]string, error) {
	m := make(map[string][]string)
	for _, addr := range addrs {
		if host, err := hostFromAddress(addr); err != nil {
			return nil, err
		} else {
			if m[host] == nil {
				m[host] = make([]string, 0, 1)
			}
			m[host] = append(m[host], addr)
		}
	}
	return m, nil
}

// Create an array of emails from the specified information. An email will be
// generated for each individual host and for each BCC recipient.
func NewEmails(from string, to, cc, bcc []string, subject string, text, html string) ([]*Email, error) {

	// Retrieve the current hostname
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	// Create the array of emails that will be returned, the buffer that will
	// be used for generating the message, and each of the individual parts
	// of the message
	var (
		emails         = make([]*Email, 0, 1)
		buff           = &bytes.Buffer{}
		id             = fmt.Sprintf("<%s@%s>", uuid.New(), hostname)
		body, boundary = createMultipartBody(text, html)
		hdrs           = map[string]string{
			"Message-ID":   id,
			"From":         from,
			"To":           strings.Join(to, ","),
			"Subject":      subject,
			"Date":         time.Now().Format(time.RubyDate),
			"MIME-Version": "1.0",
			"Content-Type": fmt.Sprintf("multipart/alternative; boundary=\"%s\"", boundary),
		}
	)

	// If any Cc addresses were provided, add them to the headers
	if len(cc) > 0 {
		hdrs["Cc"] = strings.Join(cc, ",")
	}

	// Write all of the headers followed by the body
	for h := range hdrs {
		buff.Write([]byte(fmt.Sprintf("%s: %s\r\n", h, hdrs[h])))
	}

	// Write the extra CRLF and body
	buff.Write([]byte("\r\n"))
	buff.Write(body)

	// Combine the "To", "Cc", and "Bcc" addresses and group them by hostname
	// in order to make delivery much more efficient
	addrMap, err := groupAddressesByHost(append(append(to, cc...), bcc...))
	for host := range addrMap {
		emails = append(emails, &Email{
			Id:      id,
			Host:    host,
			From:    from,
			To:      addrMap[host],
			Message: buff.Bytes(),
		})
	}

	return emails, nil
}

// Load an email from disk.
func LoadEmail(filename string) (*Email, error) {

	// Attempt to open the file
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	// Decode the JSON from the file
	var e Email
	if err = json.NewDecoder(f).Decode(&e); err != nil {
		return nil, err
	}

	return &e, nil
}

// Save the email to the specified directory.
func (e *Email) Save(directory string) error {

	// Attempt to create the , making sure nobody else can read it
	f, err := os.OpenFile(path.Join(directory, e.Id), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	// Encode the data as JSON
	if err = json.NewEncoder(f).Encode(e); err != nil {
		return err
	}

	return nil
}

// Delete the email from the specified directory.
func (e *Email) Delete(directory string) error {
	return os.Remove(path.Join(directory, e.Id))
}
