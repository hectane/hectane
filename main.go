package main

import (
	"github.com/hectane/hectane/api"
	"github.com/hectane/hectane/queue"
	"github.com/mitchellh/go-homedir"

	"flag"
	"log"
	"os"
	"path"
)

func main() {
	var (
		bind      string
		tlsCert   string
		tlsKey    string
		username  string
		password  string
		directory string
		config    queue.Config
	)
	if home, err := homedir.Dir(); err == nil {
		directory = path.Join(home, ".hectane")
	} else {
		log.Println(err)
		os.Exit(1)
	}
	flag.StringVar(&bind, "bind", ":8025", "address and port to bind to")
	flag.StringVar(&tlsCert, "tls-cert", "", "certificate for TLS")
	flag.StringVar(&tlsKey, "tls-key", "", "private key for TLS")
	flag.StringVar(&username, "username", "", "username for HTTP basic auth")
	flag.StringVar(&password, "password", "", "password for HTTP basic auth")
	flag.StringVar(&directory, "directory", directory, "directory for persistent storage")
	flag.BoolVar(&config.DisableSSLVerification, "disable-ssl-verification", false, "don't verify SSL certificates")
	flag.Parse()
	s := queue.NewStorage(directory)
	if q, err := queue.NewQueue(&config, s); err == nil {
		defer q.Stop()
		a := api.New(&api.Config{
			Addr:     bind,
			TLSCert:  tlsCert,
			TLSKey:   tlsKey,
			Username: username,
			Password: password,
		}, q, s)
		if err := a.Listen(); err != nil {
			log.Println(err)
			os.Exit(1)
		}
	} else {
		log.Println(err)
		os.Exit(1)
	}
}
