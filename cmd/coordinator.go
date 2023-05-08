package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/keep-network/keep-core/config"
	"github.com/keep-network/keep-core/pkg/bitcoin/electrum"
	"github.com/keep-network/keep-core/pkg/chain/ethereum"
	"github.com/keep-network/keep-core/pkg/coordinator"
)

var (
	walletFlagName       = "wallet"
	hideSweptFlagName    = "hide-swept"
	sortByAmountFlagName = "sort-amount"
	headFlagName         = "head"
	tailFlagName         = "tail"
)

// CoordinatorCommand contains the definition of tBTC Wallet Coordinator tools.
var CoordinatorCommand = &cobra.Command{
	Use:              "coordinator",
	Short:            "tBTC Wallet Coordinator Tools",
	Long:             "The tool exposes commands for interactions with tBTC wallets.",
	TraverseChildren: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := clientConfig.ReadConfig(
			configFilePath,
			cmd.Flags(),
			config.General, config.Ethereum, config.BitcoinElectrum,
		); err != nil {
			logger.Fatalf("error reading config: %v", err)
		}
	},
}

var depositsCommand = cobra.Command{
	Use:              "deposits",
	Short:            "get list of deposits",
	Long:             "Gets tBTC deposits details from the chain and prints them.",
	TraverseChildren: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		wallet, err := cmd.Flags().GetString(walletFlagName)
		if err != nil {
			return fmt.Errorf("failed to find wallet flag: %v", err)
		}

		hideSwept, err := cmd.Flags().GetBool(hideSweptFlagName)
		if err != nil {
			return fmt.Errorf("failed to find show swept flag: %v", err)
		}

		sortByAmount, err := cmd.Flags().GetBool(sortByAmountFlagName)
		if err != nil {
			return fmt.Errorf("failed to find show swept flag: %v", err)
		}

		head, err := cmd.Flags().GetInt(headFlagName)
		if err != nil {
			return fmt.Errorf("failed to find head flag: %v", err)
		}

		tail, err := cmd.Flags().GetInt(tailFlagName)
		if err != nil {
			return fmt.Errorf("failed to find tail flag: %v", err)
		}

		_, tbtcChain, _, _, _, err := ethereum.Connect(ctx, clientConfig.Ethereum)
		if err != nil {
			return fmt.Errorf(
				"could not connect to Ethereum chain: [%v]",
				err,
			)
		}

		btcChain, err := electrum.Connect(ctx, clientConfig.Bitcoin.Electrum)
		if err != nil {
			return fmt.Errorf("could not connect to Electrum chain: [%v]", err)
		}

		return coordinator.ListDeposits(
			tbtcChain,
			btcChain,
			wallet,
			hideSwept,
			sortByAmount,
			head,
			tail,
		)
	},
}

func init() {
	initFlags(
		CoordinatorCommand,
		&configFilePath,
		clientConfig,
		config.General, config.Ethereum, config.BitcoinElectrum,
	)

	// Coordinator Command
	CoordinatorCommand.PersistentFlags().String(
		walletFlagName,
		"",
		"wallet public key hash",
	)

	// Deposits Subcommand
	depositsCommand.Flags().Bool(
		hideSweptFlagName,
		false,
		"hide swept deposits",
	)

	depositsCommand.Flags().Bool(
		sortByAmountFlagName,
		false,
		"sort by deposit amount",
	)

	depositsCommand.Flags().Int(
		headFlagName,
		0,
		"get head of deposits",
	)

	depositsCommand.Flags().Int(
		tailFlagName,
		0,
		"get tail of deposits",
	)

	CoordinatorCommand.AddCommand(&depositsCommand)
}
