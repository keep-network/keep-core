package cmd

import (
	"github.com/spf13/cobra"

	"github.com/keep-network/keep-core/config"

	beaconcmd "github.com/keep-network/keep-core/pkg/chain/ethereum/beacon/gen/cmd"
	ecdsacmd "github.com/keep-network/keep-core/pkg/chain/ethereum/ecdsa/gen/cmd"
	tbtccmd "github.com/keep-network/keep-core/pkg/chain/ethereum/tbtc/gen/cmd"
	thresholdcmd "github.com/keep-network/keep-core/pkg/chain/ethereum/threshold/gen/cmd"
)

// EthereumCommand contains the definition of the ethereum command-line
// subcommand and its own subcommands.
var EthereumCommand *cobra.Command

const ethereumDescription = `The ethereum command allows interacting with Keep's Ethereum
contracts directly. Each subcommand corresponds to one contract, and has subcommands 
corresponding to each method on that contract, which respectively each take parameters 
based on the contract method's parameters.

See the subcommand help for additional details.`

func init() {
	EthereumCommand = &cobra.Command{
		Use:   "ethereum",
		Short: `Provides access to Keep network Ethereum contracts.`,
		Long:  ethereumDescription,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if err := clientConfig.ReadConfig(configFilePath, cmd.Flags(), config.General, config.Ethereum); err != nil {
				logger.Fatalf("error reading config: %v", err)
			}

			beaconcmd.ModuleCommand.SetConfig(&clientConfig.Ethereum)
			ecdsacmd.ModuleCommand.SetConfig(&clientConfig.Ethereum)
			tbtccmd.ModuleCommand.SetConfig(&clientConfig.Ethereum)
			thresholdcmd.ModuleCommand.SetConfig(&clientConfig.Ethereum)
		},
	}

	initFlags(EthereumCommand, &configFilePath, clientConfig, config.General, config.Ethereum)

	EthereumCommand.AddCommand(&beaconcmd.ModuleCommand.Command)
	EthereumCommand.AddCommand(&ecdsacmd.ModuleCommand.Command)
	EthereumCommand.AddCommand(&tbtccmd.ModuleCommand.Command)
	EthereumCommand.AddCommand(&thresholdcmd.ModuleCommand.Command)
}
