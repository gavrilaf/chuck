package cmds

import (
	. "chuck/handlers"
	"chuck/utils"
	"flag"
	"strings"

	"github.com/spf13/afero"
)

type IntgNoProxyTestCommand struct {
	Log utils.Logger
	Fs  afero.Fs
}

func (self *IntgNoProxyTestCommand) Run(args []string) int {
	flags := flag.NewFlagSet("intg-noproxy", flag.ContinueOnError)
	cfg := NewScenarioSeekerConfig(flags, args, "intg")
	if cfg == nil {
		return 1
	}

	self.Log.Info("Running chuck in the integrations test mode (no-proxy mode)")
	self.Log.Info("%s", cfg.String())

	handler, err := NewScenarioSeekerNoProxyHandler(cfg, self.Fs, self.Log)
	if err != nil {
		self.Log.Error("Couldn't create a handler, %v", err)
		return 1
	}

	self.Log.Info("Running stub server...")
	err = RunServer(handler, cfg.AddressAndPort())
	if err != nil {
		self.Log.Error("Couldn't run a stub server, %v", err)
	}

	return 0
}

func (self *IntgNoProxyTestCommand) Help() string {
	helpText := "Usage: chuck intg-noproxy [-address=addr] [-port=port] [-folder=folder] [-verbose]"
	return strings.TrimSpace(helpText)
}

func (self *IntgNoProxyTestCommand) Synopsis() string {
	return "Run chuck in the integration tests mode"
}
