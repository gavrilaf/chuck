package cmds

import (
	"flag"
	. "github.com/gavrilaf/chuck/handlers"
	"github.com/gavrilaf/chuck/utils"
	"strings"
)

type IntgTestCommand struct {
	Log utils.Logger
}

func (c *IntgTestCommand) Run(args []string) int {
	flags := flag.NewFlagSet("intg", flag.ContinueOnError)
	cfg := NewScenarioSeekerConfig(flags, args, "intg")
	if cfg == nil {
		return 1
	}

	c.Log.Info("Running chuck in the integrations test mode")
	c.Log.Info("%s", cfg.String())

	proxy, err := CreateProxy()
	if err != nil {
		c.Log.Panic("Couldn't create a proxy, %v", err)
	}

	handler := NewScenarioSeekerHandler(cfg, c.Log)

	c.Log.Info("Running proxy...")
	err = RunProxy(proxy, handler, cfg.AddressAndPort())
	if err != nil {
		c.Log.Panic("Couldn't run a proxy, %v", err)
	}

	return 0
}

func (c *IntgTestCommand) Help() string {
	helpText := "Usage: chuck intg [addr:port] [folder]"
	return strings.TrimSpace(helpText)
}

func (c *IntgTestCommand) Synopsis() string {
	return "Run chuck in the integration tests mode"
}
