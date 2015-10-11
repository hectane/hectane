package api

import (
	"fmt"
	"net/http"
	"testing"
)

func TestBasicAuth(t *testing.T) {
	var (
		username = "test"
		password = "test"
		config   = &Config{
			Addr:     "127.0.0.1:0",
			Username: username,
			Password: password,
		}
		a = New(config, nil)
	)
	if resp, err := http.Get(fmt.Sprintf("http://%s/v1/version", addr)); err == nil {
		t.Log(resp)
	} else {
		t.Fatal(err)
	}
}
