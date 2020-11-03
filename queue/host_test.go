package queue

import (
	"bytes"
	"context"
	"crypto/tls"
	"io"
	"io/ioutil"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/hectane/hectane/internal/mocks/queuemocks"
	"github.com/hectane/hectane/internal/mocks/smtpmocks"
)

func TestHost_receiveMessage(t *testing.T) {
	d, err := ioutil.TempDir(os.TempDir(), "")
	require.NoError(t, err)

	defer func() {
		require.NoError(t, os.RemoveAll(d))
	}()
	store := NewStorage(d)

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
	smtpClientMock.On("Mail", "from@example.org").Return(nil)
	smtpClientMock.On("Rcpt", "to1@example.com").Return(nil)
	smtpClientMock.On("Data").Run(func(args mock.Arguments) {
		go func() {
			buf := bytes.Buffer{}
			_, err := io.Copy(&buf, r)
			require.NoError(t, err)
		}()
	}).Return(w, nil)
	smtpConnecterMock := new(smtpmocks.Connecter)
	smtpConnecterMock.On("SMTPConnect", "example.com").Return(smtpClientMock, nil)

	h := Host{
		mailServerFinder: mailServerFinderMock,
	}

	message := Message{
		From: "name@example.com",
	}
	storage := Storage{}

	err := h.defaultProcessor(&message, &storage)
	require.NoError(t, err)

	smtpConnecterMock.AssertExpectations(t)
	smtpClientMock.AssertExpectations(t)
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
