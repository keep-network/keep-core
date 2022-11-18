package cmd

import (
	"context"
	"fmt"

	"github.com/keep-network/keep-core/config"
	"github.com/keep-network/keep-core/pkg/maintainer"
	"github.com/spf13/cobra"
)

// TbtcMaintainerCommand contains the definition of the tBTC maintainer
// command-line subcommand.
var TbtcMaintainerCommand = &cobra.Command{
	Use:   "maintainer",
	Short: `Starts tBTC maintainers`,
	Long:  `The maintainer command starts tBTC maintainers`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if err := clientConfig.ReadConfig(
			configFilePath,
			cmd.Flags(),
			config.Ethereum,
			config.BitcoinElectrum,
			config.TbtcMaintainer,
		); err != nil {
			logger.Fatalf("error reading config: %v", err)
		}
	},
	RunE: TbtcMaintainer,
}

func init() {
	initFlags(
		TbtcMaintainerCommand,
		&configFilePath,
		clientConfig,
		config.Ethereum,
		config.BitcoinElectrum,
		config.TbtcMaintainer,
	)
}

// TbtcMaintainer initializes maintainer tasks specified by flags passed to the
// maintainer command.
func TbtcMaintainer(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// TODO: Add connection to the Bitcoin electrum chain:
	// btcChain, err := electrum.Connect(ctx, &clientConfig.Electrum)
	// if err != nil {
	// 	return fmt.Errorf("could not connect to Electrum chain: [%v]", err)
	// }

	// TODO: Add connection to the Tbtc chain:
	// tbtcChain, err := newTbtcChain(config)
	// if err != nil {
	// 	return fmt.Errorf("could not connect to Tbtc chain: [%v]", err)
	// }

	maintainer.Initialize(ctx, clientConfig.Maintainer, nil, nil)

	<-ctx.Done()
	return fmt.Errorf("unexpected context cancellation")
}
