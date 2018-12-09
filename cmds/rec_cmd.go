package cmds

import (
	. "github.com/gavrilaf/chuck/handlers"
	"github.com/gavrilaf/chuck/utils"
	"strings"
)

type RecordCommand struct {
	Log utils.Logger
}

func (c *RecordCommand) Help() string {
	helpText := `
Usage: chuck rec [addr:port] [folder] [--force] 
`
	return strings.TrimSpace(helpText)
}

func (c *RecordCommand) Run(args []string) int {
	c.Log.Info("Running chuck in the record mode")

	addr := ":8123"
	folder := "log"

	proxy, err := CreateProxy()
	if err != nil {
		c.Log.Panic("Couldn't create a proxy, %v", err)
	}

	handler := NewRecordHandler(folder, c.Log)

	c.Log.Info("Running proxy...")
	err = RunProxy(proxy, handler, addr)
	if err != nil {
		c.Log.Panic("Couldn't run a proxy, %v", err)
	}

	return 0
}

func (c *RecordCommand) Synopsis() string {
	return "Run chuck in the recording mode"
}
