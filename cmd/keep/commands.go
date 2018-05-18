package keep

import (
	"github.com/keep-network/keep-core/cmd"
	"github.com/urfave/cli"
)

// KeepCommands is the set of actions that the Keep client application can perform
var KeepCommands = []cli.Command{
	{
		Name:   "smoke-test",
		Usage:  "smoke-test",
		Action: cmd.SmokeTestAction,
	},
}
