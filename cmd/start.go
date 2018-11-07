package cmd

import (
	"context"
	"fmt"

	"github.com/keep-network/keep-core/config"
	"github.com/keep-network/keep-core/pkg/beacon"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/chain/ethereum"
	"github.com/keep-network/keep-core/pkg/net/key"
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
				&cli.IntFlag{
					Name: portFlag + "," + portShort,
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

	var port int
	if c.Int(portFlag) > 0 {
		config.LibP2P.Port = c.Int(portFlag)
	}

	staticKey, err := loadStaticKey(config.Ethereum.Account)
	if err != nil {
		return fmt.Errorf("error loading static peer's key [%v]", err)
	}

	ctx := context.Background()
	netProvider, err := libp2p.Connect(ctx, config.LibP2P, staticKey)
	if err != nil {
		return err
	}

	isBootstrapNode := config.LibP2P.Seed != 0
	nodeHeader(isBootstrapNode, netProvider.AddrStrings(), port)

	chainProvider, err := ethereum.Connect(config.Ethereum)
	if err != nil {
		return fmt.Errorf("error connecting to Ethereum node: [%v]", err)
	}

	blockCounter, err := chainProvider.BlockCounter()
	if err != nil {
		return fmt.Errorf("error initializing blockcounter: [%v]", err)
	}

	hasMinimumStake, err := checkMinimumStake(
		chainProvider,
		config.Ethereum.Account,
	)
	if err != nil {
		return fmt.Errorf("could not check the KEEP token stake [%v]", err)
	}
	if !hasMinimumStake {
		return fmt.Errorf("KEEP token stake is below the required minimum")
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

func loadStaticKey(account ethereum.Account) (*key.StaticNetworkKey, error) {
	ethereumKey, err := ethereum.DecryptKeyFile(
		account.KeyFile,
		account.KeyFilePassword,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to read KeyFile: %s [%v]", account.KeyFile, err,
		)
	}

	return key.EthereumKeyToNetworkKey(ethereumKey), nil
}

func checkMinimumStake(
	chain chain.Handle,
	account ethereum.Account,
) (bool, error) {
	stakeMonitor, err := chain.StakeMonitor()
	if err != nil {
		return false, fmt.Errorf("error initializing stake monitor: [%v]", err)
	}

	return stakeMonitor.HasMinimumStake(account.Address)
}
