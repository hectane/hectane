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
			t.Fatal("address map doesn't match")
		}
	} else {
		t.Fatal(err)
	}
}
