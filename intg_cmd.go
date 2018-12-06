package main

import (
	"github.com/gavrilaf/chuck/utils"
	"strings"
)

type IntgTestCommand struct {
	log utils.Logger
}

func (c *IntgTestCommand) Help() string {
	helpText := `
Usage: chuck intg [addr:port] [folder]
`
	return strings.TrimSpace(helpText)
}

func (c *IntgTestCommand) Run(args []string) int {
	c.log.Info("Running chuck in the integrations test mode")
	c.log.Error("Doesn't supported yet")

	return 0
}

func (c *IntgTestCommand) Synopsis() string {
	return "Run chuck in the integration tests mode"
}
