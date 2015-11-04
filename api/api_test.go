package api

import (
	"fmt"
	"net/http"
	"testing"
)

func getStatusCode(url, username, password string) (int, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	if username != "" && password != "" {
		req.SetBasicAuth(username, password)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	return resp.StatusCode, nil
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
			{username, badPassword, http.StatusUnauthorized},
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
		s, err := getStatusCode(url, c.username, c.password)
		if err != nil {
			t.Fatal(err)
		}
		if s != c.statusCode {
			t.Fatalf("%d != %d", s, c.statusCode)
		}
	}
}
