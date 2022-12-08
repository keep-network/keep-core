package cmd

import (
	"context"
	"fmt"

	"github.com/keep-network/keep-core/config"
	"github.com/keep-network/keep-core/pkg/maintainer"
	"github.com/spf13/cobra"
)

// MaintainerCommand contains the definition of the maintainer command-line
// subcommand.
var MaintainerCommand = &cobra.Command{
	Use:   "maintainer",
	Short: `(experimental) Starts maintainers`,
	Long:  `(experimental) The maintainer command starts maintainers`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if err := clientConfig.ReadConfig(
			configFilePath,
			cmd.Flags(),
			config.MaintainerCategories...,
		); err != nil {
			logger.Fatalf("error reading config: %v", err)
		}
	},
	RunE: maintainers,
}

func init() {
	initFlags(
		MaintainerCommand,
		&configFilePath,
		clientConfig,
		config.MaintainerCategories...,
	)
}

// maintainers initializes maintainer tasks specified by flags passed to the
// maintainer command.
func maintainers(cmd *cobra.Command, args []string) error {
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
