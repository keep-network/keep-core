package cmd

import (
	"context"
	"fmt"

	"github.com/keep-network/keep-core/pkg/maintainer/config"
	"github.com/spf13/cobra"
)

var maintainerConfig = &config.Config{}

// RelayCommand contains the definition of the relay command-line subcommand.
var RelayCommand = &cobra.Command{
	Use:   "relay",
	Short: `Starts relay maintainer`,
	Long:  relayDescription,
	PreRun: func(cmd *cobra.Command, args []string) {
		if err := maintainerConfig.ReadConfig(configFilePath); err != nil {
			logger.Fatalf("error reading config: %v", err)
		}
	},
	RunE: relay,
}

const relayDescription = `The relay command starts relay maintainer`

func relay(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// TODO: uncomment when the electrum implementation is ready
	// btcChain, err := bitcoin.Connect(ctx, &maintainerConfig.Bitcoin)
	// if err != nil {
	// 	return fmt.Errorf("could not connect BTC chain: [%v]", err)
	// }

	// maintainer.NewRelay(btcChain, nil)

	// 	TODO: add metrics
	logger.Info("relay started")

	<-ctx.Done()
	return fmt.Errorf("unexpected context cancellation")
}
