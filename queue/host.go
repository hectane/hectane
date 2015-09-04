package queue

import (
	"github.com/nathan-osman/go-cannon/email"
	"github.com/nathan-osman/go-cannon/util"

	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/smtp"
	"os"
	"sync"
	"syscall"
	"time"
)

// Persistent connection to an SMTP host.
type Host struct {
	sync.Mutex
	lastActivity time.Time
	newEmail     *util.NonBlockingChan
	stop         chan bool
}

// Attempt to find the mail servers for the specified host.
func findMailServers(host string) []string {

	// First check for MX records - if one or more were found, convert the
	// records into a list of strings (already sorted by priority) - if none
	// were found, then simply return the host that was originally provided
	if mx, err := net.LookupMX(host); err == nil {
		servers := make([]string, len(mx))
		for i, r := range mx {
			servers[i] = r.Host
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
		for _, s := range servers {
			if client, err := smtp.Dial(fmt.Sprintf("%s:25", s)); err == nil {

				// Obtain the device's current hostname and say HELO
				var hostname string
				if hostname, err := os.Hostname(); err == nil {
					client.Hello(hostname)
				}

				// Check for TLS and enable if possible
				if ok, _ := client.Extension("STARTTLS"); ok {
					if err := client.StartTLS(&tls.Config{ServerName: hostname}); err != nil {
						continue
					}
				}

				// Return the client
				return client, nil
			}
		}

		// None of the connections succeeded, so it is time to wait either for
		// the specified timeout duration or a receive on the stop channel
		var d time.Duration
		if i < 2 {
			d = 30 * time.Minute
		} else {
			d = 2 * time.Hour
		}

		select {
		case <-time.After(d):
		case <-stop:
			return nil, nil
		}
	}

	// All attempts have failed, let the caller know we tried :)
	return nil, errors.New("unable to connect to a mail server")
}

// Attempt to deliver the specified email to the server.
func deliverToMailServer(client *smtp.Client, e *email.Email) error {

	// Specify the sender of the emails
	if err := client.Mail(e.From); err != nil {
		return err
	}

	// Add each of the recipients
	for _, t := range e.To {
		if err := client.Rcpt(t); err != nil {
			return err
		}
	}

	// Obtain a writer for writing the actual email
	w, err := client.Data()
	if err != nil {
		return err
	}
	defer w.Close()

	// Write the email
	if _, err = w.Write(e.Message); err != nil {
		return err
	}

	return nil
}

// Create a new host connection.
func NewHost(host string) *Host {

	// Create the host, including the channel used for delivering new emails
	h := &Host{
		newEmail: util.NewNonBlockingChan(),
		stop:     make(chan bool),
	}

	// Start a goroutine to manage the lifecycle of the host connection
	go func() {

		// Close the stop channel when the goroutine exits
		defer close(h.stop)

		// If the sight of a "goto" statement makes you cringe, you should
		// probably close your eyes and skip over the next section. Although
		// what follows could in theory be written without "goto", it wouldn't
		// be nearly as concise or easy to follow. Therefore it is used here.

		var (
			client *smtp.Client
			e      *email.Email
		)

		// Receive a new email from the channel
	receive:

		// The connection (if one exists) is considered idle while waiting for
		// a new email to be received for delivery

		h.Lock()
		h.lastActivity = time.Now()
		h.Unlock()

		select {
		case i := <-h.newEmail.Recv:
			e = i.(*email.Email)
		case <-h.stop:
			goto quit
		}

		h.Lock()
		h.lastActivity = time.Time{}
		h.Unlock()

		// Connect to the mail server (if not connected) and deliver an email
	connect:

		if client == nil {
			var err error
			client, err = connectToMailServer(host, h.stop)
			if client == nil {

				// Stop if there was no client and no error - otherwise,
				// discard the current email and wait for the next one
				if err == nil {
					goto quit
				} else {
					// TODO: log something somewhere?
					// TODO: discard remaining emails?
					goto receive
				}
			}
		}

		// Attempt to deliver the email and then wait for the next one
		if err := deliverToMailServer(client, e); err != nil {

			// If the type of error has anything to do with a syscall, assume
			// that the connection was broken and try reconnecting - otherwise,
			// discard the email
			if _, ok := err.(syscall.Errno); ok {
				client = nil
				goto connect
			}
		}

		// Receive the next message
		goto receive

		// Close the connection (if open) and quit
	quit:

		if client != nil {
			client.Quit()
		}
	}()

	return h
}

// Attempt to deliver an email to the host.
func (h *Host) Deliver(e *email.Email) {
	h.newEmail.Send <- e
}

// Retrieve the connection idle time.
func (h *Host) Idle() time.Duration {
	h.Lock()
	defer h.Unlock()
	if h.lastActivity.IsZero() {
		return 0
	} else {
		return time.Since(h.lastActivity)
	}
}

// Close the connection to the host.
func (h *Host) Stop() {

	// Send on the channel to stop it and wait for it to be closed
	h.stop <- true
	<-h.stop
}
