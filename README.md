## go-cannon

The goal of go-cannon is simple: to create an SMTP client that is extremely easy to configure and use. The application exposes a simple HTTP API that is used for sending emails.

### Features

- support for TLS encryption and HTTP basic auth
- mail queue that efficiently delivers emails to hosts
- MX records for the destination host are tried in order of priority

### Parameters

go-cannon currently recognizes the following command-line parameters (all are optional):

- `-bind` - address to bind to (in the format `host:port` - default is `:8025`)
- `-tls-cert` - path to TLS certificate
- `-tls-key` - path to TLS key
- `-username` - username for HTTP basic auth
- `-password` - password for HTTP basic auth

### Usage

go-cannon exposes an HTTP API that can be used to deliver emails. The API expects and responds with JSON data. Currently, the API consists of the following methods:

#### GET /v1/version

##### Response

    {
        "version": "x.y.z"
    }

#### POST /v1/send

##### Parameters

- `from` - sender email address
- `to` - list of recipient email addresses
- `cc` - list of carbon copy recipients
- `bcc` - list of blind carbon copy recipients
- `subject` - email subject
- `text` - body of the email as plain text
- `html` - body of the email as HTML

##### Response

- `status` - one of `delivered` or `error`

### Planned Features

The following features are planned for a future release:

- file attachments
- persisting mail queue to disk
