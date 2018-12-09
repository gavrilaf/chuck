package cmds

import (
	"github.com/gavrilaf/chuck/utils"
	"strings"
)

type IntgTestCommand struct {
	Log utils.Logger
}

func (c *IntgTestCommand) Help() string {
	helpText := `
Usage: chuck intg [addr:port] [folder]
`
	return strings.TrimSpace(helpText)
}

func (c *IntgTestCommand) Run(args []string) int {
	c.Log.Info("Running chuck in the integrations test mode")
	c.Log.Error("Doesn't supported yet")

	return 0
}

func (c *IntgTestCommand) Synopsis() string {
	return "Run chuck in the integration tests mode"
}
