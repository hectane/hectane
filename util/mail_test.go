package util

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
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

func TestExtractAddresses(t *testing.T) {
	var (
		addrFrom = "from@example.com"
		addrTo   = "to@example.com"
		addrCc   = "cc@example.com"
		addrBcc  = "bcc@example.com"
		addrs    = []string{addrTo, addrCc, addrBcc}
		msg      = []byte(fmt.Sprintf(
			"From: %s\r\nTo: %s\r\nCc: %s\r\nBcc: %s\r\n\r\n",
			addrFrom,
			addrTo,
			addrCc,
			addrBcc,
		))
		r = bytes.NewBuffer(msg)
	)
	if f, a, err := ExtractAddresses(r); err == nil {
		if f != addrFrom {
			t.Fatalf("%s != %s", f, addrFrom)
		}
		sort.Strings(a)
		sort.Strings(addrs)
		if !reflect.DeepEqual(a, addrs) {
			t.Fatalf("%v != %v", a, addrs)
		}
	} else {
		t.Fatal(err)
	}
}
