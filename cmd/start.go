package cmd

import (
	"context"
	"fmt"

	"github.com/keep-network/keep-core/config"
	"github.com/keep-network/keep-core/pkg/beacon"
	"github.com/keep-network/keep-core/pkg/chain/ethereum"
	"github.com/keep-network/keep-core/pkg/net/libp2p"
	"github.com/urfave/cli"
)

// StartCommand contains the definition of the start command-line subcommand.
var StartCommand cli.Command

const (
	bootstrapFlag = "bootstrap"
	portFlag      = "port"
	portShort     = "p"
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
				&cli.BoolFlag{
					Name: bootstrapFlag,
				},
				&cli.IntFlag{
					Name: portFlag + "," + portShort,
				},
			},
		}
}

// Start starts a node; if it's not a bootstrap node it will get the Node.URLs
// from the config file
func Start(c *cli.Context) error {
	cfg, err := config.ReadConfig(c.GlobalString("config"))
	if err != nil {
		return fmt.Errorf("error reading config file: %v", err)
	}

	var port int
	if c.Int(portFlag) > 0 {
		port = c.Int(portFlag)
	} else {
		port = cfg.Node.Port
	}

	var (
		seed          int
		bootstrapURLs []string
	)
	if c.Bool(bootstrapFlag) {
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

	nodeHeader(c.Bool(bootstrapFlag), netProvider.Addrs(), port)

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
