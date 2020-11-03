package queue

import (
	"context"
	"fmt"
	"net"
	"net/smtp"
	"time"

	"github.com/hectane/hectane/internal/smtputil"
)

type connector struct {
	helloHostname          string
	ctx                    context.Context
	disableSSLVerification bool
}

func newConnector(ctx context.Context, helloHostname string, disableSSLVerification bool) *connector {
	return &connector{
		helloHostname:          helloHostname,
		ctx:                    ctx,
		disableSSLVerification: disableSSLVerification,
	}
}

var _ smtputil.Connecter = new(connector)

func (c *connector) SMTPConnect(hostname string) (smtputil.Client, error) {
	var (
		err    error
		client smtputil.Client
		errCh  = make(chan error, 1)
	)

	go func() {
		var conn net.Conn

		addr := fmt.Sprintf("%s:25", hostname)
		deadline, ok := c.ctx.Deadline()
		if ok {
			conn, err = net.DialTimeout("tcp", addr, deadline.Sub(time.Now()))
			if err != nil {
				errCh <- err
				return
			}
		} else {
			conn, err = net.Dial("tcp", addr)
			if err != nil {
				errCh <- err
				return
			}
		}

		client, err = smtp.NewClient(conn, hostname)
		if err != nil {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	select {
	case err := <-errCh:
		if err != nil {
			return nil, err
		}
	case <-c.ctx.Done():
		return nil, c.ctx.Err()
	}

	return client, nil
}
