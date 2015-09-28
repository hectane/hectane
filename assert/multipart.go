package assert

import (
	"errors"
	"io"
	"mime"
	"mime/multipart"
)

// Description of a multipart MIME message.
type MultipartDesc struct {
	Content []byte
	Parts   map[string]*MultipartDesc
}

// Ensure that a multipart message conforms to the specified description.
func Multipart(r io.Reader, contentType string, m *MultipartDesc) error {
	if _, params, err := mime.ParseMediaType(contentType); err != nil {
		if len(m.Parts) == 0 {
			return Read(r, m.Content)
		} else {
			if boundary, ok := params["boundary"]; ok {
				reader := multipart.NewReader(r, boundary)
				for {
					if p, err := reader.NextPart(); err == nil {
						if c, _, err := mime.ParseMediaType(p.Header.Get("Content-Type")); err == nil {
							if d, ok := m.Parts[c]; ok {
								if err := Multipart(p, p.Header.Get("Content-Type"), d); err != nil {
									return err
								}
							} else {
								return errors.New("unexpected content type")
							}
						} else {
							return err
						}
					} else if err == io.EOF {
						break
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
