package email

import (
	"reflect"
	"testing"
)

func TestGroupAddressesByHost(t *testing.T) {
	var (
		addrList = []string{
			"A <a@hotmail.com>",
			"B <b@hotmail.com>",
			"C <c@gmail.com>",
		}
		addrMap = map[string][]string{
			"hotmail.com": []string{
				"a@hotmail.com",
				"b@hotmail.com",
			},
			"gmail.com": []string{
				"c@gmail.com",
			},
		}
	)
	if a, err := GroupAddressesByHost(addrList); err == nil {
		if !reflect.DeepEqual(addrMap, a) {
			t.Fatalf("%v != %v", addrMap, a)
		}
	} else {
		t.Fatal(err)
	}
}

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
