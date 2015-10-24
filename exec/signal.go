package exec

import (
	"os"
	"os/signal"
	"syscall"
)

// Run the application until a signal is received.
func execSignal(stop chan<- bool) {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT)
	<-c
	close(stop)
}
