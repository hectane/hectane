package queue

import (
	"errors"
	"fmt"
	"net"
	"net/smtp"
	"time"
)

// Persistent connection to an SMTP host.
type Host struct {
	name string
	stop chan bool
}

// Attempt to find the mail servers for the specified host.
func findMailServers(host string) []string {

	// First check for MX records - if one or more were found, convert the
	// records into a list of strings (already sorted by priority) - if none
	// were found, then simply return the host that was originally provided
	if mx, err := net.LookupMX(host); err == nil {
		servers := make([]string, len(mx))
		for i, record := range mx {
			servers[i] = record.Host
		}
		return servers
	} else {
		return []string{host}
	}
}

// Attempt to connect to the specified host.
func connectToMailServer(host string, stop chan bool) (*smtp.Client, error) {

	// Obtain the list of mail servers to try
	servers := findMailServers(host)

	// RFC 5321 (4.5.4) describes the process for retrying connections to a
	// mail server after failure. The recommended strategy is to retry twice
	// with 30 minute intervals and continue at 120 minute intervals until four
	// days have elapsed.
	for i := 0; i < 50; i++ {

		// Attempt to connect to each of the mail servers from the list in the
		// order that was provided - return immediately if a connection is made
		for _, server := range servers {
			if client, err := smtp.Dial(fmt.Sprintf("%s:25", server)); err == nil {
				return client, nil
			}
		}

		// None of the connections succeeded, so it is time to wait either for
		// the specified timeout duration or a receive on the stop channel
		var dur time.Duration
		if i < 2 {
			dur = 30 * time.Minute
		} else {
			dur = 2 * time.Hour
		}

		select {
		case <-time.After(dur):
		case <-stop:
			return nil, nil
		}
	}

	// All attempts have failed, let the caller know we tried :)
	return nil, errors.New("unable to connect to a mail server")
}

// Attempt to deliver the specified message to the server.
func deliverToMailServer(client *smtp.Client, from string, to []string, message []byte) error {
	if err := client.Mail(from); err != nil {
		return err
	}
	for _, recipient := range to {
		if err := client.Rcpt(recipient); err != nil {
			return err
		}
	}
	writer, err := client.Data()
	if err != nil {
		return err
	}
	defer writer.Close()
	if _, err = writer.Write(message); err != nil {
		return err
	}
	return nil
}

// Create a new host connection.
func NewHost(host string) *Host {

	// Create the host, including the channel used for stopping it
	h := &Host{
		name: host,
		stop: make(chan bool),
	}

	//...
	go func() {

		// Close the stop channel when the goroutine exits
		defer close(h.stop)
	}()

	return h
}

// Abort the connection.
func (h *Host) Stop() {

	// Send on the channel to stop it and wait for it to be closed
	h.stop <- true
	<-h.stop
}
