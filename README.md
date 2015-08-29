## go-cannon

The goal of go-cannon is simple: to create an SMTP client that is extremely easy to configure and use. The application exposes a simple HTTP API that is used for sending emails. go-cannon takes care of looking up the MX records for the recipient and delivering the message.

### Parameters

go-cannon currently recognizes the following command-line parameters:

- `-username` - username for HTTP basic auth
- `-password` - password for HTTP basic auth

### Usage

go-cannon exposes an HTTP API that can be used to deliver emails. All API methods expect parameters to be sent as JSON data in a POST request. Currently, the API consists of a single method:

#### /v1/send

- `from` - sender email address
- `to` - recipient email address
- `subject` - email subject
- `text` - body of the email as plain text
- `html` - body of the email as HTML

If an error occurs, a JSON object similar to the following will be returned:

    {
        "error": "brief error description"
    }

### Planned Features

The following features are planned for a future release:

- TLS
- file attachments
- mail queue for when delivery fails
