package cmds

import (
	. "github.com/gavrilaf/chuck/handlers"
	"github.com/gavrilaf/chuck/utils"
	"strings"
)

type IntgTestRecCommand struct {
	Log utils.Logger
}

func (c *IntgTestRecCommand) Help() string {
	helpText := `
Usage: chuck intg_rec [addr:port] [folder]
`
	return strings.TrimSpace(helpText)
}

func (c *IntgTestRecCommand) Run(args []string) int {
	c.Log.Info("Running chuck in the integrations test recording mode")

	addr := ":8123"
	folder := "aadhi_rec"

	proxy, err := CreateProxy()
	if err != nil {
		c.Log.Panic("Couldn't create a proxy, %v", err)
	}

	handler := NewScenarioRecorderHandler(folder, c.Log)

	c.Log.Info("Running proxy...")
	err = RunProxy(proxy, handler, addr)
	if err != nil {
		c.Log.Panic("Couldn't run a proxy, %v", err)
	}

	return 0
}

func (c *IntgTestRecCommand) Synopsis() string {
	return "Run chuck in the integration tests recording mode"
}
