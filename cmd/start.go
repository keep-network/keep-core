package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/keep-network/keep-core/build"
	"github.com/keep-network/keep-core/config"
	"github.com/keep-network/keep-core/pkg/beacon"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/chain/ethereum"
	"github.com/keep-network/keep-core/pkg/clientinfo"
	"github.com/keep-network/keep-core/pkg/firewall"
	"github.com/keep-network/keep-core/pkg/generator"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/net/libp2p"
	"github.com/keep-network/keep-core/pkg/net/retransmission"
	"github.com/keep-network/keep-core/pkg/storage"
	"github.com/keep-network/keep-core/pkg/tbtc"
)

// StartCommand contains the definition of the start command-line subcommand.
var StartCommand = &cobra.Command{
	Use:   "start",
	Short: "Starts the Keep Client",
	Long:  "Starts the Keep Client in the foreground",
	PreRun: func(cmd *cobra.Command, args []string) {
		if err := clientConfig.ReadConfig(configFilePath, cmd.Flags(), config.StartCmdCategories...); err != nil {
			logger.Fatalf("error reading config: %v", err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if err := start(cmd); err != nil {
			logger.Fatal(err)
		}
	},
}

func init() {
	initFlags(StartCommand, &configFilePath, clientConfig, config.StartCmdCategories...)

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

	logger.Infof(
		"Starting the client against [%s] ethereum network...",
		clientConfig.Ethereum.Network,
	)

	beaconChain, tbtcChain, blockCounter, signing, operatorPrivateKey, err :=
		ethereum.Connect(ctx, clientConfig.Ethereum)
	if err != nil {
		return fmt.Errorf("error connecting to Ethereum node: [%v]", err)
	}

	bootstrapPeersPublicKeys, err := libp2p.ExtractPeersPublicKeys(
		clientConfig.LibP2P.Peers,
	)
	if err != nil {
		return fmt.Errorf(
			"error extracting bootstrap peers public keys: [%v]",
			err,
		)
	}

	firewall := firewall.AnyApplicationPolicy(
		[]firewall.Application{beaconChain, tbtcChain},
		firewall.NewAllowList(bootstrapPeersPublicKeys),
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

	nodeHeader(
		netProvider.ConnectionManager().AddrStrings(),
		beaconChain.Signing().Address().String(),
		clientConfig.LibP2P.Port,
		clientConfig.Ethereum,
	)

	clientInfoRegistry := initializeClientInfo(
		ctx,
		clientConfig,
		netProvider,
		signing,
		blockCounter,
	)

	// Initialize beacon and tbtc only for non-bootstrap nodes.
	// Skip initialization for bootstrap nodes as they are only used for network
	// discovery.
	if !clientConfig.LibP2P.Bootstrap {
		storage, err := storage.Initialize(
			clientConfig.Storage,
			clientConfig.Ethereum.KeyFilePassword,
		)
		if err != nil {
			return fmt.Errorf("cannot initialize storage: [%w]", err)
		}

		beaconKeyStorePersistence, err := storage.InitializeKeyStorePersistence(
			"beacon",
		)
		if err != nil {
			return fmt.Errorf(
				"cannot initialize beacon keystore persistence: [%w]",
				err,
			)
		}

		tbtcKeyStorePersistence, err := storage.InitializeKeyStorePersistence(
			"tbtc",
		)
		if err != nil {
			return fmt.Errorf(
				"cannot initialize tbtc keystore persistence: [%w]",
				err,
			)
		}

		tbtcDataPersistence, err := storage.InitializeWorkPersistence("tbtc")
		if err != nil {
			return fmt.Errorf(
				"cannot initialize tbtc data persistence: [%w]",
				err,
			)
		}

		scheduler := generator.StartScheduler()

		err = beacon.Initialize(
			ctx,
			beaconChain,
			netProvider,
			beaconKeyStorePersistence,
			scheduler,
		)
		if err != nil {
			return fmt.Errorf("error initializing beacon: [%v]", err)
		}

		err = tbtc.Initialize(
			ctx,
			tbtcChain,
			netProvider,
			tbtcKeyStorePersistence,
			tbtcDataPersistence,
			scheduler,
			clientConfig.Tbtc,
			clientInfoRegistry,
		)
		if err != nil {
			return fmt.Errorf("error initializing TBTC: [%v]", err)
		}
	}

	<-ctx.Done()
	return fmt.Errorf("shutting down the node because its context has ended")
}

func initializeClientInfo(
	ctx context.Context,
	config *config.Config,
	netProvider net.Provider,
	signing chain.Signing,
	blockCounter chain.BlockCounter,
) *clientinfo.Registry {
	registry, isConfigured := clientinfo.Initialize(ctx, config.ClientInfo.Port)
	if !isConfigured {
		logger.Infof("client info endpoint not configured")
		return nil
	}

	registry.ObserveConnectedPeersCount(
		netProvider,
		config.ClientInfo.NetworkMetricsTick,
	)

	registry.ObserveConnectedBootstrapCount(
		netProvider,
		config.LibP2P.Peers,
		config.ClientInfo.NetworkMetricsTick,
	)

	registry.ObserveEthConnectivity(
		blockCounter,
		config.ClientInfo.EthereumMetricsTick,
	)

	registry.RegisterMetricClientInfo(build.Version)

	registry.RegisterConnectedPeersSource(netProvider, signing)
	registry.RegisterClientInfoSource(
		netProvider,
		signing,
		build.Version,
		build.Revision,
	)

	logger.Infof(
		"enabled client info endpoint on port [%v]",
		config.ClientInfo.Port,
	)

	return registry
}
