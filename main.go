package main

import (
	"github.com/hectane/hectane/api"
	"github.com/hectane/hectane/exec"
	"github.com/hectane/hectane/queue"

	"log"
)

func main() {
	c, err := exec.InitConfig()
	if err != nil {
		log.Fatal(err)
	}
	if err = exec.Init(c); err != nil {
		log.Fatal(err)
	}
	q, err := queue.NewQueue(c.Queue)
	if err != nil {
		log.Fatal(err)
	}
	defer q.Stop()
	a := api.New(c.API, q)
	if err = a.Start(); err != nil {
		log.Fatal(err)
	}
	defer a.Stop()
	exec.Exec(c)
}
