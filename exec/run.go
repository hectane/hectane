package exec

import (
	"github.com/hectane/hectane/api"
	"github.com/hectane/hectane/queue"

	"os"
	"os/signal"
	"syscall"
)

// Run the application until the SIGINT signal is received.
var RunCommand = &Command{
	Name:        "run",
	Description: "start the application",
	Exec: func(cfg *Config) error {
		q, err := queue.NewQueue(cfg.Queue)
		if err != nil {
			return err
		}
		defer q.Stop()
		a := api.New(cfg.API, q)
		if err = a.Start(); err != nil {
			return err
		}
		defer a.Stop()
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		<-c
		return nil
	},
}
