package smtputil

import (
	"crypto/tls"
	"io"
	"net/smtp"
)

//go:generate mockery -dir "." -case "underscore" -outpkg "smtpmocks" -output "../mocks/smtpmocks" -all

type Connecter interface {
	// SMTPConnect trying to connect to specific hostname.
	SMTPConnect(hostname string) (Client, error)
}

type Client interface {
	// Mail issues a MAIL command to the server using the provided email address.
	// If the server supports the 8BITMIME extension, Mail adds the BODY=8BITMIME
	// parameter.
	// This initiates a mail transaction and is followed by one or more Rcpt calls.
	Mail(from string) error
	// Rcpt issues a RCPT command to the server using the provided email address.
	// A call to Rcpt must be preceded by a call to Mail and may be followed by
	// a Data call or another Rcpt call.
	Rcpt(to string) error
	// Data issues a DATA command to the server and returns a writer that
	// can be used to write the mail headers and body. The caller should
	// close the writer before calling any more methods on c. A call to
	// Data must be preceded by one or more calls to Rcpt.
	Data() (io.WriteCloser, error)
	// Extension reports whether an extension is support by the server.
	// The extension name is case-insensitive. If the extension is supported,
	// Extension also returns a string that contains any parameters the
	// server specifies for the extension.
	Extension(ext string) (bool, string)
	// Reset sends the RSET command to the server, aborting the current mail
	// transaction.
	Reset() error
	// Noop sends the NOOP command to the server. It does nothing but check
	// that the connection to the server is okay.
	Noop() error
	// Quit sends the QUIT command and closes the connection to the server.
	Quit() error
	// Close closes the connection.
	Close() error
	// Hello sends a HELO or EHLO to the server as the given host name.
	// Calling this method is only necessary if the client needs control
	// over the host name used. The client will introduce itself as "localhost"
	// automatically otherwise. If Hello is called, it must be called before
	// any of the other methods.
	Hello(localName string) error
	// StartTLS sends the STARTTLS command and encrypts all further communication.
	// Only servers that advertise the STARTTLS extension support this function.
	StartTLS(config *tls.Config) error

	// TLSConnectionState returns the client's TLS connection state.
	// The return values are their zero values if StartTLS did
	// not succeed.
	TLSConnectionState() (state tls.ConnectionState, ok bool)

	// Verify checks the validity of an email address on the server.
	// If Verify returns nil, the address is valid. A non-nil return
	// does not necessarily indicate an invalid address. Many servers
	// will not verify addresses for security reasons.
	Verify(addr string) error

	// Auth authenticates a client using the provided authentication mechanism.
	// A failed authentication closes the connection.
	// Only servers that advertise the AUTH extension support this function.
	Auth(a smtp.Auth) error
}
