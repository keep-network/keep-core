package cmd

import (
	chaincmd "github.com/keep-network/keep-core/pkg/chain/gen/cmd"
	"github.com/urfave/cli"
)

// ContractCommand contains the definition of the contract command-line
// subcommand and its own subcommands.
var ContractCommand cli.Command

const contractDescription = `   The contract command allows interacting with Keep's contracts directly. Each
    subcommand corresponds to one contract, and has subcommands corresponding to
    each method on that contract, which respectively each take parameters based
    on the contract method's parameters.

    See the subcommand help for additional details.`

func init() {
	ContractCommand = cli.Command{
		Name:        "contract",
		Usage:       `Provides access to Keep network contracts.`,
		Description: contractDescription,
		Subcommands: chaincmd.AvailableCommands,
	}
}
