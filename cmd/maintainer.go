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
	Short: `Starts keep maintainers`,
	Long:  `The maintainer command starts keep maintainers`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if err := clientConfig.ReadConfig(
			configFilePath,
			cmd.Flags(),
			config.Ethereum,
			config.Electrs,
			config.Maintainer,
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
		config.Ethereum,
		config.Electrs,
		config.Maintainer,
	)
}

// maintainers initializes maintainer tasks specified by flags passed to the
// maintainer command.
func maintainers(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// TODO: Add connection to the Electrs chain:
	// btcChain, err := electrs.Connect(ctx, &clientConfig.Electrs)
	// if err != nil {
	// 	return fmt.Errorf("could not connect to Electrs chain: [%v]", err)
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
