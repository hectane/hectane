package main

import (
	"github.com/goji/httpauth"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"

	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/mail"
	"net/smtp"
	"strings"
)

// Find the first mail server that responds, trying MX records if present or
// (if none found) using the domain's A record.
func tryMailServers(name string) (*smtp.Client, error) {
	mxs, err := net.LookupMX(name)
	if err != nil {
		return smtp.Dial(fmt.Sprintf("%s:25", name))
	}
	for _, mx := range mxs {
		client, err := smtp.Dial(fmt.Sprintf("%s:25", mx.Host))
		if err == nil {
			return client, nil
		}
	}
	return nil, errors.New("unable to connect to any servers listed in MX records")
}

type sendParams struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Text    string `json:"text"`
	Html    string `json:"html"`
}

// Send an email with the specified parameters to the specified recipient.
func send(c web.C, w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	var params sendParams
	if err = json.Unmarshal(body, &params); err != nil {
		return
	}
	addr, err := mail.ParseAddress(params.To)
	if err != nil {
		return
	}
	parts := strings.Split(addr.Address, "@")
	if len(parts) != 2 {
		return
	}
	client, err := tryMailServers(parts[1])
	if err != nil {
		return
	}
	client.Mail(params.From)
	client.Rcpt(params.To)
	writer, _ := client.Data()
	writer.Write([]byte(fmt.Sprintf("From: %s\r\n", params.From)))
	writer.Write([]byte(fmt.Sprintf("To: %s\r\n", params.To)))
	writer.Write([]byte(fmt.Sprintf("Subject: %s\r\n", params.Subject)))
	writer.Write([]byte("Content-Type: text/plain\r\n\r\n"))
	writer.Write([]byte(params.Text))
	writer.Close()
	client.Close()
}

func main() {
	var (
		username string
		password string
	)

	flag.StringVar(&username, "username", "", "username for HTTP basic auth")
	flag.StringVar(&password, "password", "", "password for HTTP basic auth")
	flag.Parse()

	// If username and password are provided, enable HTTP basic auth
	if username != "" && password != "" {
		goji.Use(httpauth.SimpleBasicAuth(username, password))
	}

	goji.Post("/v1/send", send)
	goji.Serve()
}
