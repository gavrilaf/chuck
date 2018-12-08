package main

import (
	"github.com/gavrilaf/chuck/utils"
	"github.com/mitchellh/cli"
	"os"
)

const AppName = "chuck"
const Version = "0.0.1"

func main() {
	ui := &cli.ColoredUi{
		InfoColor:  cli.UiColorGreen,
		WarnColor:  cli.UiColorYellow,
		ErrorColor: cli.UiColorRed,
		Ui: &cli.BasicUi{
			Writer:      os.Stdout,
			ErrorWriter: os.Stderr,
		},
	}
	log := utils.NewLogger(ui)

	c := cli.NewCLI(AppName, Version)
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"rec": func() (cli.Command, error) {
			return &RecordCommand{
				log: log,
			}, nil
		},
		"dbg": func() (cli.Command, error) {
			return &DebugCommand{
				log: log,
			}, nil
		},
		"intg": func() (cli.Command, error) {
			return &IntgTestCommand{
				log: log,
			}, nil
		},
	}

	_, err := c.Run()
	if err != nil {
		log.Panic("Failed with error: %v", err)
	}
}
