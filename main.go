package main

import (
	"github.com/hectane/hectane/exec"

	"log"
)

func main() {
	exec.InitCommands()
	name, cfg, err := exec.InitConfig()
	if err != nil {
		log.Fatal(err)
	}
	cmd, ok := exec.Commands[name]
	if !ok {
		log.Fatalf("unrecognized command \"%s\"", name)
	}
	if err := cmd.Exec(cfg); err != nil {
		log.Fatal(err)
	}
}
