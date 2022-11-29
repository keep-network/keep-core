package cmd

import (
	"fmt"
	"math/big"

	"github.com/spf13/cobra"

	commonEthereum "github.com/keep-network/keep-common/pkg/chain/ethereum"
	"github.com/keep-network/keep-common/pkg/chain/ethereum/ethutil"
	"github.com/keep-network/keep-common/pkg/cmd/flag"
	"github.com/keep-network/keep-common/pkg/rate"
	"github.com/keep-network/keep-core/config"
	chainEthereum "github.com/keep-network/keep-core/pkg/chain/ethereum"
	"github.com/keep-network/keep-core/pkg/clientinfo"
	"github.com/keep-network/keep-core/pkg/net/libp2p"
	"github.com/keep-network/keep-core/pkg/tbtc"
)

func initGlobalFlags(
	cmd *cobra.Command,
	configFilePath *string,
) {
	initGlobalConfigFlags(cmd, configFilePath)
	initGlobalEthereumFlags(cmd)
}

func initFlags(
	cmd *cobra.Command,
	configFilePath *string,
	cfg *config.Config,
	categories ...config.Category,
) {
	for _, category := range categories {
		switch category {
		case config.Ethereum:
			initEthereumFlags(cmd, cfg)
		case config.BitcoinElectrum:
			initBitcoinElectrumFlags(cmd, cfg)
		case config.Network:
			initNetworkFlags(cmd, cfg)
		case config.Storage:
			initStorageFlags(cmd, cfg)
		case config.ClientInfo:
			initClientInfoFlags(cmd, cfg)
		case config.Tbtc:
			initTbtcFlags(cmd, cfg)
		case config.Developer:
			initDeveloperFlags(cmd)
		}
	}

	// Display flags in help in the same order they are defined. By default the
	// flags are ordered alphabetically which reduces readability.
	cmd.Flags().SortFlags = false
}

// Initialize flag for configuration file path.
func initGlobalConfigFlags(cmd *cobra.Command, configFilePath *string) {
	cmd.PersistentFlags().StringVarP(
		configFilePath,
		"config",
		"c",
		"", // Don't define default value as it would fail configuration reading.
		"Path to the configuration file. Supported formats: TOML, YAML, JSON.",
	)
}

// Initializes boolean flags for Ethereum network configuration. The flags can be used
// to run a client for a specific Ethereum network, e.g. add `--goerli` to the client
// start command to run the client against Görli Ethereum network. Only one flag
// from this set is allowed.
func initGlobalEthereumFlags(cmd *cobra.Command) {
	// TODO: Consider removing `--mainnet` flag. For now it's here to reduce a confusion
	// when developing and testing the client.
	cmd.PersistentFlags().Bool(
		commonEthereum.Mainnet.String(),
		false,
		"Mainnet network",
	)

	cmd.PersistentFlags().Bool(
		commonEthereum.Goerli.String(),
		false,
		"Görli network",
	)

	cmd.PersistentFlags().Bool(
		commonEthereum.Developer.String(),
		false,
		"Developer network",
	)

	cmd.MarkFlagsMutuallyExclusive(
		commonEthereum.Mainnet.String(),
		commonEthereum.Goerli.String(),
		commonEthereum.Developer.String(),
	)
}

// Initialize flags for Ethereum configuration.
func initEthereumFlags(cmd *cobra.Command, cfg *config.Config) {
	cmd.Flags().StringVar(
		&cfg.Ethereum.URL,
		"ethereum.url",
		"",
		"WS connection URL for Ethereum client.",
	)

	cmd.Flags().StringVar(
		&cfg.Ethereum.Account.KeyFile,
		"ethereum.keyFile",
		"",
		"The local filesystem path to Keep operator account keyfile.",
	)

	cmd.Flags().DurationVar(
		&cfg.Ethereum.MiningCheckInterval,
		"ethereum.miningCheckInterval",
		ethutil.DefaultMiningCheckInterval,
		"The time interval in seconds in which transaction mining status is checked. If the transaction is not mined within this time, the gas price is increased and transaction is resubmitted.",
	)

	flag.WeiVarFlag(
		cmd.Flags(),
		&cfg.Ethereum.MaxGasFeeCap,
		"ethereum.maxGasFeeCap",
		ethutil.DefaultMaxGasFeeCap,
		"The maximum gas fee the client is willing to pay for the transaction to be mined. If reached, no resubmission attempts are performed.",
	)

	cmd.Flags().IntVar(
		&cfg.Ethereum.RequestsPerSecondLimit,
		"ethereum.requestPerSecondLimit",
		rate.DefaultRequestsPerSecondLimit,
		"Request per second limit for all types of Ethereum client requests.",
	)

	cmd.Flags().IntVar(
		&cfg.Ethereum.ConcurrencyLimit,
		"ethereum.concurrencyLimit",
		rate.DefaultConcurrencyLimit,
		"The maximum number of concurrent requests which can be executed against Ethereum client.",
	)

	flag.WeiVarFlag(
		cmd.Flags(),
		&cfg.Ethereum.BalanceAlertThreshold,
		"ethereum.balanceAlertThreshold",
		*commonEthereum.WrapWei(big.NewInt(500000000000000000)), // 0.5 ether
		"The minimum balance of operator account below which client starts reporting errors in logs.",
	)
}

