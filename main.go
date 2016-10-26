package main

import (
	"github.com/hectane/hectane/api"
	"github.com/hectane/hectane/cfg"
	"github.com/hectane/hectane/cmd"
	"github.com/hectane/hectane/exec"
	"github.com/hectane/hectane/log"
	"github.com/hectane/hectane/queue"
	"github.com/hectane/hectane/smtp"

	"errors"
	"flag"
	"fmt"
	"os"
)

// Display usage information for the application.
func printUsage() {
	fmt.Fprintf(os.Stderr, "USAGE\n\thectane [flags] [command]\n\n")
	fmt.Fprintf(os.Stderr, "COMMANDS\n")
	cmd.Print()
	fmt.Fprintf(os.Stderr, "\nFLAGS\n")
	flag.PrintDefaults()
}

// Run the application.
func runApplication(config *cfg.Config) error {
	if err := log.Init(&config.Log); err != nil {
		return err
	}
	defer log.Cleanup()
	q, err := queue.NewQueue(&config.Queue)
	if err != nil {
		return err
	}
	defer q.Stop()
	a := api.New(&config.API, q)
	if err = a.Start(); err != nil {
		return err
	}
	defer a.Stop()
	s, err := smtp.New(&config.SMTP, q)
	if err != nil {
		return err
	}
	defer s.Close()
	if err = exec.Exec(); err != nil {
		return err
	}
	return nil
}

// This needs to be a separate function from main in order to ensure that the
// deferred statements are run while still allowing os.Exit() to be given an
// error code. If a single argument was specified, it's a command. Otherwise,
// run the application using the current platform's execution environment.
func run() error {
	cmd.Init()
	flag.Usage = printUsage
	config, err := cfg.Parse()
	if err != nil {
		return err
	}
	switch {
	case flag.NArg() == 0:
		return runApplication(config)
	case flag.NArg() == 1:
		return cmd.Exec(flag.Args()[0], config)
	default:
		return errors.New("single command expected")
	}
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err)
		os.Exit(1)
	}
}
