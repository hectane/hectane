package main

import (
	"github.com/hectane/hectane/api"
	"github.com/hectane/hectane/queue"

	"flag"
	"log"
	"os"
)

func main() {
	var (
		filename string
		config   Config
		err      error
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
		a := api.New(&config.API, q)
		err = a.Listen()
	}
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
