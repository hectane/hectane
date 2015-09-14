package util

import (
	"reflect"
	"testing"
)

func TestGroupAddressesByHost(t *testing.T) {
	var (
		addrList = []string{
			"a@hotmail.com",
			"b@hotmail.com",
			"c@gmail.com",
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

func TestHostFromAddress(t *testing.T) {
	var (
		addr = "a@hotmail.com"
		host = "hotmail.com"
	)
	if v, err := HostFromAddress(addr); err == nil {
		if v != host {
			t.Fatalf("%s != %s", v, host)
		}
	} else {
		t.Fatal(err)
	}
}
