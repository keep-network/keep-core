package cmd

import (
	beaconcmd "github.com/keep-network/keep-core/pkg/chain/ethereum/beacon/gen/cmd"
	ecdsacmd "github.com/keep-network/keep-core/pkg/chain/ethereum/ecdsa/gen/cmd"
	tbtccmd "github.com/keep-network/keep-core/pkg/chain/ethereum/tbtc/gen/cmd"
	thresholdcmd "github.com/keep-network/keep-core/pkg/chain/ethereum/threshold/gen/cmd"
	"github.com/urfave/cli"
)

// EthereumCommand contains the definition of the ethereum command-line
// subcommand and its own subcommands.
var EthereumCommand cli.Command

const ethereumDescription = `The ethereum command allows interacting with Keep's Ethereum
	contracts directly. Each subcommand corresponds to one contract, and has
	subcommands corresponding to each method on that contract, which respectively
	each take parameters based on the contract method's parameters.

    See the subcommand help for additional details.`

func init() {
	var subcommands []cli.Command
	subcommands = append(subcommands, beaconcmd.AvailableCommands...)
	subcommands = append(subcommands, ecdsacmd.AvailableCommands...)
	subcommands = append(subcommands, tbtccmd.AvailableCommands...)
	subcommands = append(subcommands, thresholdcmd.AvailableCommands...)

	EthereumCommand = cli.Command{
		Name:        "ethereum",
		Usage:       `Provides access to Keep network Ethereum contracts.`,
		Description: ethereumDescription,
		Subcommands: subcommands,
	}
}
