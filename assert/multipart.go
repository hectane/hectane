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
func Multipart(r io.Reader, boundary string, m *MultipartDesc) error {
	reader := multipart.NewReader(r, boundary)
	for {
		if p, err := reader.NextPart(); err == nil {
			if contentType, params, err := mime.ParseMediaType(p.Header.Get("Content-Type")); err == nil {
				if d, ok := m.Parts[contentType]; ok {
					if len(d.Parts) == 0 {
						if err := Read(p, d.Content); err != nil {
							return err
						}
					} else {
						if boundary, ok := params["boundary"]; ok {
							return Multipart(p, boundary, d)
						} else {
							return errors.New("\"boundary\" parameter missing")
						}
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
}
