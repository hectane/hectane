package email

import (
	"testing"
)

func TestToHTML(t *testing.T) {
	data := []struct {
		i, o string
	}{
		{"&<", "&amp;&lt;"},
		{"a\n\nb", "a<br><br>b"},
		{"a http://example.org b", "a <a href=\"http://example.org\">http://example.org</a> b"},
	}
	for _, d := range data {
		if o := toHTML(d.i); o != d.o {
			t.Fatalf("%s != %s", o, d.o)
		}
	}
}
