package assert

import (
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
)

// Description of a multipart MIME message.
type MultipartDesc struct {
	ContentType string
	Content     []byte
	Parts       []*MultipartDesc
}

// Ensure that a multipart message conforms to the specified description.
func Multipart(r io.Reader, contentType string, d *MultipartDesc) error {
	if c, params, err := mime.ParseMediaType(contentType); err == nil {
		if c != d.ContentType {
			return errors.New(fmt.Sprintf("%s != %s", c, d.ContentType))
		}
		if len(d.Parts) == 0 {
			return Read(r, d.Content)
		} else {
			if boundary, ok := params["boundary"]; ok {
				reader := multipart.NewReader(r, boundary)
				for _, part := range d.Parts {
					if p, err := reader.NextPart(); err == nil {
						if err := Multipart(p, p.Header.Get("Content-Type"), part); err != nil {
							return err
						}
					} else {
						return err
					}
				}
				return nil
			} else {
				return errors.New("\"boundary\" parameter missing")
			}
		}
	} else {
		return err
	}
}
