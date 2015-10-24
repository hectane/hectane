package api

import (
	"fmt"
	"net/http"
	"testing"
)

func getStatusCode(url, username, password string) (int, error) {
	if r, err := http.NewRequest("GET", url, nil); err == nil {
		if username != "" && password != "" {
			r.SetBasicAuth(username, password)
		}
		if resp, err := http.DefaultClient.Do(r); err == nil {
			return resp.StatusCode, nil
		} else {
			return 0, err
		}
	} else {
		return 0, err
	}
}

func TestBasicAuth(t *testing.T) {
	var (
		addr        = "127.0.0.1:9000"
		url         = fmt.Sprintf("http://%s/v1/version", addr)
		username    = "test"
		password    = "test"
		badPassword = "test2"
		testCases   = []struct {
			username   string
			password   string
			statusCode int
		}{
			{"", "", http.StatusUnauthorized},
			{username, badPassword, http.StatusForbidden},
			{username, password, http.StatusOK},
		}
		config = &Config{
			Addr:     addr,
			Username: username,
			Password: password,
		}
		a = New(config, nil)
	)
	if err := a.Start(); err != nil {
		t.Fatal(err)
	}
	defer a.Stop()
	for _, c := range testCases {
		if s, err := getStatusCode(url, c.username, c.password); err == nil {
			if s != c.statusCode {
				t.Fatalf("%d != %d", s, c.statusCode)
			}
		} else {
			t.Fatal(err)
		}
	}
}
