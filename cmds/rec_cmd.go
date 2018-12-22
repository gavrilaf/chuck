package cmds

import (
	"flag"
	. "github.com/gavrilaf/chuck/handlers"
	"github.com/gavrilaf/chuck/utils"
	"strings"
)

type RecordCommand struct {
	Log utils.Logger
}

func (c *RecordCommand) Run(args []string) int {
	flags := flag.NewFlagSet("rec", flag.ContinueOnError)
	cfg := NewRecorderConfig(flags, args, "log")
	if cfg == nil {
		return 1
	}

	c.Log.Info("Running chuck in the recording mode...")
	c.Log.Info("%s", cfg.String())

	proxy, err := CreateProxy()
	if err != nil {
		c.Log.Panic("Couldn't create a proxy, %v", err)
	}

	handler := NewRecorderHandler(cfg, c.Log)

	c.Log.Info("Running proxy...")
	err = RunProxy(proxy, handler, cfg.AddressAndPort())
	if err != nil {
		c.Log.Panic("Couldn't run a proxy, %v", err)
	}

	return 0
}

func (c *RecordCommand) Help() string {
	helpText := "Usage: chuck rec [-address=addr] [-port=port] [-folder=folder] [-new_folder] [-prevent_304] [-focused] [-print_only]"
	return strings.TrimSpace(helpText)
}

func (c *RecordCommand) Synopsis() string {
	return "Run chuck in the recording mode"
}
