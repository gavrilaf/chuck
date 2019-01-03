package cmds

import (
	"flag"
	. "github.com/gavrilaf/chuck/handlers"
	"github.com/gavrilaf/chuck/utils"
	"github.com/spf13/afero"
	"strings"
)

type IntgTestCommand struct {
	Log utils.Logger
	Fs  afero.Fs
}

func (self *IntgTestCommand) Run(args []string) int {
	flags := flag.NewFlagSet("intg", flag.ContinueOnError)
	cfg := NewScenarioSeekerConfig(flags, args, "intg")
	if cfg == nil {
		return 1
	}

	self.Log.Info("Running chuck in the integrations test mode")
	self.Log.Info("%s", cfg.String())

	proxy, err := CreateProxy()
	if err != nil {
		self.Log.Error("Couldn't create a proxy, %v", err)
		return 1
	}

	handler, err := NewScenarioSeekerHandler(cfg, self.Fs, self.Log)
	if err != nil {
		self.Log.Error("Couldn't create a handler, %v", err)
		return 1
	}

	self.Log.Info("Running proxy...")
	err = RunProxy(proxy, handler, cfg.AddressAndPort())
	if err != nil {
		self.Log.Error("Couldn't run a proxy, %v", err)
	}

	return 0
}

func (self *IntgTestCommand) Help() string {
	helpText := "Usage: chuck intg [addr:port] [folder]"
	return strings.TrimSpace(helpText)
}

func (self *IntgTestCommand) Synopsis() string {
	return "Run chuck in the integration tests mode"
}
