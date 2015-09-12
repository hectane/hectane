package email

import (
	"fmt"
	"io"
)

// Map of email headers.
type EmailHeaders map[string]string

// Write the headers to the specified io.Writer.
func (e EmailHeaders) Write(w io.Writer) error {
	for k, v := range e {
		if _, err := w.Write([]byte(fmt.Sprintf("%s: %s\r\n", k, v))); err != nil {
			return err
		}
	}
	return nil
}
