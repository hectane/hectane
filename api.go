package main

import (
	"github.com/zenazn/goji/web"

	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/smtp"
)

// Write the specified error message as JSON to the client. No error checking
// is done since little could be done if these methods fail anyway.
func respondWithError(w http.ResponseWriter, message string) {
	data, _ := json.Marshal(map[string]string{
		"error": message,
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// Find the first mail server that responds, trying MX records if present or
// (if none found) using the domain's A record.
func tryMailServers(host string) (*smtp.Client, error) {
	mxs, err := net.LookupMX(host)
	if err != nil {
		return smtp.Dial(fmt.Sprintf("%s:25", host))
	}
	for _, mx := range mxs {
		client, err := smtp.Dial(fmt.Sprintf("%s:25", mx.Host))
		if err == nil {
			return client, nil
		}
	}
	return nil, errors.New("unable to connect to any servers listed in MX records")
}

// Attempt to send an email to the specified recipient.
func Send(c web.C, w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, err.Error())
		return
	}
	email, err := NewEmailFromJson(body)
	if err != nil {
		respondWithError(w, err.Error())
		return
	}
	host, err := email.Host()
	if err != nil {
		respondWithError(w, err.Error())
		return
	}
	client, err := tryMailServers(host)
	if err != nil {
		respondWithError(w, err.Error())
		return
	}
	defer client.Close()
	client.Mail(email.From)
	client.Rcpt(email.To)
	writer, err := client.Data()
	if err != nil {
		respondWithError(w, err.Error())
		return
	}
	defer writer.Close()
	email.Write(writer)
}