// Initialize flags for Bitcoin electrum configuration.
func initBitcoinElectrumFlags(cmd *cobra.Command, cfg *config.Config) {
	cmd.Flags().StringVar(
		&cfg.Bitcoin.Electrum.URL,
		"bitcoin.electrum.url",
		"",
		"URL for Bitcoin electrum client connection.",
	)
	// TODO: Initialize other config options from `electrum.Config`
}

// Initialize flags for Network configuration.
func initNetworkFlags(cmd *cobra.Command, cfg *config.Config) {
	cmd.Flags().BoolVar(
		&cfg.LibP2P.Bootstrap,
		"network.bootstrap",
		false,
		"Run the client in bootstrap mode.",
	)

	cmd.Flags().StringSliceVar(
		&cfg.LibP2P.Peers,
		"network.peers",
		[]string{},
		"Addresses of the network bootstrap nodes.",
	)

	cmd.Flags().IntVarP(
		&cfg.LibP2P.Port,
		"network.port",
		"p",
		libp2p.DefaultPort,
		"Keep client listening port.",
	)

	cmd.Flags().StringSliceVar(
		&cfg.LibP2P.AnnouncedAddresses,
		"network.announcedAddresses",
		[]string{},
		"Overwrites the default Keep client address announced in the network. Should be used for NAT or when more advanced firewall rules are applied.",
	)

	cmd.Flags().IntVar(
		&cfg.LibP2P.DisseminationTime,
		"network.disseminationTime",
		0,
		"Specifies courtesy message dissemination time in seconds for topics the node is not subscribed to. Should be used only on selected bootstrap nodes. (0 = none)",
	)
}

// Initialize flags for Storage configuration.
func initStorageFlags(cmd *cobra.Command, cfg *config.Config) {
	cmd.Flags().StringVar(
		&cfg.Storage.Dir,
		"storage.dir",
		"",
		"Location to store the Keep client key shares and other sensitive data.",
	)
}

// Initialize flags for ClientInfo configuration.
func initClientInfoFlags(cmd *cobra.Command, cfg *config.Config) {
	cmd.Flags().IntVar(
		&cfg.ClientInfo.Port,
		"clientInfo.port",
		9601,
		"Client Info HTTP server listening port.",
	)

	cmd.Flags().DurationVar(
		&cfg.ClientInfo.NetworkMetricsTick,
		"clientInfo.networkMetricsTick",
		clientinfo.DefaultNetworkMetricsTick,
		"Client Info network metrics check tick in seconds.",
	)

	cmd.Flags().DurationVar(
		&cfg.ClientInfo.EthereumMetricsTick,
		"clientInfo.ethereumMetricsTick",
		clientinfo.DefaultEthereumMetricsTick,
		"Client info Ethereum metrics check tick in seconds.",
	)
}

func initTbtcFlags(cmd *cobra.Command, cfg *config.Config) {
	cmd.Flags().IntVar(
		&cfg.Tbtc.PreParamsPoolSize,
		"tbtc.preParamsPoolSize",
		tbtc.DefaultPreParamsPoolSize,
		"tECDSA pre-parameters pool size.",
	)

	cmd.Flags().DurationVar(
		&cfg.Tbtc.PreParamsGenerationTimeout,
		"tbtc.preParamsGenerationTimeout",
		tbtc.DefaultPreParamsGenerationTimeout,
		"tECDSA pre-parameters generation timeout.",
	)

	cmd.Flags().DurationVar(
		&cfg.Tbtc.PreParamsGenerationDelay,
		"tbtc.preParamsGenerationDelay",
		tbtc.DefaultPreParamsGenerationDelay,
		"tECDSA pre-parameters generation delay.",
	)

	cmd.Flags().IntVar(
		&cfg.Tbtc.PreParamsGenerationConcurrency,
		"tbtc.preParamsGenerationConcurrency",
		tbtc.DefaultPreParamsGenerationConcurrency,
		"tECDSA pre-parameters generation concurrency.",
	)

	cmd.Flags().IntVar(
		&cfg.Tbtc.KeyGenerationConcurrency,
		"tbtc.keyGenerationConcurrency",
		tbtc.DefaultKeyGenerationConcurrency,
		"tECDSA key generation concurrency.",
	)

	cmd.Flags().BoolVar(
		&cfg.Tbtc.Maintainer.BitcoinDifficulty,
		"bitcoinDifficulty",
		false,
		"start Bitcoin difficulty maintainer",
	)
}

// Initialize flags for Developer configuration.
func initDeveloperFlags(command *cobra.Command) {
	initContractAddressFlag := func(contractName string) {
		command.Flags().String(
			config.GetDeveloperContractAddressKey(contractName),
			"",
			fmt.Sprintf(
				"Address of the %s smart contract",
				contractName,
			),
		)
	}

	initContractAddressFlag(chainEthereum.BridgeContractName)
	initContractAddressFlag(chainEthereum.LightRelayContractName)
	initContractAddressFlag(chainEthereum.RandomBeaconContractName)
	initContractAddressFlag(chainEthereum.TokenStakingContractName)
	initContractAddressFlag(chainEthereum.WalletRegistryContractName)
}
