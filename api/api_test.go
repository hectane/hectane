package api

import (
	"github.com/hectane/go-attest"

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
	if err := attest.HttpStatusCode(req, http.StatusUnauthorized); err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth(username, badPassword)
	if err := attest.HttpStatusCode(req, http.StatusUnauthorized); err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth(username, password)
	if err := attest.HttpStatusCode(req, http.StatusOK); err != nil {
		t.Fatal(err)
	}
}

func TestShutdown(t *testing.T) {
	a, req, err := createServer("", "")
	if err != nil {
		t.Fatal(err)
	}
	a.Stop()
	req.URL.Path = "/v1/version"
	if err := attest.HttpStatusCode(req, 0); err == nil {
		t.Fatal("error expected")
	}
}
