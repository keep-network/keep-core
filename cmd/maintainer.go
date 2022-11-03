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
			config.Electrum,
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
		config.Electrum,
		config.Maintainer,
	)
}

// maintainers initializes maintainer tasks specified by flags passed to the
// maintainer command.
func maintainers(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	err := maintainer.Initialize(clientConfig.Maintainer, ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize maintainers [%v]", err)
	}

	<-ctx.Done()
	return fmt.Errorf("unexpected context cancellation")
}
