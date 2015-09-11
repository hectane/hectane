package util

import (
	"net"
	"net/mail"
	"strings"
)

// Attempt to find the mail servers for the specified host. MX records are
// checked first. If one or more were found, the records are converted into an
// array of strings (sorted by priority). If none were found, the original host
// is returned.
func FindMailServers(host string) []string {
	if mx, err := net.LookupMX(host); err == nil {
		servers := make([]string, len(mx))
		for i, r := range mx {
			servers[i] = strings.TrimSuffix(r.Host, ".")
		}
		return servers
	} else {
		return []string{host}
	}
}

// Group a list of email addresses by their host.
func GroupAddressesByHost(addrs []string) (map[string][]string, error) {
	m := make(map[string][]string)
	for _, addr := range addrs {
		if host, err := HostFromAddress(addr); err != nil {
			return nil, err
		} else {
			if m[host] == nil {
				m[host] = make([]string, 0, 1)
			}
			m[host] = append(m[host], addr)
		}
	}
	return m, nil
}

// Attempt to extract the host from the specified email address.
func HostFromAddress(data string) (string, error) {
	if addr, err := mail.ParseAddress(data); err != nil {
		return "", err
	} else {
		return strings.Split(addr.Address, "@")[1], nil
	}
}
