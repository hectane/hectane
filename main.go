package main

import (
	"github.com/hectane/hectane/api"
	"github.com/hectane/hectane/queue"

	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var (
		filename string
		config   Config
	)
	config.RegisterFlags()
	flag.StringVar(&filename, "config", "", "file containing configuration")
	flag.Parse()
	if filename != "" {
		if err := config.LoadFromFile(filename); err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}
	if q, err := queue.NewQueue(&config.Queue); err == nil {
		defer q.Stop()
		var (
			a = api.New(&config.API, q)
			c = make(chan os.Signal)
		)
		signal.Notify(c, syscall.SIGINT)
		go func() {
			<-c
			a.Stop()
		}()
		if err = a.Listen(); err != nil {
			log.Println(err)
			os.Exit(1)
		}
	} else {
		log.Println(err)
		os.Exit(1)
	}
}
