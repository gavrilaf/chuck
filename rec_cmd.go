package main

import (
	"github.com/gavrilaf/chuck/utils"
	"strings"
)

type RecordCommand struct {
	log utils.Logger
}

func (c *RecordCommand) Help() string {
	helpText := `
Usage: chuck rec [addr:port] [folder] [--force] 
`
	return strings.TrimSpace(helpText)
}

func (c *RecordCommand) Run(args []string) int {
	c.log.Info("Running chuck in the record mode")

	addr := ":8123"

	proxy, err := CreateProxy()
	if err != nil {
		c.log.Panic("Couldn't create a proxy, %v", err)
	}

	handler := NewRecordHandler(c.log)

	c.log.Info("Running proxy...")
	err = RunProxy(proxy, handler, addr)
	if err != nil {
		c.log.Panic("Couldn't run a proxy, %v", err)
	}

	return 0
}

func (c *RecordCommand) Synopsis() string {
	return "Run chuck in the recording mode"
}
