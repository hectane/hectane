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
		filename = flag.String("config", "", "file containing configuration")
		config   Config
	)
	config.RegisterFlags()
	flag.Parse()
	if *filename != "" {
		if err := config.LoadFromFile(*filename); err != nil {
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
		if err = a.Start(); err == nil {
			defer a.Stop()
		} else {
			log.Println(err)
			os.Exit(1)
		}
		signal.Notify(c, syscall.SIGINT)
		<-c
	} else {
		log.Println(err)
		os.Exit(1)
	}
}
