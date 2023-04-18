package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/keep-network/keep-core/config"
	"github.com/keep-network/keep-core/pkg/bitcoin/electrum"
	"github.com/keep-network/keep-core/pkg/maintainer"
)

// TbtcMaintainerCommand contains the definition of the tBTC maintainer
// command-line subcommand.
var TbtcMaintainerCommand = &cobra.Command{
	Use:   "maintainer",
	Short: `(experimental) Starts tBTC maintainers`,
	Long:  `(experimental) The maintainer command starts tBTC maintainers`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if err := clientConfig.ReadConfig(
			configFilePath,
			cmd.Flags(),
			config.MaintainerCategories...,
		); err != nil {
			logger.Fatalf("error reading config: %v", err)
		}
	},
	RunE: tbtcMaintainer,
}

func init() {
	initFlags(
		TbtcMaintainerCommand,
		&configFilePath,
		clientConfig,
		config.MaintainerCategories...,
	)
}

// tbtcMaintainer initializes maintainer tasks specified by flags passed to the
// maintainer command.
func tbtcMaintainer(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	bitcoinChain, err := electrum.Connect(ctx, clientConfig.Bitcoin.Electrum)
	if err != nil {
		return fmt.Errorf("could not connect to Electrum chain: [%v]", err)
	}

	// TODO: Add connection to the Tbtc chain:
	// tbtcChain, err := newTbtcChain(config)
	// if err != nil {
	// 	return fmt.Errorf("could not connect to Tbtc chain: [%v]", err)
	// }

	maintainer.Initialize(ctx, clientConfig.Maintainer, bitcoinChain, nil, nil)

	<-ctx.Done()
	return fmt.Errorf("unexpected context cancellation")
}
