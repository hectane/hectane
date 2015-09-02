package queue

import (
	"github.com/nathan-osman/go-cannon/util"

	"errors"
	"fmt"
	"net"
	"net/smtp"
	"syscall"
	"time"
)

// Individual message for sending to a host.
type Message struct {
	From    string
	To      []string
	Message []byte
}

// Persistent connection to an SMTP host.
type Host struct {
	newMessage *util.NonBlockingChan
	stop       chan bool
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

		// Enter a loop that will continue to deliver messages
	queue:
		for {

			// Obtain the next message for delivery
			var msg *Message

			select {
			case i := <-h.newMessage.Recv:
				msg = i.(*Message)
			case <-h.stop:
				break queue
			}

		connect:
			for {

				// Attempt to connect to a server
				client, err := connectToMailServer(host, h.stop)
				if client == nil {

					// Stop if there was no client and no error - otherwise,
					// discard the current message and enter the loop again
					if err == nil {
						break queue
					} else {
						// TODO: log something somewhere?
						// TODO: discard remaining messages?
						continue queue
					}
				}

			deliver:
				for {

					// Attempt to deliver the message
					if err = deliverToMailServer(client, msg); err != nil {

						// If the type of error has anything to do with a syscall,
						// assume that the connection was broken and try
						// reconnecting - otherwise, discard the message
						if _, ok := err.(syscall.Errno); ok {
							continue connect
						} else {
							continue queue
						}
					}

					// If a new message comes in then send it, if 5 minutes
					// elapses, close the connection and wait for a new message
					select {
					case i := <-h.newMessage.Recv:
						msg = i.(*Message)
						continue deliver
					case <-time.After(5 * time.Minute):
						client.Quit()
						continue queue
					case <-h.stop:
						break queue
					}
				}
			}
		}
	}()

	return h
}

// Attempt to deliver a message to the host.
func (h *Host) Deliver(msg *Message) {
	h.newMessage.Send <- msg
}

// Abort the connection to the host..
func (h *Host) Stop() {

	// Send on the channel to stop it and wait for it to be closed
	h.stop <- true
	<-h.stop
}
