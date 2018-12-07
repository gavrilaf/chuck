package main

import (
	"github.com/gavrilaf/chuck/utils"
	"strings"
)

type DebugCommand struct {
	log utils.Logger
}

func (c *DebugCommand) Help() string {
	helpText := `
Usage: chuck dbg [addr:port] [folder]
`
	return strings.TrimSpace(helpText)
}

func (c *DebugCommand) Run(args []string) int {
	c.log.Info("Running chuck in the debug mode")

	addr := ":8123"
	folder := "dbg"

	proxy, err := CreateProxy()
	if err != nil {
		c.log.Panic("Couldn't create a proxy, %v", err)
	}

	handler := NewSeekerHandler(folder, c.log)

	c.log.Info("Running proxy...")
	err = RunProxy(proxy, handler, addr)
	if err != nil {
		c.log.Panic("Couldn't run a proxy, %v", err)
	}

	return 0
}

func (c *DebugCommand) Synopsis() string {
	return "Run chuck in the debug mode"
}
