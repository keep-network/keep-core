package cmd

import (
	"github.com/keep-network/keep-core/cmd/flag"
	"github.com/urfave/cli"
)

const (
	blockFlag        string = "block"
	blockShort       string = "b"
	transactionFlag  string = "transaction"
	transactionShort string = "t"
	submitFlag       string = "submit"
	submitShort      string = "s"
	valueFlag        string = "value"
	valueShort       string = "v"
)

var (
	transactionFlagValue *flag.TransactionHash = &flag.TransactionHash{}
	blockFlagValue       *flag.Uint256         = &flag.Uint256{}
	valueFlagValue       *flag.Uint256         = &flag.Uint256{}
)

// AvailableCommands is the exported list of generated commands that can be
// installed on a CLI app. Generated contract command files set up init
// functions that add the contract's command and subcommands to this global
// variable, and any top-level command that wishes to include these commands can
// reference this variable and expect it to contain all generated contract
// commands.
var AvailableCommands []cli.Command
