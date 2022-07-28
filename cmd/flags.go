package cmd

import (
	"github.com/keep-network/keep-common/pkg/chain/ethereum"
	"github.com/keep-network/keep-common/pkg/chain/ethereum/ethutil"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

var EthereumFlags = []cli.Flag{
	altsrc.NewStringFlag(EthereumURLFlag),
	altsrc.NewPathFlag(EthereumKeyFileFlag),
	altsrc.NewDurationFlag(EthereumMiningCheckInterval),
	altsrc.NewGenericFlag(EthereumMaxGasFeeCap),
	altsrc.NewIntFlag(EthereumRequestPerSecondLimit),
	altsrc.NewIntFlag(EthereumConcurrencyLimit),
	altsrc.NewGenericFlag(EthereumBalanceAlertThreshold),
}

var NetworkFlags = []cli.Flag{
	NetworkPeersFlag, // TODO: Implement reading from file to workaround https://github.com/urfave/cli/discussions/1443
	altsrc.NewIntFlag(NetworkPortFlag),
	NetworkAnnouncedAddressesFlag, // TODO: Implement reading from file to workaround https://github.com/urfave/cli/discussions/1443
	altsrc.NewIntFlag(NetworkDisseminationTimeFlag),
}

var (
	ConfigFileFlag = &cli.PathFlag{
		Name:        "config",
		Aliases:     []string{"c"},
		Destination: &configFilePath,
		Usage:       "Full path to the configuration `FILE`",
		Value:       defaultConfigPath,
		TakesFile:   true,
		Category:    GeneralCategory,
	}
	StorageDataDirFlag = &cli.PathFlag{
		Name:        "storage.dataDir",
		Aliases:     []string{"d"},
		Destination: &clientConfig.Storage.DataDir,
		Usage:       "`DIRECTORY` to store the Keep client key shares and other sensitive data",
		Category:    StorageCategory,
	}
	EthereumURLFlag = &cli.StringFlag{
		Name:        "ethereum.url",
		Destination: &clientConfig.Ethereum.URL,
		Usage:       "WS connection `URL` for Ethereum client",
		Category:    EthereumCategory,
	}
	EthereumKeyFileFlag = &cli.PathFlag{
		Name:        "ethereum.keyFile",
		Destination: &clientConfig.Ethereum.Account.KeyFile,
		Usage:       "The local filesystem path to the Keep operator account key `FILE`",
		TakesFile:   true,
		Category:    EthereumCategory,
	}
	EthereumMiningCheckInterval = &cli.DurationFlag{
		Name:        "ethereum.miningCheckInterval",
		Destination: &clientConfig.Ethereum.MiningCheckInterval,
		Usage:       "The time interval in seconds in which transaction mining status is checked. If the transaction is not mined within this time, the gas price is increased and transaction is resubmitted",
		Value:       ethutil.DefaultMiningCheckInterval,
		Category:    EthereumCategory,
	}
	EthereumMaxGasFeeCap = &cli.GenericFlag{
		Name:        "ethereum.maxGasFeeCap",
		Destination: &clientConfig.Ethereum.MaxGasFeeCap,
		Usage:       "The maximum gas fee the client is willing to pay for the transaction to be mined. If reached, no resubmission attempts are performed",
		Value:       ethutil.DefaultMaxGasFeeCap,
		DefaultText: "500 Gwei", // TODO: Implement ethereum.Wei type to String conversion with units
		Category:    EthereumCategory,
	}
	EthereumRequestPerSecondLimit = &cli.IntFlag{
		Name:        "ethereum.requestPerSecondLimit",
		Destination: &clientConfig.Ethereum.RequestsPerSecondLimit,
		Usage:       "Request per second limit for all types of Ethereum client requests",
		Value:       150,
		Category:    EthereumCategory,
	}
	EthereumConcurrencyLimit = &cli.IntFlag{
		Name:        "ethereum.concurrencyLimit",
		Destination: &clientConfig.Ethereum.ConcurrencyLimit,
		Usage:       "The maximum number of concurrent requests which can be executed against Ethereum client",
		Value:       30,
		Category:    EthereumCategory,
	}
	EthereumBalanceAlertThreshold = &cli.GenericFlag{
		Name:        "ethereum.balanceAlertThreshold",
		Destination: &clientConfig.Ethereum.BalanceAlertThreshold,
		Usage:       "The minimum balance of operator account below which client starts reporting errors in logs",
		Value:       ethereum.WeiFromString("0.5 ether"),
		DefaultText: "0.5 ether",
		Category:    EthereumCategory,
	}
	// TODO: Figure out password
	// EthereumKeyFilePasswordFlag = &cli.StringFlag{
	// 	Destination: &clientConfig.Ethereum.Account.KeyFilePassword,
	// 	EnvVars:     []string{config.EthereumPasswordEnvVariable},
	// 	Required:    true,
	// 	Usage:       "aaa",
	// 	Category:    EthereumCategory,
	// }
	NetworkPeersFlag = &cli.MultiStringFlag{
		Target: &cli.StringSliceFlag{
			Name:  "network.peers",
			Usage: "Addresses of the network bootstrap nodes",
			// Value: , TODO: Load default peers addresses
			DefaultText: "auto",
			Category:    NetworkCategory,
		},
		Destination: &clientConfig.LibP2P.Peers,
	}
	NetworkPortFlag = &cli.IntFlag{
		Name:        "network.port",
		Aliases:     []string{"p"},
		Usage:       "Keep client listening `port`",
		Value:       1234, // TODO: Configure in lipp2p package
		Destination: &clientConfig.LibP2P.Port,
		Category:    NetworkCategory,
	}
	NetworkAnnouncedAddressesFlag = &cli.MultiStringFlag{
		Target: &cli.StringSliceFlag{
			Name:  "network.announcedAddresses",
			Usage: "Overwrites the default Keep client `addresses` announced in the network. Should be used for NAT or when more advanced firewall rules are applied",
			// Value: , TODO: Load default announced addresses
			DefaultText: "auto",
			Category:    NetworkCategory,
		},
		Destination: &clientConfig.LibP2P.AnnouncedAddresses,
	}
	NetworkDisseminationTimeFlag = &cli.IntFlag{
		Name:        "network.disseminationTime",
		Usage:       "Specifies courtesy message dissemination time in `seconds` for topics the node is not subscribed to. Should be used only on selected bootstrap nodes",
		Value:       0,
		DefaultText: "none",
		Destination: &clientConfig.LibP2P.DisseminationTime,
		Category:    NetworkCategory,
	}
)
