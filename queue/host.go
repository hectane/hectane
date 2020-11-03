package queue

import (
	"context"

	"github.com/hectane/go-nonblockingchan"
	"github.com/sirupsen/logrus"

	"crypto/tls"
	"fmt"
	"github.com/hectane/hectane/internal/smtputil"
	"io"
	"net"
	"net/mail"
	"net/smtp"
	"net/textproto"
	"strings"
	"sync"
	"syscall"
	"time"
)

// Host status information.
type HostStatus struct {
	Active bool `json:"active"`
	Length int  `json:"length"`
}

type ProcessFunc func(m *Message, s *Storage) error

// Persistent connection to an SMTP host.
type Host struct {
	m            sync.Mutex
	config       *Config
	storage      *Storage
	log          *logrus.Entry
	host         string
	newMessage   *nbc.NonBlockingChan
	lastActivity time.Time
	ctx          context.Context
	stopFunc     context.CancelFunc
	wg           *sync.WaitGroup
	process      ProcessFunc

	mailServerFinder MailServerFinder
	smtpConnecter    smtputil.Connecter
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
	select {
	case m := <-h.newMessage.Recv:
		return m.(*Message)
	case <-h.ctx.Done():
		return nil
	}
}

// Parse an email address and extract the hostname.
func (h *Host) parseHostname(addr string) (string, error) {
	a, err := mail.ParseAddress(addr)
	if err != nil {
		return "", err
	}
	return strings.Split(a.Address, "@")[1], nil
}

// Attempt to connect to the specified server. The connection attempt is
// performed in a separate goroutine, allowing it to be aborted if the host
// queue is shut down.
func (h *Host) tryMailServer(server, hostname string) (*smtp.Client, error) {
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
	case <-h.ctx.Done():
		return nil, h.ctx.Err()
	}
	if err != nil {
		return nil, err
	}
	if err := c.Hello(hostname); err != nil {
		return nil, err
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

// Attempt to find the mail servers for the specified host. MX records are
// checked first. If one or more were found, the records are converted into an
// array of strings (sorted by priority). If none were found, the original host
// is returned.
func (h *Host) findMailServers(host string) []string {
	r, err := net.LookupMX(host)
	if err != nil {
		return []string{host}
	}
	servers := make([]string, len(r))
	for i, r := range r {
		servers[i] = strings.TrimSuffix(r.Host, ".")
	}
	return servers
}

// Attempt to connect to one of the mail servers.
func (h *Host) connectToMailServer(hostname string) (smtputil.Client, error) {
	mxServers, err := h.mailServerFinder.FindServers(hostname)
	if err != nil {
		return nil, err
	}
	for _, mxServer := range mxServers {
		c, err := h.smtpConnecter.SMTPConnect(mxServer)
		if err != nil {
			h.log.WithError(err).Debugf("unable to connect to %s", mxServer)
			continue
		}

		if err := c.Hello(h.config.Hostname); err != nil {
			return nil, err
		}

		if ok, _ := c.Extension("STARTTLS"); ok {
			config := &tls.Config{ServerName: mxServer}
			if h.config.DisableSSLVerification {
				config.InsecureSkipVerify = true
			}
			if err := c.StartTLS(config); err != nil {
				return nil, err
			}
		}
		return c, nil
	}
	return nil, fmt.Errorf("unable to reach any mail server for %s", hostname)
}

// Attempt to send the specified message to the specified client.
func (h *Host) deliverToMailServer(c smtputil.Client, m *Message) error {
	r, err := h.storage.GetMessageBody(m)
	if err != nil {
		return err
	}
	defer r.Close()
	r, err = dkimSigned(m.From, r, h.config)
	if err != nil {
		return err
	}
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
	var (
		m        *Message
		hostname string
		c        smtputil.Client
		err      error
		tries    int
		duration = time.Minute
	)

	defer func() {
		h.log.Debug("shutting down")
		if c != nil {
			c.Close()
		}
		h.wg.Done()
	}()

receive:
	if m == nil {
		m = h.receiveMessage()
		if m == nil {
			return
		}
		h.log.Info("message received in queue")
	}
	if err := h.process(m, h.storage); err != nil {
		h.log.WithError(err).Error("failed to process message")
		goto wait
	} else {
		goto cleanup
	}
	hostname, err = h.parseHostname(m.From)
	if err != nil {
		h.log.Error(err)
		goto cleanup
	}

deliver:
	if c == nil {
		h.log.Debug("connecting to mail server")
		c, err = h.connectToMailServer(hostname)
		if c == nil {
			if err != nil {
				h.log.Error(err)
				goto wait
			} else {
				return
			}
		}
		h.log.Debug("connection established")
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
				if closeError := c.Close(); closeError != nil {
					h.log.WithError(err).Error("close error")
				}
				c = nil
				goto wait
			}
			if rstErr := c.Reset(); rstErr != nil {
				h.log.WithError(err).Error("reset error")
			}
		}
		h.log.Error(err.Error())
		goto cleanup
	}
	h.log.Info("message delivered successfully")
cleanup:
	h.log.Debug("deleting message from disk")
	err = h.storage.DeleteMessage(m)
	if err != nil {
		h.log.Error(err.Error())
	}
	m = nil
	tries = 0
	goto receive
wait:
	// We differ a tiny bit from the RFC spec here but this should work well
	// enough - the goal is to retry lots of times early on and space out the
	// remaining attempts as time goes on. (Roughly 48 hours total.)
	tries++
	switch {
	case tries < 8:
		duration *= 2
	case tries < 18:
	default:
		h.log.Error("maximum retry count exceeded")
		goto cleanup
	}
	select {
	case <-h.ctx.Done():
	case <-time.After(duration):
		goto receive
	}
}

func (h *Host) defaultProcessor(m *Message, s *Storage) error {
	hostname, err := h.parseHostname(m.To[0])
	if err != nil {
		return err
	}

	c, err := h.connectToMailServer(hostname)
	if err != nil {
		return err
	}

	if err := h.deliverToMailServer(c, m); err != nil {
		return err
	}

	return nil
}

// NewHost creates a new host connection.
func NewHost(host string, s *Storage, c *Config) *Host {
	ctx, cancel := context.WithCancel(context.Background())
	h := &Host{
		config:     c,
		storage:    s,
		log:        logrus.WithField("context", host),
		host:       host,
		newMessage: nbc.New(),
		ctx:        ctx,
		stopFunc:   cancel,
		wg:         &sync.WaitGroup{},
		process:    c.ProcessFunc,
	}
	if h.process == nil {
		h.process = h.defaultProcessor
	}

	h.wg.Add(1)
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
	}
	return time.Since(h.lastActivity)
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
	h.stopFunc()
	h.wg.Wait()
}
