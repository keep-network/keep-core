package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/keep-network/keep-core/pkg/diagnostics"
	"github.com/keep-network/keep-core/pkg/metrics"
	"github.com/keep-network/keep-core/pkg/net"

	"github.com/ipfs/go-log"
	"github.com/keep-network/keep-common/pkg/chain/ethereum/ethutil"
	"github.com/keep-network/keep-common/pkg/persistence"
	"github.com/keep-network/keep-core/config"
	"github.com/keep-network/keep-core/pkg/beacon"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/chain/ethereum"
	"github.com/keep-network/keep-core/pkg/firewall"
	"github.com/keep-network/keep-core/pkg/net/key"
	"github.com/keep-network/keep-core/pkg/net/libp2p"
	"github.com/keep-network/keep-core/pkg/net/retransmission"
	"github.com/keep-network/keep-core/pkg/operator"
	"github.com/urfave/cli"
)

// StartCommand contains the definition of the start command-line subcommand.
var (
	StartCommand cli.Command
	logger       = log.Logger("keep-start")
)

const (
	bootstrapFlag     = "bootstrap"
	portFlag          = "port"
	portShort         = "p"
	waitForStakeFlag  = "wait-for-stake"
	waitForStakeShort = "w"
)

const startDescription = `Starts the Keep client in the foreground. Currently this only consists of the
   threshold relay client for the Keep random beacon.`

func init() {
	StartCommand =
		cli.Command{
			Name:        "start",
			Usage:       `Starts the Keep client in the foreground`,
			Description: startDescription,
			Action:      Start,
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name: portFlag + "," + portShort,
				},
				&cli.IntFlag{
					Name: waitForStakeFlag + "," + waitForStakeShort,
				},
			},
		}
}

// Start starts a node; if it's not a bootstrap node it will get the Node.URLs
// from the config file
func Start(c *cli.Context) error {
	ctx := context.Background()

	config, err := config.ReadConfig(c.GlobalString("config"))
	if err != nil {
		return fmt.Errorf("error reading config file: %v", err)
	}

	if c.Int(portFlag) > 0 {
		config.LibP2P.Port = c.Int(portFlag)
	}

	ethereumKey, err := ethutil.DecryptKeyFile(
		config.Ethereum.Account.KeyFile,
		config.Ethereum.Account.KeyFilePassword,
	)
	if err != nil {
		return fmt.Errorf(
			"failed to read key file [%s]: [%v]",
			config.Ethereum.Account.KeyFile,
			err,
		)
	}

	chainProvider, err := ethereum.Connect(ctx, config.Ethereum)
	if err != nil {
		return fmt.Errorf("error connecting to Ethereum node: [%v]", err)
	}

	blockCounter, err := chainProvider.BlockCounter()
	if err != nil {
		return err
	}

	stakeMonitor, err := chainProvider.StakeMonitor()
	if err != nil {
		return fmt.Errorf("error obtaining stake monitor handle [%v]", err)
	}
	// FIXME: Update for V2
	// if c.Int(waitForStakeFlag) != 0 {
	// 	err = waitForStake(stakeMonitor, ethereumKey.Address.Hex(), c.Int(waitForStakeFlag))
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	// hasMinimumStake, err := stakeMonitor.HasMinimumStake(
	// 	ethereumKey.Address.Hex(),
	// )
	// if err != nil {
	// 	return fmt.Errorf("could not check the stake [%v]", err)
	// }
	// if !hasMinimumStake {
	// 	return fmt.Errorf(
	// 		"no minimum KEEP stake or operator is not authorized to use it; " +
	// 			"please make sure the operator address in the configuration " +
	// 			"is correct and it has KEEP tokens delegated and the operator " +
	// 			"contract has been authorized to operate on the stake",
	// 	)
	// }

	networkPrivateKey, _ := key.OperatorKeyToNetworkKey(
		operator.ChainKeyToOperatorKey(ethereumKey),
	)
	netProvider, err := libp2p.Connect(
		ctx,
		config.LibP2P,
		networkPrivateKey,
		libp2p.ProtocolBeacon,
		firewall.MinimumStakePolicy(stakeMonitor),
		retransmission.NewTicker(blockCounter.WatchBlocks(ctx)),
	)
	if err != nil {
		return err
	}

	nodeHeader(netProvider.ConnectionManager().AddrStrings(), config.LibP2P.Port)

	handle, err := persistence.NewDiskHandle(config.Storage.DataDir)
	if err != nil {
		return fmt.Errorf("failed while creating a storage disk handler: [%v]", err)
	}
	persistence := persistence.NewEncryptedPersistence(
		handle,
		config.Ethereum.Account.KeyFilePassword,
	)

	err = beacon.Initialize(
		ctx,
		ethereumKey.Address.Hex(),
		chainProvider,
		netProvider,
		persistence,
	)
	if err != nil {
		return fmt.Errorf("error initializing beacon: [%v]", err)
	}

	initializeMetrics(ctx, config, netProvider, stakeMonitor, ethereumKey.Address.Hex())
	initializeDiagnostics(ctx, config, netProvider)

	select {
	case <-ctx.Done():
		if err != nil {
			return err
		}

		return fmt.Errorf("uh-oh, we went boom boom for no reason")
	}
}

func waitForStake(stakeMonitor chain.StakeMonitor, address string, timeout int) error {
	waitMins := 0
	for waitMins < timeout {
		hasMinimumStake, err := stakeMonitor.HasMinimumStake(address)
		if err != nil {
			return fmt.Errorf("could not check the stake [%v]", err)
		}
		if hasMinimumStake {
			return nil
		}
		logger.Warningf("%s below min stake for %d min \n", address, waitMins)
		time.Sleep(time.Minute)
		waitMins++
	}
	return fmt.Errorf("timed out waiting for %s to have required minimum stake", address)
}

func initializeMetrics(
	ctx context.Context,
	config *config.Config,
	netProvider net.Provider,
	stakeMonitor chain.StakeMonitor,
	ethereumAddress string,
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
		stakeMonitor,
		ethereumAddress,
		time.Duration(config.Metrics.EthereumMetricsTick)*time.Second,
	)
}

func initializeDiagnostics(
	ctx context.Context,
	config *config.Config,
	netProvider net.Provider,
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

	diagnostics.RegisterConnectedPeersSource(registry, netProvider)
	diagnostics.RegisterClientInfoSource(registry, netProvider)
}
