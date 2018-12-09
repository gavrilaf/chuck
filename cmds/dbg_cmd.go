package cmds

import (
	. "github.com/gavrilaf/chuck/handlers"
	"github.com/gavrilaf/chuck/utils"
	"strings"
)

type DebugCommand struct {
	Log utils.Logger
}

func (c *DebugCommand) Help() string {
	helpText := `
Usage: chuck dbg [addr:port] [folder]
`
	return strings.TrimSpace(helpText)
}

func (c *DebugCommand) Run(args []string) int {
	c.Log.Info("Running chuck in the debug mode")

	addr := ":8123"
	folder := "dbg"

	proxy, err := CreateProxy()
	if err != nil {
		c.Log.Panic("Couldn't create a proxy, %v", err)
	}

	handler := NewSeekerHandler(folder, c.Log)

	c.Log.Info("Running proxy...")
	err = RunProxy(proxy, handler, addr)
	if err != nil {
		c.Log.Panic("Couldn't run a proxy, %v", err)
	}

	return 0
}

func (c *DebugCommand) Synopsis() string {
	return "Run chuck in the debug mode"
}
