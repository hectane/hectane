package email

import (
	"bytes"
	"io"
	"mime"
	"mime/multipart"
	"net/mail"
	"testing"
)

func TestEmailMessages(t *testing.T) {
	var (
		from    = "a@hotmail.com"
		to      = "b@hotmail.com"
		cc      = "c@hotmail.com"
		bcc     = "d@yahoo.com"
		subject = "Test"
		text    = "test"
		html    = "<em>test</em>"
		e       = &Email{
			From:    from,
			To:      []string{to},
			Cc:      []string{cc},
			Bcc:     []string{bcc},
			Subject: subject,
			Text:    text,
			Html:    html,
		}
	)
	if messages, err := e.Messages(); err == nil {
		hosts := make([]string, 0, 2)
		for _, m := range messages {
			buff := bytes.NewBuffer(m.Message)
			if msg, err := mail.ReadMessage(buff); err == nil {
				hosts = append(hosts, m.Host)
				if _, params, err := mime.ParseMediaType(msg.Header.Get("Content-Type")); err == nil {
					if boundary, ok := params["boundary"]; ok {
						r := multipart.NewReader(msg.Body, boundary)
						for {
							_, err := r.NextPart()
							if err == io.EOF {
								break
							} else if err != nil {
								t.Fatal(err)
							} else {
								// TODO
							}
						}
					} else {
						t.Fatal("content type missing boundary")
					}
				} else {
					t.Fatal(err)
				}
			} else {
				t.Fatal(err)
			}
		}
	} else {
		t.Fatal(err)
	}
}
