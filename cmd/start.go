package cmd

import (
	"context"
	"fmt"
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
	StartCommand = &cobra.Command{
		Use:   "start",
		Short: "Starts the Keep Client",
		Long:  startDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := start(cmd); err != nil {
				logger.Fatal(err)
			}
			return nil
		},
	}

	logger = log.Logger("keep-start")
)

const startDescription = `Starts the Keep Client in the foreground`

func init() {
	config.InitFlags(StartCommand)
}

// start starts a node
func start(cmd *cobra.Command) error {
	ctx := context.Background()

	filePath, err := cmd.Flags().GetString("config")
	if err != nil {
		return fmt.Errorf("error getting config flag: %w", err)
	}

	config, err := config.ReadConfig(filePath)
	if err != nil {
		return fmt.Errorf("error reading config: %w", err)
	}

	beaconChain, tbtcChain, err := ethereum.Connect(ctx, config.Ethereum)
	if err != nil {
		return fmt.Errorf("error connecting to Ethereum node: [%v]", err)
	}

	operatorPrivateKey, _, err := beaconChain.OperatorKeyPair()
	if err != nil {
		return fmt.Errorf("failed to get operator key pair: [%v]", err)
	}

	blockCounter, err := beaconChain.BlockCounter()
	if err != nil {
		return fmt.Errorf("failed to get block counter: [%v]", err)
	}

	firewall := firewall.AnyApplicationPolicy(
		[]firewall.Application{beaconChain, tbtcChain},
	)

	netProvider, err := libp2p.Connect(
		ctx,
		config.LibP2P,
		operatorPrivateKey,
		libp2p.ProtocolBeacon,
		firewall,
		retransmission.NewTicker(blockCounter.WatchBlocks(ctx)),
	)
	if err != nil {
		return fmt.Errorf("failed while creating the network provider: [%v]", err)
	}

	nodeHeader(netProvider.ConnectionManager().AddrStrings(), config.LibP2P.Port)

	handle, err := persistence.NewDiskHandle(config.Storage.DataDir)
	if err != nil {
		return fmt.Errorf("failed while creating a storage disk handler: [%v]", err)
	}
	encryptedPersistence := persistence.NewEncryptedPersistence(
		handle,
		config.Ethereum.Account.KeyFilePassword,
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

	initializeMetrics(ctx, config, netProvider, blockCounter)
	initializeDiagnostics(ctx, config, netProvider, beaconChain.Signing())

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
