## Hectane

[![Build Status - Linux](https://travis-ci.org/hectane/hectane.svg)](https://travis-ci.org/hectane/hectane)
[![Build status - Windows](https://ci.appveyor.com/api/projects/status/h3r46k12llvw18u6?svg=true)](https://ci.appveyor.com/project/nathan-osman/hectane)
[![GoDoc](https://godoc.org/github.com/hectane/hectane?status.svg)](https://godoc.org/github.com/hectane/hectane)
[![MIT License](http://img.shields.io/badge/license-MIT-9370d8.svg?style=flat)](http://opensource.org/licenses/MIT)

The goal of Hectane is simple: to create an SMTP client that is extremely easy to configure and use. The application exposes a simple HTTP API that is used for sending emails.

### Features

- ability to attach files to emails
- support for TLS encryption and HTTP basic auth
- mail queue that efficiently delivers emails to hosts
- emails in the queue are stored on disk until delivery
- MX records for the destination host are tried in order of priority
- run the application as a service on Windows

### Parameters

Hectane currently recognizes the following command-line parameters (all are optional):

- `-config` - JSON file containing configuration
- `-bind` - address to bind to (in the format `host:port` - default is `:8025`)
- `-tls-cert` - path to TLS certificate
- `-tls-key` - path to TLS key
- `-username` - username for HTTP basic auth
- `-password` - password for HTTP basic auth
- `-directory` - storage location for emails awaiting delivery
- `-disable-ssl-verification` - disables verification of server SSL certificates
- `-logfile` - file to write log output to
- `-debug` - show debug log messages

### Usage

Hectane exposes an HTTP API that can be used to deliver emails. The API expects and responds with JSON data. Currently, the API consists of the following methods:

#### POST /v1/raw

##### Parameters

- `from` - sender email address
- `to` - list of recipient email addresses
- `body` - UTF-8 encoded message body

##### Response

The response is either an empty JSON object (indicating success) or a JSON object with an `error` key describing the problem.

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

#### GET /v1/status

##### Sample Response

    {
        "uptime": 1615,
        "hosts": {
            "example.com": {
                "active": false,
                "length": 2
            }
        }
    }

#### GET /v1/version

##### Sample Response

    {
        "version": "x.y.z"
    }
