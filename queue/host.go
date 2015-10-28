package queue

import (
	"github.com/Sirupsen/logrus"
	"github.com/hectane/hectane/util"

	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/smtp"
	"net/textproto"
	"os"
	"sync"
	"syscall"
	"time"
)

// Host status information.
type HostStatus struct {
	Active bool `json:"active"`
	Length int  `json:"length"`
}

// Persistent connection to an SMTP host.
type Host struct {
	m            sync.Mutex
	config       *Config
	storage      *Storage
	log          *logrus.Entry
	host         string
	newMessage   *util.NonBlockingChan
	lastActivity time.Time
	stop         chan bool
}

// Receive the next message in the queue. The host queue is considered
// "inactive" while waiting for new messages to arrive. The current time is
// recorded before entering the select{} block so that the Idle() method can
// calculate the idle time.
func (h *Host) receiveMessage() *Message {
	h.m.Lock()
	h.lastActivity = time.Now()
	h.m.Unlock()
	defer func() {
		h.m.Lock()
		h.lastActivity = time.Time{}
		h.m.Unlock()
	}()
	for {
		select {
		case i := <-h.newMessage.Recv:
			return i.(*Message)
		case <-h.stop:
			return nil
		}
	}
}

// Attempt to connect to the specified server. The connection attempt is
// performed in a separate goroutine, allowing it to be aborted if the host
// queue is shut down.
func (h *Host) tryMailServer(server string) (*smtp.Client, error) {
	var (
		c    *smtp.Client
		err  error
		done = make(chan bool)
	)
	go func() {
		c, err = smtp.Dial(fmt.Sprintf("%s:25", server))
		close(done)
	}()
	select {
	case <-done:
	case <-h.stop:
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if hostname, err := os.Hostname(); err == nil {
		if err := c.Hello(hostname); err != nil {
			return nil, err
		}
	}
	if ok, _ := c.Extension("STARTTLS"); ok {
		config := &tls.Config{ServerName: server}
		if h.config.DisableSSLVerification {
			config.InsecureSkipVerify = true
		}
		if err := c.StartTLS(config); err != nil {
			return nil, err
		}
	}
	return c, nil
}

// Attempt to connect to one of the mail servers.
func (h *Host) connectToMailServer() (*smtp.Client, error) {
	for _, s := range util.FindMailServers(h.host) {
		c, err := h.tryMailServer(s)
		if err != nil {
			h.log.Warningf("unable to connect to %s", s)
			continue
		}
		return c, nil
	}
	return nil, errors.New("unable to connect to a mail server")
}

// Attempt to send the specified message to the specified client.
func (h *Host) deliverToMailServer(c *smtp.Client, m *Message) error {
	r, err := h.storage.GetMessageBody(m)
	if err != nil {
		return err
	}
	defer r.Close()
	if err := c.Mail(m.From); err != nil {
		return err
	}
	for _, t := range m.To {
		if err := c.Rcpt(t); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	defer w.Close()
	if _, err := io.Copy(w, r); err != nil {
		return err
	}
	return nil
}

// Receive message and deliver them to their recipients. Due to the complicated
// algorithm for message delivery, the body of the method is broken up into a
// sequence of labeled sections.
func (h *Host) run() {
	defer close(h.stop)
	var (
		m        *Message
		c        *smtp.Client
		err      error
		tries    int
		duration time.Duration
	)
receive:
	if m == nil {
		m = h.receiveMessage()
		if m == nil {
			goto shutdown
		}
		h.log.Info("message received in queue")
	}
deliver:
	if c == nil {
		h.log.Info("connecting to mail server...")
		c, err = h.connectToMailServer()
		if c == nil {
			if err != nil {
				h.log.Error(err)
				goto wait
			} else {
				goto shutdown
			}
		}
		h.log.Info("connection established")
	}
	err = h.deliverToMailServer(c, m)
	if err != nil {
		h.log.Error(err)
		if _, ok := err.(syscall.Errno); ok {
			c = nil
			goto deliver
		}
		if e, ok := err.(*textproto.Error); ok {
			if e.Code >= 400 && e.Code <= 499 {
				c.Close()
				c = nil
				goto wait
			}
			c.Reset()
		}
		goto cleanup
	}
	h.log.Info("message delivered successfully")
cleanup:
	h.log.Info("deleting message from disk")
	err = h.storage.DeleteMessage(m)
	if err != nil {
		h.log.Error(err.Error())
	}
	m = nil
	tries = 0
	goto receive
wait:
	tries++
	// We differ a tiny bit from the RFC spec here but this should work well
	// enough - retry once after a minute, twice on the half-hour, and 16 more
	// times every three hours. This is roughly 48 hours.
	switch {
	case tries == 1:
		duration = time.Minute
	case tries < 4:
		duration = 30 * time.Minute
	case tries < 20:
		duration = 3 * time.Hour
	default:
		h.log.Warning("maximum retry count exceeded")
		goto cleanup
	}
	select {
	case <-h.stop:
	case <-time.After(duration):
		goto receive
	}
shutdown:
	h.log.Info("shutting down")
	if c != nil {
		c.Close()
	}
}

// Create a new host connection.
func NewHost(host string, s *Storage, c *Config) *Host {
	h := &Host{
		config:     c,
		storage:    s,
		log:        logrus.WithField("context", host),
		host:       host,
		newMessage: util.NewNonBlockingChan(),
		stop:       make(chan bool),
	}
	go h.run()
	return h
}

// Attempt to deliver a message to the host.
func (h *Host) Deliver(m *Message) {
	h.newMessage.Send <- m
}

// Retrieve the connection idle time.
func (h *Host) Idle() time.Duration {
	h.m.Lock()
	defer h.m.Unlock()
	if h.lastActivity.IsZero() {
		return 0
	} else {
		return time.Since(h.lastActivity)
	}
}

// Return the status of the host connection.
func (h *Host) Status() *HostStatus {
	return &HostStatus{
		Active: h.Idle() == 0,
		Length: h.newMessage.Len(),
	}
}

// Close the connection to the host.
func (h *Host) Stop() {
	h.stop <- true
	<-h.stop
}
