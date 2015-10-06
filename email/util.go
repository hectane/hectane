package email

import (
	"fmt"
	"html"
	"regexp"
	"strings"
)

// The link RegExp is adapted from http://stackoverflow.com/a/3809435/193619.
var (
	replaceCRLF  = regexp.MustCompile(`\r?\n`)
	replaceLinks = regexp.MustCompile(`(?:https?:\/\/)?[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b(?:[-a-zA-Z0-9@:%_\+.~#?&//=]*)`)
)

// Convert the specified text to its HTML equivalent, preserving formatting
// where possible and converting URLs to <a> elements.
func toHTML(data string) string {
	data = html.EscapeString(data)
	data = replaceCRLF.ReplaceAllString(data, "<br>")
	data = replaceLinks.ReplaceAllStringFunc(data, func(m string) string {
		if !strings.HasPrefix(m, "http://") && !strings.HasPrefix(m, "https") {
			m = fmt.Sprintf("http://%s", m)
		}
		return fmt.Sprintf("<a href=\"%s\">%s</a>", m, m)
	})
	return data
}
