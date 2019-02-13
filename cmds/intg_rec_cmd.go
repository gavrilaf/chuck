package cmds

import (
	. "chuck/handlers"
	"chuck/utils"
	"flag"
	"github.com/spf13/afero"
	"strings"
)

type IntgTestRecCommand struct {
	Log utils.Logger
	Fs  afero.Fs
}

func (self *IntgTestRecCommand) Run(args []string) int {
	flags := flag.NewFlagSet("intg_rec", flag.ContinueOnError)
	cfg := NewScenarioRecorderConfig(flags, args, "intg_rec")
	if cfg == nil {
		return 1
	}

	self.Log.Info("Running chuck in the integrations test recording mode")
	self.Log.Info("%s", cfg.String())

	proxy, err := CreateProxy()
	if err != nil {
		self.Log.Error("Couldn't create a proxy, %v", err)
	}

	handler, err := NewScenarioRecorderHandler(cfg, self.Fs, self.Log)
	if err != nil {
		self.Log.Error("Couldn't create a proxy, %v", err)
		return 1
	}

	self.Log.Info("Running proxy...")
	err = RunProxy(proxy, handler, cfg.AddressAndPort())
	if err != nil {
		self.Log.Error("Couldn't run a proxy, %v", err)
	}

	return 0
}

func (self *IntgTestRecCommand) Synopsis() string {
	return "Run chuck in the integration tests recording mode"
}

func (self *IntgTestRecCommand) Help() string {
	helpText := "Usage: chuck intg_rec [-address=addr] [-port=port] [-folder=folder] [-new_folder] [-prevent_304]"
	return strings.TrimSpace(helpText)
}
