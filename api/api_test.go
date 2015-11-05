package api

import (
	"net/http"
	"net/url"
	"testing"
)

func createServer(username, password string) (*API, *http.Request, error) {
	a := New(&Config{
		Addr:     "127.0.0.1:0",
		Username: username,
		Password: password,
	}, nil)
	if err := a.Start(); err != nil {
		return nil, nil, err
	}
	return a, &http.Request{
		URL: &url.URL{
			Scheme: "http",
			Host:   a.listener.Addr().String(),
		},
	}, nil
}

func TestBasicAuth(t *testing.T) {
	var (
		username    = "test"
		password    = "test"
		badPassword = "test2"
	)
	a, req, err := createServer(username, password)
	if err != nil {
		t.Fatal(err)
	}
	defer a.Stop()
	req.URL.Path = "/v1/version"
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("%d != %d", resp.StatusCode, http.StatusUnauthorized)
	}
	req.SetBasicAuth(username, badPassword)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("%d != %d", resp.StatusCode, http.StatusUnauthorized)
	}
	req.SetBasicAuth(username, password)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("%d != %d", resp.StatusCode, http.StatusOK)
	}
}
