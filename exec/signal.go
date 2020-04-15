package exec

import (
	"os"
	"os/signal"
	"syscall"
)

// Run until SIGINT is received.
func execSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
}
