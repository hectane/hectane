## go-cannon

[![MIT License](http://img.shields.io/badge/license-MIT-9370d8.svg?style=flat)](http://opensource.org/licenses/MIT)
[![GoDoc](https://godoc.org/github.com/nathan-osman/go-cannon?status.svg)](https://godoc.org/github.com/nathan-osman/go-cannon)
[![Build Status](https://travis-ci.org/nathan-osman/go-cannon.svg)](https://travis-ci.org/nathan-osman/go-cannon)

The goal of go-cannon is simple: to create an SMTP client that is extremely easy to configure and use. The application exposes a simple HTTP API that is used for sending emails.

### Features

- ability to attach files to emails
- support for TLS encryption and HTTP basic auth
- mail queue that efficiently delivers emails to hosts
- emails in the queue are stored on disk until delivery
- MX records for the destination host are tried in order of priority

### Parameters

go-cannon currently recognizes the following command-line parameters (all are optional):

- `-bind` - address to bind to (in the format `host:port` - default is `:8025`)
- `-tls-cert` - path to TLS certificate
- `-tls-key` - path to TLS key
- `-username` - username for HTTP basic auth
- `-password` - password for HTTP basic auth
- `-directory` - storage location for emails awaiting delivery

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
- `attachments` - a list of attachments:
    - `filename` - filename of the attachment
    - `content_type` - MIME type of the attachment (for example, `text/plain`)
    - `content` - UTF-8 or base64 encoded content of the attachment
    - `encoded` - `true` if `content` is base64 encoded

##### Response

The response is either an empty JSON object (indicating success) or a JSON object with an `error` key describing the problem.
