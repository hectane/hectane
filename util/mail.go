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

// Group a list of email addresses by their host. An error will be returned if
// any of the addresses are invalid.
func GroupAddressesByHost(addrs []string) (map[string][]string, error) {
	m := make(map[string][]string)
	for _, a := range addrs {
		if addr, err := mail.ParseAddress(a); err == nil {
			parts := strings.Split(addr.Address, "@")
			if m[parts[1]] == nil {
				m[parts[1]] = make([]string, 0, 1)
			}
			m[parts[1]] = append(m[parts[1]], addr.Address)
		} else {
			return nil, err
		}
	}
	return m, nil
}
