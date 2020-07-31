package main

import (
	"log"
	"os"

	"github.com/mitchellh/cli"
	"github.com/yasuoza/switchbot/cmd/switchbot/command"
)

var Version = "current"

func main() {
	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stdout,
	}

	c := cli.NewCLI("switchbot", Version)
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"scan": func() (cli.Command, error) {
			return &command.ScanCommand{UI: ui}, nil
		},
		"press": func() (cli.Command, error) {
			return &command.PressCommand{UI: ui}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
