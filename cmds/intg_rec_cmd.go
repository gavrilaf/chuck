package cmds

import (
	"flag"
	. "github.com/gavrilaf/chuck/handlers"
	"github.com/gavrilaf/chuck/utils"
	"strings"
)

type IntgTestRecCommand struct {
	Log utils.Logger
}

func (c *IntgTestRecCommand) Run(args []string) int {
	flags := flag.NewFlagSet("intg_rec", flag.ContinueOnError)
	cfg := NewScenarioRecorderConfig(flags, args, "intg_rec")
	if cfg == nil {
		return 1
	}

	c.Log.Info("Running chuck in the integrations test recording mode")
	c.Log.Info("%s", cfg.String())

	proxy, err := CreateProxy()
	if err != nil {
		c.Log.Panic("Couldn't create a proxy, %v", err)
	}

	handler := NewScenarioRecorderHandler(cfg, c.Log)

	c.Log.Info("Running proxy...")
	err = RunProxy(proxy, handler, cfg.AddressAndPort())
	if err != nil {
		c.Log.Panic("Couldn't run a proxy, %v", err)
	}

	return 0
}

func (c *IntgTestRecCommand) Synopsis() string {
	return "Run chuck in the integration tests recording mode"
}

func (c *IntgTestRecCommand) Help() string {
	helpText := "Usage: chuck intg_rec [-address=addr] [-port=port] [-folder=folder] [-new_folder] [-prevent_304]"
	return strings.TrimSpace(helpText)
}
