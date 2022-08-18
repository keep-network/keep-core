package cmd

import (
	"github.com/keep-network/keep-core/config"
	beaconcmd "github.com/keep-network/keep-core/pkg/chain/ethereum/beacon/gen/cmd"
	ecdsacmd "github.com/keep-network/keep-core/pkg/chain/ethereum/ecdsa/gen/cmd"
	tbtccmd "github.com/keep-network/keep-core/pkg/chain/ethereum/tbtc/gen/cmd"
	thresholdcmd "github.com/keep-network/keep-core/pkg/chain/ethereum/threshold/gen/cmd"
	"github.com/spf13/cobra"
)

// EthereumCommand contains the definition of the ethereum command-line
// subcommand and its own subcommands.
var EthereumCommand *cobra.Command

const ethereumDescription = `The ethereum command allows interacting with Keep's Ethereum
	contracts directly. Each subcommand corresponds to one contract, and has
	subcommands corresponding to each method on that contract, which respectively
	each take parameters based on the contract method's parameters.

    See the subcommand help for additional details.`

func init() {
	EthereumCommand = &cobra.Command{
		Use:   "ethereum",
		Short: `Provides access to Keep network Ethereum contracts.`,
		Long:  ethereumDescription,
	}

	initFlags(EthereumCommand, &configFilePath, clientConfig, config.General, config.Ethereum)

	EthereumCommand.AddCommand(beaconcmd.Command)
	EthereumCommand.AddCommand(ecdsacmd.Command)
	EthereumCommand.AddCommand(tbtccmd.Command)
	EthereumCommand.AddCommand(thresholdcmd.Command)
}
