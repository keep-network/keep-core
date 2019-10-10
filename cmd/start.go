package cmd

import (
	"context"
	"fmt"

	"github.com/keep-network/keep-common/pkg/chain/ethereum/ethutil"
	"github.com/keep-network/keep-common/pkg/persistence"
	"github.com/keep-network/keep-core/config"
	"github.com/keep-network/keep-core/pkg/beacon"
	"github.com/keep-network/keep-core/pkg/chain/ethereum"
	"github.com/keep-network/keep-core/pkg/net/key"
	"github.com/keep-network/keep-core/pkg/net/libp2p"
	"github.com/keep-network/keep-core/pkg/operator"
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

	if c.Int(portFlag) > 0 {
		config.LibP2P.Port = c.Int(portFlag)
	}

	operatorPrivateKey, operatorPublicKey, err := loadStaticKey(config.Ethereum.Account)
	if err != nil {
		return fmt.Errorf("error loading static peer's key [%v]", err)
	}

	chainProvider, err := ethereum.Connect(config.Ethereum)
	if err != nil {
		return fmt.Errorf("error connecting to Ethereum node: [%v]", err)
	}

	stakeMonitor, err := chainProvider.StakeMonitor()
	if err != nil {
		return fmt.Errorf("error obtaining stake monitor handle [%v]", err)
	}
	hasMinimumStake, err := stakeMonitor.HasMinimumStake(
		config.Ethereum.Account.Address,
	)
	if err != nil {
		return fmt.Errorf("could not check the stake [%v]", err)
	}
	if !hasMinimumStake {
		return fmt.Errorf("stake is below the required minimum")
	}

	ctx := context.Background()
	networkPrivateKey, _ := key.OperatorKeyToNetworkKey(
		operatorPrivateKey, operatorPublicKey,
	)
	netProvider, err := libp2p.Connect(
		ctx,
		config.LibP2P,
		networkPrivateKey,
		stakeMonitor,
	)
	if err != nil {
		return err
	}

	nodeHeader(netProvider.AddrStrings(), config.LibP2P.Port)

	persistence := persistence.NewEncryptedPersistence(
		persistence.NewDiskHandle(config.Storage.DataDir),
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

	select {
	case <-ctx.Done():
		if err != nil {
			return err
		}

		return fmt.Errorf("uh-oh, we went boom boom for no reason")
	}
}

func loadStaticKey(
	account ethereum.Account,
) (*operator.PrivateKey, *operator.PublicKey, error) {
	ethereumKey, err := ethutil.DecryptKeyFile(
		account.KeyFile,
		account.KeyFilePassword,
	)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"failed to read KeyFile: %s [%v]", account.KeyFile, err,
		)
	}

	privateKey, publicKey := operator.EthereumKeyToOperatorKey(ethereumKey)

	return privateKey, publicKey, nil
}
