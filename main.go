package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/hectane/hectane/api"
	"github.com/hectane/hectane/cfg"
	"github.com/hectane/hectane/cmd"
	"github.com/hectane/hectane/exec"
	"github.com/hectane/hectane/queue"

	"flag"
	"fmt"
	"os"
)

// Display usage information for the application.
func printUsage() {
	fmt.Fprintf(os.Stderr, "USAGE\n\thectane [options] [command]\n\n")
	fmt.Fprintf(os.Stderr, "COMMANDS\n")
	cmd.Print()
	fmt.Fprintf(os.Stderr, "\nOPTIONS\n")
	flag.PrintDefaults()
}

// Log the specified message and return a suitable error code.
func logError(err error) int {
	logrus.Error(err)
	return 1
}

// This needs to be a separate function from main in order to ensure that the
// deferred statements are run while still allowing os.Exit() to be given an
// error code.
func run() int {

	// Initialize the execution environment
	if err := exec.Init(); err != nil {
		logrus.Fatal(err)
	}
	defer exec.Cleanup()

	// Initialize the list of commands
	cmd.Init()

	// Set up the usage function for flags and parse them
	flag.Usage = printUsage
	config, err := cfg.Parse()
	if err != nil {
		return logError(err)
	}

	// If a single argument was specified, it's a command - otherwise, run the
	// application using the current platform's execution environment
	switch {
	case flag.NArg() == 0:
		q, err := queue.NewQueue(&config.Queue)
		if err != nil {
			return logError(err)
		}
		defer q.Stop()
		a := api.New(&config.API, q)
		if err = a.Start(); err != nil {
			return logError(err)
		}
		defer a.Stop()
		if err = exec.Exec(); err != nil {
			return logError(err)
		}
	case flag.NArg() == 1:
		if err := cmd.Exec(flag.Args()[0], config); err != nil {
			return logError(err)
		}
	default:
		logrus.Error("single command expected")
		return 1
	}

	// If execution reaches this point, there were no errors
	return 0
}

func main() {
	os.Exit(run())
}
