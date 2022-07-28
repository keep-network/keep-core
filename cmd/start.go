package cmd

import (
	"context"
	"fmt"
	"github.com/keep-network/keep-core/pkg/tbtc"
	"time"

	"github.com/keep-network/keep-core/config"
	"github.com/keep-network/keep-core/pkg/chain/ethereum"
	"github.com/keep-network/keep-core/pkg/diagnostics"
	"github.com/keep-network/keep-core/pkg/metrics"
	"github.com/keep-network/keep-core/pkg/net"

	"github.com/ipfs/go-log"
	"github.com/keep-network/keep-common/pkg/persistence"
	"github.com/keep-network/keep-core/pkg/beacon"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/firewall"
	"github.com/keep-network/keep-core/pkg/net/libp2p"
	"github.com/keep-network/keep-core/pkg/net/retransmission"

	"github.com/spf13/cobra"
)

var (
	// StartCommand contains the definition of the start command-line subcommand.
	StartCommand *cobra.Command

	logger = log.Logger("keep-start")
)

func init() {
	StartCommand = &cobra.Command{
		Use:   "start",
		Short: "Starts the Keep Client",
		Long:  "Starts the Keep Client in the foreground",
		PreRun: func(cmd *cobra.Command, args []string) {
			if err := clientConfig.ReadConfig(configFilePath, cmd.Flags()); err != nil {
				logger.Fatalf("error reading config: %v", err)
			}

		},
		Run: func(cmd *cobra.Command, args []string) {
			if err := start(cmd); err != nil {
				logger.Fatal(err)
			}
		},
	}

	initFlags(StartCommand, allCategories, &configFilePath, clientConfig)

	StartCommand.SetUsageTemplate(
		fmt.Sprintf(`%s
Environment variables:
    %s    Password for Keep operator account keyfile decryption.
    %s                 Space-delimited set of log level directives; set to "help" for help.
`,
			StartCommand.UsageString(),
			config.EthereumPasswordEnvVariable,
			config.LogLevelEnvVariable,
		),
	)
}

// start starts a node
func start(cmd *cobra.Command) error {
	ctx := context.Background()

	beaconChain, tbtcChain, blockCounter, signing, operatorPrivateKey, err :=
		ethereum.Connect(ctx, clientConfig.Ethereum)
	if err != nil {
		return fmt.Errorf("error connecting to Ethereum node: [%v]", err)
	}

	firewall := firewall.AnyApplicationPolicy(
		[]firewall.Application{beaconChain, tbtcChain},
	)

	netProvider, err := libp2p.Connect(
		ctx,
		clientConfig.LibP2P,
		operatorPrivateKey,
		firewall,
		retransmission.NewTicker(blockCounter.WatchBlocks(ctx)),
	)
	if err != nil {
		return fmt.Errorf("failed while creating the network provider: [%v]", err)
	}

	nodeHeader(netProvider.ConnectionManager().AddrStrings(), clientConfig.LibP2P.Port)

	handle, err := persistence.NewDiskHandle(clientConfig.Storage.DataDir)
	if err != nil {
		return fmt.Errorf("failed while creating a storage disk handler: [%v]", err)
	}
	encryptedPersistence := persistence.NewEncryptedPersistence(
		handle,
		clientConfig.Ethereum.Account.KeyFilePassword,
	)

	err = beacon.Initialize(
		ctx,
		beaconChain,
		netProvider,
		encryptedPersistence,
	)
	if err != nil {
		return fmt.Errorf("error initializing beacon: [%v]", err)
	}

	err = tbtc.Initialize(
		ctx,
		tbtcChain,
		netProvider,
		nil, // TODO: Pass a proper persistence handle.
	)
	if err != nil {
		return fmt.Errorf("error initializing TBTC: [%v]", err)
	}

	initializeMetrics(ctx, clientConfig, netProvider, blockCounter)
	initializeDiagnostics(ctx, clientConfig, netProvider, signing)

	select {
	case <-ctx.Done():
		if err != nil {
			return err
		}

		return fmt.Errorf("uh-oh, we went boom boom for no reason")
	}
}

func initializeMetrics(
	ctx context.Context,
	config *config.Config,
	netProvider net.Provider,
	blockCounter chain.BlockCounter,
) {
	registry, isConfigured := metrics.Initialize(
		config.Metrics.Port,
	)
	if !isConfigured {
		logger.Infof("metrics are not configured")
		return
	}

	logger.Infof(
		"enabled metrics on port [%v]",
		config.Metrics.Port,
	)

	metrics.ObserveConnectedPeersCount(
		ctx,
		registry,
		netProvider,
		time.Duration(config.Metrics.NetworkMetricsTick)*time.Second,
	)

	metrics.ObserveConnectedBootstrapCount(
		ctx,
		registry,
		netProvider,
		config.LibP2P.Peers,
		time.Duration(config.Metrics.NetworkMetricsTick)*time.Second,
	)

	metrics.ObserveEthConnectivity(
		ctx,
		registry,
		blockCounter,
		time.Duration(config.Metrics.EthereumMetricsTick)*time.Second,
	)
}

func initializeDiagnostics(
	ctx context.Context,
	config *config.Config,
	netProvider net.Provider,
	signing chain.Signing,
) {
	registry, isConfigured := diagnostics.Initialize(
		config.Diagnostics.Port,
	)
	if !isConfigured {
		logger.Infof("diagnostics are not configured")
		return
	}

	logger.Infof(
		"enabled diagnostics on port [%v]",
		config.Diagnostics.Port,
	)

	diagnostics.RegisterConnectedPeersSource(registry, netProvider, signing)
	diagnostics.RegisterClientInfoSource(registry, netProvider, signing)
}
