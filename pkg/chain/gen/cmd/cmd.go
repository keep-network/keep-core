package cmd

import (
	"github.com/urfave/cli"
)

// AvailableCommands is the exported list of generated commands that can be
// installed on a CLI app. Generated contract command files set up init
// functions that add the contract's command and subcommands to this global
// variable, and any top-level command that wishes to include these commands can
// reference this variable and expect it to contain all generated contract
// commands.
var AvailableCommands []cli.Command
