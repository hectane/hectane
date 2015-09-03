package email

import (
	"bytes"
	"mime/multipart"
	"net/mail"
	"net/textproto"
	"strings"
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
