package cmds

import (
	"flag"
	. "github.com/gavrilaf/chuck/handlers"
	"github.com/gavrilaf/chuck/utils"
	"strings"
)

type DebugCommand struct {
	Log utils.Logger
}

func (c *DebugCommand) Run(args []string) int {
	flags := flag.NewFlagSet("dbg", flag.ContinueOnError)
	cfg := NewSeekerConfig(flags, args, "dbg")
	if cfg == nil {
		return 1
	}

	c.Log.Info("Running chuck in the debug mode")
	c.Log.Info("%s", cfg.String())

	proxy, err := CreateProxy()
	if err != nil {
		c.Log.Panic("Couldn't create a proxy, %v", err)
	}

	handler := NewSeekerHandler(cfg, c.Log)

	c.Log.Info("Running proxy...")
	err = RunProxy(proxy, handler, cfg.AddressAndPort())
	if err != nil {
		c.Log.Panic("Couldn't run a proxy, %v", err)
	}

	return 0
}

func (c *DebugCommand) Help() string {
	helpText := "Usage: chuck dbg [-address=addr] [-port=port] [-folder=folder]"
	return strings.TrimSpace(helpText)
}

func (c *DebugCommand) Synopsis() string {
	return "Run chuck in the debug mode"
}
