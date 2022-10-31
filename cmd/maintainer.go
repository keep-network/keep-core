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
			config.Bitcoin,
			config.Maintainer,
		); err != nil {
			logger.Fatalf("error reading config: %v", err)
		}
	},
	RunE: launchMaintainers,
}

func init() {
	initFlags(
		MaintainerCommand,
		&configFilePath,
		clientConfig,
		config.Ethereum,
		config.Bitcoin,
		config.Maintainer,
	)
}

// launchMaintainers launches maintainer tasks specified by flags passed to
// the maintainer command. If there were no tasks specified, all the maintainer
// tasks are launched.
func launchMaintainers(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	launchRelay, err := cmd.Flags().GetBool("relay")
	if err != nil {
		return fmt.Errorf("failed to check relay flag: [%w]", err)
	}

	if launchRelay {
		go maintainer.LaunchRelay(ctx)
	}

	// TODO: Launch other maintainer tasks if necessary, e.g. spv. If no task
	//       has been specified - launch all the maintainer tasks.
	// TODO: Cancel all launched tasks if one of the tasks is unable to be
	//       launched, e.g. due to configuration errors.

	<-ctx.Done()
	return fmt.Errorf("unexpected context cancellation")
}
