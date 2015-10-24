package exec

import (
	"os"
	"os/signal"
	"syscall"
)

// Run the application until SIGINT is received.
func execSignal() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT)
	<-c
}
