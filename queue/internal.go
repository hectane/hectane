package queue

//go:generate mockery -dir "." -case "underscore" -outpkg "queuemocks" -output "../internal/mocks/queuemocks" -all

type MailServerFinder interface {
	// Attempt to find the mail servers for the specified host. MX records are
	// checked first. If one or more were found, the records are converted into an
	// array of strings (sorted by priority). If none were found, the original host
	// is returned.
	FindServers(host string) ([]string, error)
}
