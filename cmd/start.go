package cmd

import (
	"context"
	"fmt"
	"time"

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
	config, err := config.ReadConfig(c.GlobalString("config"))
	if err != nil {
		return fmt.Errorf("error reading config file: %v", err)
	}

	if c.Int(portFlag) > 0 {
		config.LibP2P.Port = c.Int(portFlag)
	}

	// FIXME This needs to happen inside the `pkg/chain/ethereum` scope,
	// FIXME probably.
	operatorPrivateKey, operatorPublicKey, err := loadStaticKey(
		config.Ethereum.Account.KeyFile,
		config.Ethereum.Account.KeyFilePassword,
	)
	if err != nil {
		return fmt.Errorf("error loading static peer's key [%v]", err)
	}

	chainProvider, err := ethereum.Connect(config.Ethereum)
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
	if c.Int(waitForStakeFlag) != 0 {
		err = waitForStake(stakeMonitor, config.Ethereum.Account.Address, c.Int(waitForStakeFlag))
		if err != nil {
			return err
		}
	}
	hasMinimumStake, err := stakeMonitor.HasMinimumStake(
		config.Ethereum.Account.Address,
	)
	if err != nil {
		return fmt.Errorf("could not check the stake [%v]", err)
	}
	if !hasMinimumStake {
		return fmt.Errorf(
			"no minimum KEEP stake or operator is not authorized to use it; " +
				"please make sure the operator address in the configuration " +
				"is correct and it has KEEP tokens delegated and the operator " +
				"contract has been authorized to operate on the stake",
		)
	}

	ctx := context.Background()
	networkPrivateKey, _ := key.OperatorKeyToNetworkKey(
		operatorPrivateKey, operatorPublicKey,
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
		config.Ethereum.Account.Address,
		chainProvider,
		netProvider,
		persistence,
	)
	if err != nil {
		return fmt.Errorf("error initializing beacon: [%v]", err)
	}

	initializeMetrics(ctx, config, netProvider, stakeMonitor)

	select {
	case <-ctx.Done():
		if err != nil {
			return err
		}

		return fmt.Errorf("uh-oh, we went boom boom for no reason")
	}
}

func loadStaticKey(
	keyFile string,
	keyFilePassword string,
) (*operator.PrivateKey, *operator.PublicKey, error) {
	ethereumKey, err := ethutil.DecryptKeyFile(
		keyFile,
		keyFilePassword,
	)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"failed to read KeyFile: %s [%v]", keyFile, err,
		)
	}

	privateKey, publicKey := operator.EthereumKeyToOperatorKey(ethereumKey)

	return privateKey, publicKey, nil
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
		config.Ethereum.Account.Address,
		time.Duration(config.Metrics.EthereumMetricsTick)*time.Second,
	)

	metrics.ExposeLibP2PInfo(
		registry,
		netProvider,
	)
}
