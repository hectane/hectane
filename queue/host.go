package queue

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/mail"
	"net/textproto"
	"strings"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/hectane/go-nonblockingchan"
	"github.com/sirupsen/logrus"

	"github.com/hectane/hectane/internal/smtputil"
)

type SMTPError struct {
	Err error
}

func (e *SMTPError) Error() string {
	return e.Err.Error()
}

func (e *SMTPError) Unwrap() error {
	return e.Err
}

func (e *SMTPError) IsPermanent() bool {
	var err *textproto.Error
	if errors.As(e.Err, &err) {
		return 500 <= err.Code && err.Code <= 599
	}

	return false
}

func (e *SMTPError) IsTemporary() bool {
	var err *textproto.Error
	if errors.As(e.Err, &err) {
		return 400 <= err.Code && err.Code <= 499
	}

	return false
}

// Host status information.
type HostStatus struct {
	Active bool `json:"active"`
	Length int  `json:"length"`
}

type ProcessFunc func(m *Message, s *Storage) error

// Persistent connection to an SMTP host.
type Host struct {
	m            sync.RWMutex
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
	back             backoff.BackOff
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

type dnsServerFinder struct{}

// FindServers looking for the mail servers for the specified host. MX records are
// checked first. If one or more were found, the records are converted into an
// array of strings (sorted by priority).
func (d *dnsServerFinder) FindServers(host string) ([]string, error) {
	r, err := net.LookupMX(host)
	if err != nil {
		return nil, err
	}
	servers := make([]string, len(r))
	for i, r := range r {
		servers[i] = strings.TrimSuffix(r.Host, ".")
	}
	return servers, nil
}

var _ MailServerFinder = new(dnsServerFinder)

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
		m   *Message
		c   smtputil.Client
		err error
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
		var smtpErr *SMTPError
		if errors.As(err, &smtpErr) && smtpErr.IsPermanent() {
			h.log.WithError(err).Error("got permanent error")
			goto cleanup
		}
		h.log.WithError(err).Error("failed to process message")
		goto wait
	}
	h.log.Info("message delivered successfully")
	goto cleanup
cleanup:
	h.log.Debug("deleting message from disk")
	err = h.storage.DeleteMessage(m)
	if err != nil {
		h.log.Error(err.Error())
	}
	m = nil
	h.back.Reset()
	goto receive
wait:
	duration := h.back.NextBackOff()
	if duration == backoff.Stop {
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
		var protoError *textproto.Error
		if errors.As(err, &protoError) {
			return &SMTPError{Err: err}
		}
		return err
	}

	if err := h.deliverToMailServer(c, m); err != nil {
		var protoError *textproto.Error
		if errors.As(err, &protoError) {
			return &SMTPError{Err: err}
		}
		return err
	}
	if err := c.Quit(); err != nil {
		h.log.WithError(err).Warning("failed to send quit command")
	}

	return nil
}

// NewHost creates a new host connection.
func NewHost(host string, s *Storage, c *Config) *Host {
	ctx, cancel := context.WithCancel(context.Background())

	back := backoff.NewExponentialBackOff()
	back.InitialInterval = c.InitialInterval
	back.MaxElapsedTime = c.MaxElapsedTime
	back.MaxInterval = c.MaxInterval
	back.Multiplier = c.Multiplier
	back.RandomizationFactor = c.RandomizationFactor

	h := &Host{
		config:           c,
		storage:          s,
		log:              logrus.WithField("context", host),
		host:             host,
		newMessage:       nbc.New(),
		ctx:              ctx,
		stopFunc:         cancel,
		wg:               &sync.WaitGroup{},
		process:          c.ProcessFunc,
		mailServerFinder: &dnsServerFinder{},
		smtpConnecter:    newConnector(ctx, c.Hostname, c.DisableSSLVerification),
		back:             back,
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
	h.m.RLock()
	defer h.m.RUnlock()
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
