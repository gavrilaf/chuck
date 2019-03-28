package main

import (
	. "chuck/cmds"
	"chuck/utils"

	"github.com/mitchellh/cli"
	"github.com/spf13/afero"
	"os"
)

const AppName = "chuck"
const Version = "0.1.0"

func main() {
	ui := &cli.ColoredUi{
		InfoColor:   cli.UiColorGreen,
		WarnColor:   cli.UiColorYellow,
		ErrorColor:  cli.UiColorRed,
		OutputColor: cli.UiColorMagenta,
		Ui: &cli.BasicUi{
			Writer:      os.Stdout,
			ErrorWriter: os.Stderr,
		},
	}

	log := utils.NewLogger(ui)
	fs := afero.NewOsFs()

	c := cli.NewCLI(AppName, Version)
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"rec": func() (cli.Command, error) {
			return &RecordCommand{
				Log: log,
				Fs:  fs,
			}, nil
		},
		"dbg": func() (cli.Command, error) {
			return &DebugCommand{
				Log: log,
				Fs:  fs,
			}, nil
		},
		"intg": func() (cli.Command, error) {
			return &IntgTestCommand{
				Log: log,
				Fs:  fs,
			}, nil
		},
		"intg_rec": func() (cli.Command, error) {
			return &IntgTestRecCommand{
				Log: log,
				Fs:  fs,
			}, nil
		},
	}

	_, err := c.Run()
	if err != nil {
		log.Panic("Failed with error: %v", err)
	}
}
