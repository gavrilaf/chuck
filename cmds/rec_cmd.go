package cmds

import (
	. "chuck/handlers"
	"chuck/utils"
	"flag"
	"github.com/spf13/afero"
	"strings"
)

type RecordCommand struct {
	Log utils.Logger
	Fs  afero.Fs
}

func (self *RecordCommand) Run(args []string) int {
	flags := flag.NewFlagSet("rec", flag.ContinueOnError)
	cfg := NewRecorderConfig(flags, args, "log")
	if cfg == nil {
		return 1
	}

	self.Log.Info("Running chuck in the recording mode...")
	self.Log.Info("%s", cfg.String())

	proxy, err := CreateProxy()
	if err != nil {
		self.Log.Error("Couldn't create a proxy, %v", err)
		return 1
	}

	handler, err := NewRecorderHandler(cfg, self.Fs, self.Log)
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

func (self *RecordCommand) Help() string {
	helpText := "Usage: chuck rec [-address=addr] [-port=port] [-folder=folder] [-new_folder] [-prevent_304] [-focused]"
	return strings.TrimSpace(helpText)
}

func (self *RecordCommand) Synopsis() string {
	return "Run chuck in the recording mode"
}
