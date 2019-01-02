package cmds

import (
	"flag"
	. "github.com/gavrilaf/chuck/handlers"
	"github.com/gavrilaf/chuck/utils"
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
		self.Log.Panic("Couldn't create a proxy, %v", err)
	}

	handler := NewRecorderHandler(cfg, self.Fs, self.Log)

	self.Log.Info("Running proxy...")
	err = RunProxy(proxy, handler, cfg.AddressAndPort())
	if err != nil {
		self.Log.Panic("Couldn't run a proxy, %v", err)
	}

	return 0
}

func (self *RecordCommand) Help() string {
	helpText := "Usage: chuck rec [-address=addr] [-port=port] [-folder=folder] [-new_folder] [-prevent_304] [-focused] [-print_only]"
	return strings.TrimSpace(helpText)
}

func (self *RecordCommand) Synopsis() string {
	return "Run chuck in the recording mode"
}
