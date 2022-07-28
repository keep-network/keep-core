package cmd

import (
	"fmt"
	"math/big"

	"github.com/keep-network/keep-common/pkg/chain/ethereum"
	"github.com/keep-network/keep-common/pkg/chain/ethereum/ethutil"
	"github.com/keep-network/keep-core/cmd/flag"
	"github.com/keep-network/keep-core/config"
	"github.com/keep-network/keep-core/pkg/net/libp2p"
	"github.com/spf13/cobra"

	chainEthereum "github.com/keep-network/keep-core/pkg/chain/ethereum"

	ethereumBeacon "github.com/keep-network/keep-core/pkg/chain/ethereum/beacon/gen"
	ethereumEcdsa "github.com/keep-network/keep-core/pkg/chain/ethereum/ecdsa/gen"
	ethereumThreshold "github.com/keep-network/keep-core/pkg/chain/ethereum/threshold/gen"
)

type category int

const (
	General category = iota
	Ethereum
	Network
	Storage
	Metrics
	Diagnostics
	Developer
)

var allCategories = []category{
	General,
	Ethereum,
	Network,
	Storage,
	Metrics,
	Diagnostics,
	Developer,
}

func initFlags(
	cmd *cobra.Command,
	categories []category,
	configFilePath *string,
	cfg *config.Config,
) {
	for _, category := range categories {
		switch category {
		case General:
			initConfigFlags(cmd, configFilePath)
			break
		case Ethereum:
			initEthereumFlags(cmd, cfg)
			break
		case Network:
			initNetworkFlags(cmd, cfg)
			break
		case Storage:
			initStorageFlags(cmd, cfg)
			break
		case Metrics:
			initMetricsFlags(cmd, cfg)
			break
		case Diagnostics:
			initDiagnosticsFlags(cmd, cfg)
			break
		case Developer:
			initDeveloperFlags(cmd)
			break
		}
	}

	// Display flags in help in the same order they are defined. By default the
	// flags are ordered alphabetically which reduces readability.
	cmd.Flags().SortFlags = false
}

// Initialize flag for configuration file path.
func initConfigFlags(cmd *cobra.Command, configFilePath *string) {
	cmd.Flags().StringVarP(
		configFilePath,
		"config",
		"c",
		"", // Don't define default value as it would fail configuration reading.
		"Path to the configuration file.",
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
		150, // TODO: Can it be read from some place?
		"Request per second limit for all types of Ethereum client requests.",
	)

	cmd.Flags().IntVar(
		&cfg.Ethereum.ConcurrencyLimit,
		"ethereum.concurrencyLimit",
		30, // TODO: Can it be read from some place?
		"The maximum number of concurrent requests which can be executed against Ethereum client.",
	)

	flag.WeiVarFlag(
		cmd.Flags(),
		&cfg.Ethereum.BalanceAlertThreshold,
		"ethereum.balanceAlertThreshold",
		*ethereum.WrapWei(big.NewInt(500000000000000000)), // 0.5 ether
		"The minimum balance of operator account below which client starts reporting errors in logs.",
	)
}

// Initialize flags for Network configuration.
func initNetworkFlags(cmd *cobra.Command, cfg *config.Config) {
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
		&cfg.Storage.DataDir,
		"storage.dataDir",
		"",
		"Location to store the Keep client key shares and other sensitive data.",
	)
}

// Initialize flags for Metrics configuration.
func initMetricsFlags(cmd *cobra.Command, cfg *config.Config) {
	cmd.Flags().IntVar(
		&cfg.Metrics.Port,
		"metrics.port",
		8080,
		"Metrics HTTP server listening port.",
	)

	cmd.Flags().IntVar(
		&cfg.Metrics.NetworkMetricsTick,
		"metrics.networkMetricsTick",
		60,
		"Network metrics check tick in seconds.",
	)

	cmd.Flags().IntVar(
		&cfg.Metrics.EthereumMetricsTick,
		"metrics.ethereumMetricsTick",
		60,
		"Ethereum metrics check tick in seconds.",
	)
}

// Initialize flags for Diagnostics configuration.
func initDiagnosticsFlags(cmd *cobra.Command, cfg *config.Config) {
	cmd.Flags().IntVar(
		&cfg.Diagnostics.Port,
		"diagnostics.port",
		8081,
		"Diagnostics HTTP server listening port.",
	)
}

// Initialize flags for Developer configuration.
func initDeveloperFlags(command *cobra.Command) {
	initContractAddressFlag := func(contractName string, defaultAddress string) {
		command.Flags().String(
			config.GetDeveloperContractAddressKey(contractName),
			defaultAddress,
			fmt.Sprintf(
				"Address of the %s smart contract",
				contractName,
			),
		)
	}

	initContractAddressFlag(
		chainEthereum.RandomBeaconContractName,
		ethereumBeacon.RandomBeaconAddress,
	)

	initContractAddressFlag(
		chainEthereum.TokenStakingContractName,
		ethereumThreshold.TokenStakingAddress,
	)

	initContractAddressFlag(
		chainEthereum.WalletRegistryContractName,
		ethereumEcdsa.WalletRegistryAddress,
	)
}
