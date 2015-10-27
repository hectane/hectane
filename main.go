package main

import (
	"github.com/hectane/hectane/api"
	"github.com/hectane/hectane/cfg"
	"github.com/hectane/hectane/cmd"
	"github.com/hectane/hectane/exec"
	"github.com/hectane/hectane/queue"

	"flag"
	"fmt"
	"log"
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

func main() {
	cmd.Init()
	flag.Usage = printUsage
	config, err := cfg.Parse()
	if err != nil {
		log.Fatal(err)
	}
	switch {
	case flag.NArg() == 0:
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
		if err = exec.Exec(); err != nil {
			log.Fatal(err)
		}
	case flag.NArg() == 1:
		if err := cmd.Exec(flag.Args()[0]); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("single command expected")
	}
}
