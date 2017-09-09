package queue

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/smtp"
	"net/textproto"
	"os"
	"sort"
	"time"
)

var (
	resolver = &net.Resolver{}
	dialer   = &net.Dialer{
		Timeout:  30 * time.Second,
		Resolver: resolver,
	}
)

// findMailServers locates the MX records for a host and returns a list of them
// in order. If no mail servers were found, the host's name itself is used.
func findMailServers(ctx context.Context, host string) ([]string, error) {
	s, err := resolver.LookupMX(ctx, host)
	if err != nil {
		if err == context.Canceled {
			return nil, err
		}
		return []string{host}, nil
	}
	sort.Slice(s, func(i, j int) bool {
		return s[i].Pref < s[j].Pref
	})
	servers := []string{}
	for _, m := range s {
		servers = append(servers, m.Host)
	}
	return servers, nil
}

// tryMailServers tries the provided list of mail servers until one of them
// successfully connects or the entire list is exhausted.
func tryMailServers(ctx context.Context, servers []string) (*smtp.Client, error) {
	for _, s := range servers {
		conn, err := dialer.DialContext(ctx, "tcp", fmt.Sprintf("%s:25", s))
		if err != nil {
			if err == context.Canceled {
				return nil, err
			}
			continue
		}
		c, err := smtp.NewClient(conn, s)
		if err != nil {
			return nil, err
		}
		return c, nil
	}
	return nil, errors.New("unable to connect to a mail server")
}

// initServer sends the HELO message and enables TLS if available. Verification
// is disabled since so many popular hosts don't have proper certificates.
func initServer(c *smtp.Client) error {
	h, err := os.Hostname()
	if err != nil {
		return err
	}
	if err := c.Hello(h); err != nil {
		return err
	}
	if ok, _ := c.Extension("STARTTLS"); ok {
		if err := c.StartTLS(&tls.Config{InsecureSkipVerify: true}); err != nil {
			return err
		}
	}
	return nil
}

// isHostError determines if an error message pertains to the host (no other
// deliveries are likely to succeed now).
func isHostError(err error) bool {
	if e, ok := err.(*textproto.Error); ok {
		if e.Code >= 400 && e.Code < 500 {
			return false
		}
	}
	return true
}
