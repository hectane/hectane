FROM golang:latest
MAINTAINER Nathan Osman <nathan@quickmediasolutions.com>

# Grab the source files and build them
RUN go get github.com/hectane/hectane

# Set a few configuration defaults
ENV DIRECTORY=/data \
        DISABLE_SSL_VERIFICATION=0 \
        LOGFILE=/var/log/hectane.log \
        DEBUG=0

# Specify the executable to run
CMD hectane \
        -tls-cert="$TLS_CERT" \
        -tls-key="$TLS_KEY" \
        -username="$USERNAME" \
        -password="$PASSWORD" \
        -directory="$DIRECTORY" \
        -disable-ssl-verification="$DISABLE_SSL_VERIFICATION" \
        -logfile="$LOGFILE" \
        -debug="$DEBUG"

# Expose the SMTP and HTTP API ports
EXPOSE 25
EXPOSE 8025
