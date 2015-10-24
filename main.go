package main

import (
	"github.com/hectane/hectane/api"
	"github.com/hectane/hectane/exec"
	"github.com/hectane/hectane/queue"

	"flag"
	"log"
)

func main() {
	var (
		config   Config
		filename = flag.String("config", "", "file containing configuration")
	)
	config.RegisterFlags()
	flag.Parse()
	if *filename != "" {
		if err := config.LoadFromFile(*filename); err != nil {
			log.Fatal(err)
		}
	}
	q, err := queue.NewQueue(&config.Queue)
	if err != nil {
		log.Fatal(err)
	}
	defer q.Stop()
	a := api.New(&config.API, q)
	if err = a.Start(); err != nil {
		log.Fatal(err)
	}
	defer a.Stop()
	exec.Exec()
}
