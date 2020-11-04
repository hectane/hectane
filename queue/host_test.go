package queue

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"io"
	"io/ioutil"
	"net/textproto"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/cenkalti/backoff/v4"
	nbc "github.com/hectane/go-nonblockingchan"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/hectane/hectane/internal/mocks/queuemocks"
	"github.com/hectane/hectane/internal/mocks/smtpmocks"
)

func newStorage(t *testing.T) (storage *Storage, deleter func()) {
	d, err := ioutil.TempDir(os.TempDir(), "")
	require.NoError(t, err)

	storage = NewStorage(d)
	return storage, func() {
		require.NoError(t, os.RemoveAll(d))
	}
}

func saveMessage(t *testing.T, messageBody string, storage *Storage, from string, to []string) *Message {
	r := strings.NewReader(messageBody)
	w, body, err := storage.NewBody()
	require.NoError(t, err)
	_, err = io.Copy(w, r)
	require.NoError(t, err)
	require.NoError(t, w.Close())

	m := &Message{
		From: from,
		To:   to,
	}

	require.NoError(t, storage.SaveMessage(m, body))
	return m
}

func TestHost_receiveMessage(t *testing.T) {
	store, deleter := newStorage(t)
	defer deleter()

	h := NewHost("example.com", store, new(Config))
	// we stop processor so it cannot interrupt in our test
	h.Stop()
	ctx, cancel := context.WithCancel(context.Background())
	h.ctx = ctx
	h.stopFunc = cancel

	require.True(t, h.lastActivity.IsZero())

	h.newMessage.Send <- &Message{
		id: "1",
	}

	m := h.receiveMessage()

	assert.True(t, h.lastActivity.IsZero())
	assert.Equal(t, "1", m.id)

	h.Stop()
	assert.Nil(t, h.receiveMessage())
}

func TestHost_receiveMessage_lastActivity(t *testing.T) {
	h := NewHost("example.com", new(Storage), new(Config))
	defer h.Stop()

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		wg.Done()
		_ = h.receiveMessage()
	}()

	time.Sleep(10 * time.Millisecond)
	assert.False(t, h.lastActivity.IsZero())
}

