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

const (
	StatusDelivered = "delivered"
	StatusError     = "error"
)

// Write the specified JSON object to the client. No error checking is done
// since little could be done if the methods were to fail anyway.
func respondWithJSON(w http.ResponseWriter, o interface{}) {
	data, _ := json.Marshal(o)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// Indicate status to the client.
func respondWithStatus(w http.ResponseWriter, status string) {
	respondWithJSON(w, map[string]string{
		"status": status,
	})
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

// Return version information
func Version(c web.C, w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, map[string]string{
		"version": "0.1.0",
	})
}

// Attempt to send an email to the specified recipient.
func Send(c web.C, w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithStatus(w, StatusError)
		return
	}
	email, err := NewEmailFromJson(body)
	if err != nil {
		respondWithStatus(w, StatusError)
		return
	}
	host, err := email.Host()
	if err != nil {
		respondWithStatus(w, StatusError)
		return
	}
	client, err := tryMailServers(host)
	if err != nil {
		respondWithStatus(w, StatusError)
		return
	}
	defer client.Close()
	client.Mail(email.From)
	client.Rcpt(email.To)
	writer, err := client.Data()
	if err != nil {
		respondWithStatus(w, StatusError)
		return
	}
	defer writer.Close()
	email.Write(writer)
	respondWithStatus(w, StatusDelivered)
}
