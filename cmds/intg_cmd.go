package cmds

import (
	. "github.com/gavrilaf/chuck/handlers"
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

	addr := ":8123"
	folder := "aadhi"

	proxy, err := CreateProxy()
	if err != nil {
		c.Log.Panic("Couldn't create a proxy, %v", err)
	}

	handler := NewScenarioHandler(folder, c.Log)

	c.Log.Info("Running proxy...")
	err = RunProxy(proxy, handler, addr)
	if err != nil {
		c.Log.Panic("Couldn't run a proxy, %v", err)
	}

	return 0
}

func (c *IntgTestCommand) Synopsis() string {
	return "Run chuck in the integration tests mode"
}