func TestHost_parseHostname(t *testing.T) {
	type args struct {
		addr string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "valid email",
			args: args{
				addr: "name@example.com",
			},
			want: "example.com",
		},
		{
			name: "with name",
			args: args{
				addr: "John <name@example.com>",
			},
			want: "example.com",
		},
		{
			name: "invalid email",
			args: args{
				addr: "name.example.com",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := Host{}
			got, err := h.parseHostname(tt.args.addr)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDefaultProcessor(t *testing.T) {
	r, w := io.Pipe()

	wg := sync.WaitGroup{}
	wg.Add(1)
	mailServerFinderMock := new(queuemocks.MailServerFinder)
	mailServerFinderMock.On("FindServers", "example.com").Return([]string{"mx1.example.com", "mx2.example.com"}, nil)
	clientMock := new(smtpmocks.Client)
	clientMock.On("Hello", "forwarder1.example.org").Return(nil)
	clientMock.On("Extension", "STARTTLS").Return(true, "")
	clientMock.On("StartTLS", mock.MatchedBy(func(conf *tls.Config) bool {
		assert.Equal(t, "mx1.example.com", conf.ServerName)
		assert.Equal(t, true, conf.InsecureSkipVerify)
		return true
	})).Return(nil)
	clientMock.On("Mail", "from@example.org").Return(nil)
	clientMock.On("Rcpt", "to1@example.com").Return(nil)
	clientMock.On("Data").Run(func(args mock.Arguments) {
		go func() {
			defer wg.Done()
			buf := bytes.Buffer{}
			_, err := io.Copy(&buf, r)
			require.NoError(t, err)
			assert.Equal(t, "some body1", buf.String())
		}()
	}).Return(w, nil)
	smtpConnecterMock := new(smtpmocks.Connecter)
	smtpConnecterMock.On("SMTPConnect", "mx1.example.com").Return(clientMock, nil)

	storage, deleter := newStorage(t)
	defer deleter()

	message := saveMessage(t, "some body1", storage, "from@example.org", []string{"to1@example.com"})

	h := Host{
		storage:          storage,
		mailServerFinder: mailServerFinderMock,
		smtpConnecter:    smtpConnecterMock,
		config: &Config{
			Hostname:               "forwarder1.example.org",
			DisableSSLVerification: true,
		},
	}

	err := h.defaultProcessor(message, storage)
	require.NoError(t, err)

	smtpConnecterMock.AssertExpectations(t)
	clientMock.AssertExpectations(t)
	mailServerFinderMock.AssertExpectations(t)
}

func TestDefaultProcessorTemporaryError(t *testing.T) {
	temporaryError := textproto.Error{
		Msg:  "Service is unavailable",
		Code: 421,
	}

	mailServerFinderMock := new(queuemocks.MailServerFinder)
	mailServerFinderMock.On("FindServers", "example.com").Return([]string{"mx1.example.com", "mx2.example.com"}, nil)
	clientMock := new(smtpmocks.Client)
	clientMock.On("Hello", "forwarder1.example.org").Return(&temporaryError)
	smtpConnecterMock := new(smtpmocks.Connecter)
	smtpConnecterMock.On("SMTPConnect", "mx1.example.com").Return(clientMock, nil)

	storage, deleter := newStorage(t)
	defer deleter()

	message := saveMessage(t, "some body1", storage, "from@example.org", []string{"to1@example.com"})

	h := Host{
		storage:          storage,
		mailServerFinder: mailServerFinderMock,
		smtpConnecter:    smtpConnecterMock,
		config: &Config{
			Hostname: "forwarder1.example.org",
		},
	}

	err := h.defaultProcessor(message, storage)
	require.Error(t, err)
	var e *SMTPError
	require.True(t, errors.As(err, &e))
	assert.True(t, e.IsTemporary())

	smtpConnecterMock.AssertExpectations(t)
	clientMock.AssertExpectations(t)
	mailServerFinderMock.AssertExpectations(t)
}

func TestConnectToMailServer(t *testing.T) {
	mailServerFinderMock := new(queuemocks.MailServerFinder)
	mailServerFinderMock.On("FindServers", "example.com").Return([]string{"mx1.example.com", "mx2.example.com"}, nil)
	smtpClientMock := new(smtpmocks.Client)
	smtpClientMock.On("Hello", "forwarder1.example.org").Return(nil)
	smtpClientMock.On("Extension", "STARTTLS").Return(true, "")
	smtpClientMock.On("StartTLS", mock.MatchedBy(func(conf *tls.Config) bool {
		assert.Equal(t, "mx1.example.com", conf.ServerName)
		assert.Equal(t, true, conf.InsecureSkipVerify)
		return true
	})).Return(nil)
	smtpConnecterMock := new(smtpmocks.Connecter)
	smtpConnecterMock.On("SMTPConnect", "mx1.example.com").Return(smtpClientMock, nil)

	h := Host{
		mailServerFinder: mailServerFinderMock,
		smtpConnecter:    smtpConnecterMock,
		config: &Config{
			Hostname:               "forwarder1.example.org",
			DisableSSLVerification: true,
		},
	}

	_, err := h.connectToMailServer("example.com")
	require.NoError(t, err)

	smtpConnecterMock.AssertExpectations(t)
	smtpClientMock.AssertExpectations(t)
	mailServerFinderMock.AssertExpectations(t)
}

func TestDeliverToMailServer(t *testing.T) {
	rPipe, wPipe := io.Pipe()

	wg := sync.WaitGroup{}
	wg.Add(1)
	clientMock := new(smtpmocks.Client)
	clientMock.On("Mail", "from@example.org").Return(nil)
	clientMock.On("Rcpt", "to@example.com").Return(nil)
	clientMock.On("Data").Run(func(args mock.Arguments) {
		go func() {
			defer wg.Done()
			buf := bytes.Buffer{}
			_, err := io.Copy(&buf, rPipe)
			require.NoError(t, err)
			assert.Equal(t, "some body", buf.String())
		}()
	}).Return(wPipe, nil)

	storage, deleter := newStorage(t)
	defer deleter()

	m := saveMessage(t, "some body", storage, "from@example.org", []string{"to@example.com"})

	h := Host{
		storage: storage,
		config:  &Config{},
	}

	require.NoError(t, h.deliverToMailServer(clientMock, m))

	clientMock.AssertExpectations(t)

	wg.Wait()
}

func TestRun(t *testing.T) {
	testCases := []struct {
		name       string
		beforeTest func(t *testing.T) (host *Host, postFunc func())
	}{
		{
			name: "success send",
			beforeTest: func(t *testing.T) (*Host, func()) {
				r, w := io.Pipe()

				mailServerFinderMock := new(queuemocks.MailServerFinder)
				mailServerFinderMock.On("FindServers", "example.com").Return([]string{"mx1.example.com", "mx2.example.com"}, nil)
				clientMock := new(smtpmocks.Client)
				clientMock.On("Hello", "forwarder1.example.org").Return(nil)
				clientMock.On("Extension", "STARTTLS").Return(true, "")
				clientMock.On("StartTLS", mock.MatchedBy(func(conf *tls.Config) bool {
					assert.Equal(t, "mx1.example.com", conf.ServerName)
					assert.Equal(t, true, conf.InsecureSkipVerify)
					return true
				})).Return(nil)
				clientMock.On("Mail", "from@example.org").Return(nil)
				clientMock.On("Rcpt", "to@example.com").Return(nil)
				clientMock.On("Data").Run(func(args mock.Arguments) {
					go func() {
						buf := bytes.Buffer{}
						_, err := io.Copy(&buf, r)
						require.NoError(t, err)
						assert.Equal(t, "some body1", buf.String())
					}()
				}).Return(w, nil)
				smtpConnecterMock := new(smtpmocks.Connecter)
				smtpConnecterMock.On("SMTPConnect", "mx1.example.com").Return(clientMock, nil)

				ctx, cancel := context.WithCancel(context.Background())

				storage, deleter := newStorage(t)
				hostWg := sync.WaitGroup{}
				hostWg.Add(1)
				h := Host{
					ctx:              ctx,
					log:              logrus.NewEntry(logrus.StandardLogger()),
					storage:          storage,
					wg:               &hostWg,
					newMessage:       nbc.New(),
					mailServerFinder: mailServerFinderMock,
					smtpConnecter:    smtpConnecterMock,
					config: &Config{
						Hostname:               "forwarder1.example.org",
						DisableSSLVerification: true,
					},
					back: &backoff.ZeroBackOff{},
				}
				h.process = h.defaultProcessor

				message := saveMessage(t, "some body1", h.storage, "from@example.org", []string{"to@example.com"})

				h.newMessage.Send <- message

				go func() {
					time.Sleep(100 * time.Millisecond)
					cancel()
				}()

				return &h, func() {
					deleter()
					smtpConnecterMock.AssertExpectations(t)
					clientMock.AssertExpectations(t)
					mailServerFinderMock.AssertExpectations(t)
				}
			},
		},
		{
			name: "success send after first retry",
			beforeTest: func(t *testing.T) (*Host, func()) {
				ctx, cancel := context.WithCancel(context.Background())
				r, w := io.Pipe()

				wg := sync.WaitGroup{}
				wg.Add(1)
				temporaryError := textproto.Error{
					Msg:  "Service is unavailable",
					Code: 421,
				}

				mailServerFinderMock := new(queuemocks.MailServerFinder)
				mailServerFinderMock.On("FindServers", "example.com").Return([]string{"mx1.example.com", "mx2.example.com"}, nil)
				clientMock := new(smtpmocks.Client)
				clientMock.On("Hello", "forwarder1.example.org").Return(&temporaryError).Once()
				clientMock.On("Hello", "forwarder1.example.org").Return(nil)
				clientMock.On("Extension", "STARTTLS").Return(true, "")
				clientMock.On("StartTLS", mock.MatchedBy(func(conf *tls.Config) bool {
					assert.Equal(t, "mx1.example.com", conf.ServerName)
					assert.Equal(t, true, conf.InsecureSkipVerify)
					return true
				})).Return(nil)
				clientMock.On("Mail", "from@example.org").Return(nil)
				clientMock.On("Rcpt", "to@example.com").Return(nil)
				clientMock.On("Data").Run(func(args mock.Arguments) {
					go func() {
						defer wg.Done()
						buf := bytes.Buffer{}
						_, err := io.Copy(&buf, r)
						require.NoError(t, err)
						assert.Equal(t, "some body1", buf.String())
						cancel()
					}()
				}).Return(w, nil)
				smtpConnecterMock := new(smtpmocks.Connecter)
				smtpConnecterMock.On("SMTPConnect", "mx1.example.com").Return(clientMock, nil)

				storage, deleter := newStorage(t)
				hostWg := sync.WaitGroup{}
				hostWg.Add(1)
				h := Host{
					ctx:              ctx,
					log:              logrus.NewEntry(logrus.StandardLogger()),
					storage:          storage,
					wg:               &hostWg,
					newMessage:       nbc.New(),
					mailServerFinder: mailServerFinderMock,
					smtpConnecter:    smtpConnecterMock,
					config: &Config{
						Hostname:               "forwarder1.example.org",
						DisableSSLVerification: true,
					},
					back: &backoff.ZeroBackOff{},
				}
				h.process = h.defaultProcessor

				message := saveMessage(t, "some body1", h.storage, "from@example.org", []string{"to@example.com"})

				h.newMessage.Send <- message

				return &h, func() {
					deleter()
					smtpConnecterMock.AssertExpectations(t)
					clientMock.AssertExpectations(t)
					mailServerFinderMock.AssertExpectations(t)
				}
			},
		},
		{
			name: "failed send after first try",
			beforeTest: func(t *testing.T) (*Host, func()) {
				ctx, cancel := context.WithCancel(context.Background())
				permanentError := textproto.Error{
					Msg:  "IP is blocked",
					Code: 525,
				}

				mailServerFinderMock := new(queuemocks.MailServerFinder)
				mailServerFinderMock.On("FindServers", "example.com").Return([]string{"mx1.example.com", "mx2.example.com"}, nil).Once()
				clientMock := new(smtpmocks.Client)
				clientMock.On("Hello", "forwarder1.example.org").Run(func(args mock.Arguments) {
					go func() {
						time.Sleep(100 * time.Millisecond)
						cancel()
					}()
				}).Return(&permanentError).Once()
				smtpConnecterMock := new(smtpmocks.Connecter)
				smtpConnecterMock.On("SMTPConnect", "mx1.example.com").Return(clientMock, nil).Once()

				storage, deleter := newStorage(t)
				hostWg := sync.WaitGroup{}
				hostWg.Add(1)
				h := Host{
					ctx:              ctx,
					log:              logrus.NewEntry(logrus.StandardLogger()),
					storage:          storage,
					wg:               &hostWg,
					newMessage:       nbc.New(),
					mailServerFinder: mailServerFinderMock,
					smtpConnecter:    smtpConnecterMock,
					config: &Config{
						Hostname:               "forwarder1.example.org",
						DisableSSLVerification: true,
					},
					back: &backoff.ZeroBackOff{},
				}
				h.process = h.defaultProcessor

				message := saveMessage(t, "some body1", h.storage, "from@example.org", []string{"to@example.com"})

				h.newMessage.Send <- message

				return &h, func() {
					deleter()
					smtpConnecterMock.AssertExpectations(t)
					clientMock.AssertExpectations(t)
					mailServerFinderMock.AssertExpectations(t)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.NotNil(t, tc.beforeTest, "beforeTest must be implemented")
			h, postFunc := tc.beforeTest(t)
			defer postFunc()
			h.run()
		})
	}
}

func TestSMTPError_IsTemporary(t *testing.T) {
	err := textproto.Error{
		Code: 421,
		Msg:  "The service is unavailable, try again later",
	}
	e := SMTPError{Err: &err}
	assert.False(t, e.IsPermanent())
	assert.True(t, e.IsTemporary())
}

func TestSMTPError_IsPermanent(t *testing.T) {
	err := textproto.Error{
		Code: 501,
		Msg:  "Syntax error",
	}
	e := SMTPError{Err: &err}
	assert.True(t, e.IsPermanent())
	assert.False(t, e.IsTemporary())
}
