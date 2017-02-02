## Hectane

[![Build Status - Linux](https://travis-ci.org/hectane/hectane.svg)](https://travis-ci.org/hectane/hectane)
[![Build status - Windows](https://ci.appveyor.com/api/projects/status/h3r46k12llvw18u6?svg=true)](https://ci.appveyor.com/project/nathan-osman/hectane)
[![GoDoc](https://godoc.org/github.com/hectane/hectane?status.svg)](https://godoc.org/github.com/hectane/hectane)
[![MIT License](http://img.shields.io/badge/license-MIT-9370d8.svg?style=flat)](http://opensource.org/licenses/MIT)

Hectane is both a Go package providing an SMTP queue for sending emails and a standalone application that exposes this functionality via an HTTP API.

### Features

- Ability to attach files to emails
- Support for TLS encryption and HTTP basic auth
- Mail queue that efficiently delivers emails to hosts
- Emails in the queue are stored on disk until delivery
- MX records for the destination host are tried in order of priority
- Run the application as a service on Windows

### Documentation

Documentation for Hectane can be found below:

- [Using Hectane in a Go application](https://github.com/hectane/hectane/wiki/Hectane%20Package)
- [Using Hectane in another language or on a server](https://github.com/hectane/hectane/wiki/Hectane%20Daemon)

### Installation

In addition to the [files on the releases page](https://github.com/hectane/hectane/releases), Hectane can be installed from any of the sources below:

- PPA: [stable](https://launchpad.net/~hectane/+archive/ubuntu/hectane) | [daily](https://launchpad.net/~hectane/+archive/ubuntu/hectane-dev)
- [Juju charm store](https://jujucharms.com/hectane/)
- [Docker Hub](https://hub.docker.com/r/hectane/hectane/)
