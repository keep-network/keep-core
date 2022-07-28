package cmd

import (
	"context"
	"fmt"
	"time"

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
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

// StartCommand contains the definition of the start command-line subcommand.
var (
	StartCommand cli.Command
	logger       = log.Logger("keep-start")
)

const (
	portFlag  = "port"
	portShort = "p"
)

const startDescription = `Starts the Keep client in the foreground`

func init() {
	flags := append(
		[]cli.Flag{
			ConfigFileFlag,
			altsrc.NewPathFlag(StorageDataDirFlag),
		},
		EthereumFlags...,
	)

	flags = append(flags, NetworkFlags...)

	StartCommand =
		cli.Command{
			Name:        "start",
			Usage:       `Starts the Keep client in the foreground`,
			Description: startDescription,
			Action:      Start,
			Before:      altsrc.InitInputSourceWithContext(flags, altsrc.NewTomlSourceFromFlagFunc("config")),
			Flags:       flags,
		}
}

// Start starts a node; if it's not a bootstrap node it will get the node.URLs
// from the config file
func Start(c *cli.Context) error {
	if err := clientConfig.Validate(); err != nil {
		return err
	}

	ctx := context.Background()

	logger.Infof("TEST PEERS %v", clientConfig.LibP2P.Peers)

	beaconChain, tbtcChain, err := ethereum.Connect(ctx, clientConfig.Ethereum)
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
		clientConfig.LibP2P,
		operatorPrivateKey,
		libp2p.ProtocolBeacon,
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

	initializeMetrics(ctx, netProvider, blockCounter)
	initializeDiagnostics(ctx, netProvider, beaconChain.Signing())

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
	netProvider net.Provider,
	blockCounter chain.BlockCounter,
) {
	registry, isConfigured := metrics.Initialize(
		clientConfig.Metrics.Port,
	)
	if !isConfigured {
		logger.Infof("metrics are not configured")
		return
	}

	logger.Infof(
		"enabled metrics on port [%v]",
		clientConfig.Metrics.Port,
	)

	metrics.ObserveConnectedPeersCount(
		ctx,
		registry,
		netProvider,
		time.Duration(clientConfig.Metrics.NetworkMetricsTick)*time.Second,
	)

	metrics.ObserveConnectedBootstrapCount(
		ctx,
		registry,
		netProvider,
		clientConfig.LibP2P.Peers,
		time.Duration(clientConfig.Metrics.NetworkMetricsTick)*time.Second,
	)

	metrics.ObserveEthConnectivity(
		ctx,
		registry,
		blockCounter,
		time.Duration(clientConfig.Metrics.EthereumMetricsTick)*time.Second,
	)
}

func initializeDiagnostics(
	ctx context.Context,
	netProvider net.Provider,
	signing chain.Signing,
) {
	registry, isConfigured := diagnostics.Initialize(
		clientConfig.Diagnostics.Port,
	)
	if !isConfigured {
		logger.Infof("diagnostics are not configured")
		return
	}

	logger.Infof(
		"enabled diagnostics on port [%v]",
		clientConfig.Diagnostics.Port,
	)

	diagnostics.RegisterConnectedPeersSource(registry, netProvider, signing)
	diagnostics.RegisterClientInfoSource(registry, netProvider, signing)
}
