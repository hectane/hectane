package queue

import (
	"github.com/nathan-osman/go-cannon/util"

	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"sync"
	"syscall"
	"time"
)

// Persistent connection to an SMTP host.
type Host struct {
	sync.Mutex
	host         string
	lastActivity time.Time
	newMessage   *util.NonBlockingChan
	stop         chan bool
}

// Log the specified message for the specified host.
func (h *Host) log(msg string) {
	log.Printf("[%s] %s", h.host, msg)
}

// Receive the next message in the queue.
func (h *Host) receiveMessage() *Message {

	// The host queue is considered "inactive" while waiting for new messages
	// to arrive - record the current time before entering the select{} block
	// so that the Idle() method can calculate the idle time
	h.Lock()
	h.lastActivity = time.Now()
	h.Unlock()

	// When the function exits, reset the inactive timer to 0
	defer func() {
		h.Lock()
		h.lastActivity = time.Time{}
		h.Unlock()
	}()

	// Either receive a new message or stop the queue
	select {
	case i := <-h.newMessage.Recv:
		return i.(*Message)
	case <-h.stop:
		return nil
	}
}

// Attempt to connect to the specified server.
func (h *Host) tryMailServer(server string) (*smtp.Client, error) {

	var (
		c    *smtp.Client
		err  error
		done = make(chan bool)
	)

	// Because Dial() is a blocking function, it must be run in a separate
	// goroutine so that it can be aborted immediately
	go func() {
		c, err = smtp.Dial(fmt.Sprintf("%s:25", server))
		close(done)
	}()

	// Wait for either the goroutine to complete or a stop request
	select {
	case <-done:
	case <-h.stop:
		return nil, nil
	}

	// Attempt to establish a TCP connection to port 25
	if err == nil {

		// Obtain this machine's hostname and say HELO
		if hostname, err := os.Hostname(); err == nil {
			c.Hello(hostname)
		}

		// If the server advertises TLS, attempt to use it - if the server
		// advertises TLS but it doesn't work, that's an error
		if ok, _ := c.Extension("STARTTLS"); ok {
			if err := c.StartTLS(&tls.Config{ServerName: server}); err != nil {
				return nil, err
			}
		}

		return c, nil

	} else {
		return nil, err
	}
}

// Attempt to connect to a mail server until the timeout is reached or until
// stopped.
func (h *Host) connectToMailServer() (*smtp.Client, error) {

	// Obtain the list of mail servers to try
	servers := util.FindMailServers(h.host)

	// RFC 5321 (4.5.4) describes the process for retrying connections to a
	// mail server after failure. The recommended strategy is to retry twice
	// with 30 minute intervals and continue at 120 minute intervals until four
	// days have elapsed. That's *roughly* what is done here.
	for i := 0; i < 50; i++ {

		// Try each of the servers in order
		for _, s := range servers {
			if c, err := h.tryMailServer(s); err == nil {
				return c, nil
			} else {
				h.log(fmt.Sprintf("unable to connect to %s", s))
			}
		}

		// None of the servers could be reached - wait a few minutes before
		// trying to connect again
		var d time.Duration
		if i < 2 {
			d = 30 * time.Minute
		} else {
			d = 2 * time.Hour
		}

		select {
		case <-time.After(d):
		case <-h.stop:
			return nil, nil
		}
	}

	// All attempts have failed, let the caller know we tried :)
	return nil, errors.New("unable to connect to a mail server")
}

// Attempt to deliver the specified message to the specifed client
func (h *Host) deliverToMailServer(c *smtp.Client, m *Message) error {

	// Specify the sender of the message
	if err := c.Mail(m.From); err != nil {
		return err
	}

	// Add each of the recipients
	for _, t := range m.To {
		if err := c.Rcpt(t); err != nil {
			return err
		}
	}

	// Obtain a writer for writing the actual message
	w, err := c.Data()
	if err != nil {
		return err
	}
	defer w.Close()

	// Write the message
	if _, err = w.Write(m.Message); err != nil {
		return err
	}

	return nil
}

// Receive message and deliver them to their recipients.
func (h *Host) run() {

	// Close the stop channel when the goroutine exits
	defer close(h.stop)

	// The client must be declared here so that it can be closed after the loop
	var (
		c   *smtp.Client
		m   *Message
		err error
	)

	for {

		// Receive the next message from the queue if one hasn't already been
		// retrieved
		if m == nil {
			h.log("waiting for message in queue...")
			if m = h.receiveMessage(); m == nil {
				break
			}
			h.log("message received in queue")
		}

		// Connect to the server if a connection does not already exist
		if c == nil {
			h.log("connecting to mail server...")
			if c, err = h.connectToMailServer(); c == nil {
				if err == nil {
					break
				} else {
					h.log(err.Error())
					continue
				}
			}
			h.log("connection established")
		}

		// Attempt delivery of the message
		if err = h.deliverToMailServer(c, m); err != nil {

			// If the type of error has anything to do with a syscall, assume
			// that the connection was broken and try reconnecting - otherwise,
			// discard the message - either way, reset the error
			if _, ok := err.(syscall.Errno); ok {
				h.log("connection to server lost")
				c, err = nil, nil
			} else {
				h.log(err.Error())
				m.Delete()
				m, err = nil, nil
			}
			continue

		} else {
			h.log("message successfully delivered")
			m.Delete()
			m = nil
		}
	}

	// Close the connection if it is still open
	if c != nil {
		c.Quit()
	}

	h.log("shutting down queue")
}

// Create a new host connection.
func NewHost(host string) *Host {

	// Create the host, including the channel used for delivering new messages
	h := &Host{
		host:       host,
		newMessage: util.NewNonBlockingChan(),
		stop:       make(chan bool),
	}

	// Start a goroutine to manage the lifecycle of the host connection
	go h.run()

	return h
}

// Attempt to deliver a message to the host.
func (h *Host) Deliver(m *Message) {
	h.newMessage.Send <- m
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
