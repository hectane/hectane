package util

import (
	"net"
	"strings"
)

// Attempt to find the mail servers for the specified host.
func FindMailServers(host string) []string {

	// First check for MX records - if one or more were found, convert the
	// records into a list of strings (already sorted by priority) - if none
	// were found, then simply return the host that was originally provided
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
