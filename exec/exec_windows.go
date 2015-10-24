package exec

import (
	"flag"
)

var service = flag.Bool("service", false, "run as a Windows service")

// Run the application either using signals or as a Windows service.
func Exec() {
	if *service {
		go execService()
	} else {
		go execSignal()
	}
}
