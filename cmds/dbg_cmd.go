package cmds

import (
	"flag"
	. "github.com/gavrilaf/chuck/handlers"
	"github.com/gavrilaf/chuck/utils"
	"github.com/spf13/afero"
	"strings"
)

type DebugCommand struct {
	Log utils.Logger
	Fs  afero.Fs
}

func (self *DebugCommand) Run(args []string) int {
	flags := flag.NewFlagSet("dbg", flag.ContinueOnError)
	cfg := NewSeekerConfig(flags, args, "dbg")
	if cfg == nil {
		return 1
	}

	self.Log.Info("Running chuck in the debug mode")
	self.Log.Info("%s", cfg.String())

	proxy, err := CreateProxy()
	if err != nil {
		self.Log.Error("Couldn't create a proxy, %v", err)
		return 1
	}

	handler, err := NewSeekerHandler(cfg, self.Fs, self.Log)
	if err != nil {
		self.Log.Error("Couldn't create a handler, %v", err)
		return 1
	}

	self.Log.Info("Running proxy...")
	err = RunProxy(proxy, handler, cfg.AddressAndPort())
	if err != nil {
		self.Log.Error("Couldn't run a proxy, %v", err)
		return 1
	}

	return 0
}

func (self *DebugCommand) Help() string {
	helpText := "Usage: chuck dbg [-address=addr] [-port=port] [-folder=folder]"
	return strings.TrimSpace(helpText)
}

func (self *DebugCommand) Synopsis() string {
	return "Run chuck in the debug mode"
}
