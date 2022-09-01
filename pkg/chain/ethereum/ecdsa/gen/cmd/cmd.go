package cmd

import (
	"github.com/spf13/cobra"

	"github.com/keep-network/keep-core/pkg/chain/ethereum"
)

// ModuleCommand is the exported list of generated commands that can be
// installed on a CLI app. Generated contract command files set up init
// functions that add the contract's command and subcommands to this global
// variable, and any top-level command that wishes to include these commands can
// reference this variable and expect it to contain all generated contract
// commands.
var ModuleCommand = ethereum.Command{
	Command: cobra.Command{
		Use:   "ecdsa",
		Short: `Provides access to Keep ECDSA contracts.`,
	},
}
