package queue

import (
	"github.com/nathan-osman/go-cannon/util"

	"errors"
	"fmt"
	"net"
	"net/smtp"
	"sync"
	"syscall"
	"time"
)

// Individual message for sending to a host.
type Message struct {
	Host    string
	From    string
	To      []string
	Message []byte
}

// Persistent connection to an SMTP host.
type Host struct {
	sync.Mutex
	lastActivity time.Time
	newMessage   *util.NonBlockingChan
	stop         chan bool
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
func deliverToMailServer(client *smtp.Client, msg *Message) error {

	// Specify the sender of the message
	if err := client.Mail(msg.From); err != nil {
		return err
	}

	// Add each of the recipients
	for _, to := range msg.To {
		if err := client.Rcpt(to); err != nil {
			return err
		}
	}

	// Obtain a writer for writing the actual message
	writer, err := client.Data()
	if err != nil {
		return err
	}
	defer writer.Close()

	// Write the message
	if _, err = writer.Write(msg.Message); err != nil {
		return err
	}

	return nil
}

// Create a new host connection.
func NewHost(host string) *Host {

	// Create the host, including the channel used for delivering new messages
	h := &Host{
		newMessage: util.NewNonBlockingChan(),
		stop:       make(chan bool),
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
			msg    *Message
		)

		// Receive a new message from the channel
	receive:

		// The connection (if one exists) is considered idle while waiting for
		// a new message to be received for delivery

		h.Lock()
		h.lastActivity = time.Now()
		h.Unlock()

		select {
		case i := <-h.newMessage.Recv:
			msg = i.(*Message)
		case <-h.stop:
			goto quit
		}

		h.Lock()
		h.lastActivity = time.Time{}
		h.Unlock()

		// Connect to the mail server (if not connected) and deliver a message
	connect:

		if client == nil {
			client, err := connectToMailServer(host, h.stop)
			if client == nil {

				// Stop if there was no client and no error - otherwise,
				// discard the current message and wait for the next one
				if err == nil {
					goto quit
				} else {
					// TODO: log something somewhere?
					// TODO: discard remaining messages?
					goto receive
				}
			}
		}

		// Attempt to deliver the message and then wait for the next one
		if err := deliverToMailServer(client, msg); err != nil {

			// If the type of error has anything to do with a syscall, assume
			// that the connection was broken and try reconnecting - otherwise,
			// discard the message
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

// Attempt to deliver a message to the host.
func (h *Host) Deliver(msg *Message) {
	h.newMessage.Send <- msg
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
