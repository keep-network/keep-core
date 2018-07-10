package cmd

import (
	"context"
	"fmt"

	"github.com/keep-network/keep-core/config"
	"github.com/keep-network/keep-core/pkg/beacon"
	"github.com/keep-network/keep-core/pkg/chain/ethereum"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/net/libp2p"
	"github.com/urfave/cli"
)

const (
	sampleText             = "sample text"
	broadcastChannelName   = "test"
	resetBroadcastTimerSec = 5
)

// StartFlags for bootstrap and port
var StartFlags []cli.Flag

type recvParams struct {
	port     int
	ipaddr   string
	recvChan chan net.Message
}

type broadcastParams struct {
	port      int
	ipaddr    string
	bcastChan net.BroadcastChannel
}

func init() {
	StartFlags = []cli.Flag{
		&cli.BoolFlag{
			Name: "bootstrap",
		},
		&cli.IntFlag{
			Name: "port",
		},
	}
}

// StartNode starts a node; if it's not a bootstrap node it will get the Node.URLs from the config file
func StartNode(c *cli.Context) error {
	cfg, err := config.ReadConfig(c.GlobalString("config"))
	if err != nil {
		return fmt.Errorf("error reading config file: %v", err)
	}

	var port int
	if c.Int("port") > 0 {
		port = c.Int("port")
	} else {
		port = cfg.Node.Port
	}

	var (
		seed          int
		bootstrapURLs []string
	)
	if c.Bool("bootstrap") {
		seed = cfg.Bootstrap.Seed
	} else {
		bootstrapURLs = cfg.Bootstrap.URLs
	}

	ctx := context.Background()
	netProvider, err := libp2p.Connect(ctx, &libp2p.Config{
		Port:  port,
		Peers: bootstrapURLs,
		Seed:  seed,
	})
	if err != nil {
		return err
	}

	nodeHeader(c.Bool("bootstrap"), netProvider.Addrs(), port)

	chainProvider, err := ethereum.Connect(cfg.Ethereum)
	if err != nil {
		return fmt.Errorf("error connecting to Ethereum node: [%v]", err)
	}

	blockCounter, err := chainProvider.BlockCounter()
	if err != nil {
		return fmt.Errorf("error initializing blockcounter: [%v]", err)
	}

	err = beacon.Initialize(
		ctx,
		chainProvider.ThresholdRelay(),
		blockCounter,
		netProvider,
	)
	if err != nil {
		return fmt.Errorf("error initializing beacon: [%v]", err)
	}

	select {
	case <-ctx.Done():
		if err != nil {
			return err
		}

		return fmt.Errorf("uh-oh, we went boom boom for no reason")
	}
}
