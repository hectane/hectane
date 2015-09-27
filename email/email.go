package email

import (
	"github.com/nathan-osman/go-cannon/queue"
	"github.com/nathan-osman/go-cannon/util"

	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/mail"
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

// Write the headers for the email to the specified writer.
func (e *Email) writeHeaders(w io.Writer, id, boundary string) error {
	headers := EmailHeaders{
		"Message-Id":   fmt.Sprintf("<%s@go-cannon>", id),
		"From":         e.From,
		"To":           strings.Join(e.To, ", "),
		"Subject":      e.Subject,
		"Date":         time.Now().Format("Mon, 02 Jan 2006 15:04:05 -0700"),
		"MIME-Version": "1.0",
		"Content-Type": fmt.Sprintf("multipart/mixed; boundary=%s", boundary),
	}
	if len(e.Cc) > 0 {
		headers["Cc"] = strings.Join(e.Cc, ", ")
	}
	return headers.Write(w)
}

// Write the body of the email to the specified writer.
func (e *Email) writeBody(w *multipart.Writer) error {
	var (
		buff      = &bytes.Buffer{}
		altWriter = multipart.NewWriter(buff)
		header    = textproto.MIMEHeader{
			"Content-Type": []string{
				fmt.Sprintf("multipart/alternative; boundary=%s", altWriter.Boundary()),
			},
		}
	)
	if p, err := w.CreatePart(header); err == nil {
		if err := (Attachment{
			ContentType: "text/plain; charset=utf-8",
			Content:     e.Text,
		}.Write(w)); err != nil {
			return err
		}
		if err := (Attachment{
			ContentType: "text/html; charset=utf-8",
			Content:     e.Html,
		}.Write(w)); err != nil {
			return err
		}
		if err := altWriter.Close(); err != nil {
			return err
		}
		if _, err := io.Copy(p, buff); err != nil {
			return err
		}
		return nil
	} else {
		return err
	}
}

// Create an array of messages with the specified body.
func (e *Email) newMessages(s *queue.Storage, from, body string) ([]*queue.Message, error) {
	addresses := append(append(e.To, e.Cc...), e.Bcc...)
	if m, err := util.GroupAddressesByHost(addresses); err == nil {
		messages := make([]*queue.Message, 0, 1)
		for h, to := range m {
			msg := &queue.Message{
				Host: h,
				From: from,
				To:   to,
			}
			if err := s.SaveMessage(msg, body); err != nil {
				return nil, err
			}
			messages = append(messages, msg)
		}
		return messages, nil
	} else {
		return nil, err
	}
}

// Convert the email into an array of messages grouped by host suitable for
// delivery to the mail queue.
func (e *Email) Messages(s *queue.Storage) ([]*queue.Message, error) {
	if from, err := mail.ParseAddress(e.From); err == nil {
		if w, body, err := s.NewBody(); err == nil {
			mpWriter := multipart.NewWriter(w)
			if err := e.writeHeaders(w, body, mpWriter.Boundary()); err != nil {
				return nil, err
			}
			if err := e.writeBody(mpWriter); err != nil {
				return nil, err
			}
			for _, a := range e.Attachments {
				if err := a.Write(mpWriter); err != nil {
					return nil, err
				}
			}
			if err := mpWriter.Close(); err != nil {
				return nil, err
			}
			if err := w.Close(); err != nil {
				return nil, err
			}
			return e.newMessages(s, from.Address, body)
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}
