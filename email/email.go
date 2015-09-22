package email

import (
	"github.com/nathan-osman/go-cannon/queue"
	"github.com/nathan-osman/go-cannon/util"

	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	"strings"
	"time"
)

// Abstract representation of an email.
type Email struct {
	From        string       `json:"from"`
	To          []string     `json:"to"`
	Cc          []string     `json:"cc"`
	Bcc         []string     `json:"bcc"`
	Subject     string       `json:"subject"`
	Text        string       `json:"text"`
	Html        string       `json:"html"`
	Attachments []Attachment `json:"attachments"`
}

// Create a multipart body with the specified text and HTML and write it to the
// specified writer. A temporary buffer is used to work around a cyclical
// dependency with respect to the writer, header, and part.
func writeMultipartBody(w *multipart.Writer, text, html string) error {
	var (
		buff      = &bytes.Buffer{}
		altWriter = multipart.NewWriter(buff)
		headers   = textproto.MIMEHeader{
			"Content-Type": []string{
				fmt.Sprintf("multipart/alternative; boundary=\"%s\"", altWriter.Boundary()),
			},
		}
		textPart = &Attachment{
			ContentType: "text/plain; charset=\"utf-8\"",
			Content:     text,
		}
		htmlPart = &Attachment{
			ContentType: "text/html; charset=\"utf-8\"",
			Content:     html,
		}
	)
	part, err := w.CreatePart(headers)
	if err != nil {
		return err
	}
	if err := textPart.Write(altWriter); err != nil {
		return err
	}
	if err := htmlPart.Write(altWriter); err != nil {
		return err
	}
	if err := altWriter.Close(); err != nil {
		return err
	}
	_, err = io.Copy(part, buff)
	return err
}

// Convert the email into an array of messages grouped by host suitable for
// delivery to the mail queue.
func (e *Email) Messages(q *queue.Queue) ([]*queue.Message, error) {
	if w, id, err := q.NewBody(); err == nil {
		var (
			m       = multipart.NewWriter(w)
			headers = EmailHeaders{
				"Message-Id":   fmt.Sprintf("<%s@go-cannon>", id),
				"From":         e.From,
				"To":           strings.Join(e.To, ", "),
				"Subject":      e.Subject,
				"Date":         time.Now().Format("Mon, 02 Jan 2006 15:04:05 -0700"),
				"MIME-Version": "1.0",
				"Content-Type": fmt.Sprintf("multipart/mixed; boundary=%s", m.Boundary()),
			}
			addresses = append(append(e.To, e.Cc...), e.Bcc...)
		)
		if len(e.Cc) > 0 {
			headers["Cc"] = strings.Join(e.Cc, ",")
		}
		if err := headers.Write(w); err != nil {
			return nil, err
		}
		if err := writeMultipartBody(m, e.Text, e.Html); err != nil {
			return nil, err
		}
		for _, a := range e.Attachments {
			if err := a.Write(m); err != nil {
				return nil, err
			}
		}
		if err := m.Close(); err != nil {
			return nil, err
		}
		if err := w.Close(); err != nil {
			return nil, err
		}
		if addrMap, err := util.GroupAddressesByHost(addresses); err == nil {
			messages := make([]*queue.Message, 0, 1)
			for h, to := range addrMap {
				messages = append(messages, &queue.Message{
					Host: h,
					From: e.From,
					To:   to,
					Body: id,
				})
			}
			return messages, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}
