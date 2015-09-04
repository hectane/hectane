package email

import (
	"code.google.com/p/go-uuid/uuid"

	"bytes"
	"fmt"
	"mime/multipart"
	"net/mail"
	"net/textproto"
	"os"
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

// Generate a Message-ID.
func generateMessageId(hostname string) string {
	return fmt.Sprintf("<%s@%s>", uuid.New(), hostname)
}

// Create a message from the specified headers and data.
func createMessage(hdrs map[string]string, body []byte) []byte {
	buff := &bytes.Buffer{}
	for hdr := range hdrs {
		buff.Write([]byte(fmt.Sprintf("%s: %s\r\n", hdr, hdrs[hdr])))
	}
	buff.Write([]byte("\r\n"))
	buff.Write(body)
	return buff.Bytes()
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

	// Create the list of emails that will be returned, build a map of headers,
	// and generate the message body - it will be identical for all Email
	// instances returned
	var (
		emails         = make([]*Email, 0, 1)
		body, boundary = createMultipartBody(text, html)
		id             = generateMessageId(hostname)
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

	// Add Cc addresses if any were provided
	if len(cc) > 0 {
		hdrs["Cc"] = strings.Join(cc, ",")
	}

	// The first email will be sent to all "To" and "Cc" addresses and does not
	// contain any of the Bcc addresses
	addrMap, err := groupAddressesByHost(append(to, cc...))
	if err != nil {
		return nil, err
	}

	// Generate emails for each of the hosts
	for host := range addrMap {
		emails = append(emails, &Email{
			Id:      id,
			Host:    host,
			From:    from,
			To:      addrMap[host],
			Message: createMessage(hdrs, body),
		})
	}

	// Generate emails for each of the Bcc addresses
	for _, addr := range bcc {

		// Reset a couple of headers
		id = generateMessageId(hostname)
		hdrs["Message-ID"] = id
		hdrs["Bcc"] = addr

		// Fetch the host
		host, err := hostFromAddress(addr)
		if err != nil {
			return nil, err
		}

		// Add the email to the list
		emails = append(emails, &Email{
			Id:      id,
			Host:    host,
			From:    from,
			To:      []string{addr},
			Message: createMessage(hdrs, body),
		})
	}

	return emails, nil
}
