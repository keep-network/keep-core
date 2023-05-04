package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	walletcmd "github.com/keep-network/keep-core/cmd/wallet"
	"github.com/keep-network/keep-core/config"
	"github.com/keep-network/keep-core/pkg/bitcoin/electrum"
	"github.com/keep-network/keep-core/pkg/chain/ethereum"
)

var (
	walletFlagName       = "wallet"
	hideSweptFlagName    = "hide-swept"
	sortByAmountFlagName = "sort-amount"
	feeFlagName          = "fee"
)

// WalletCommand contains the definition of tBTC wallets tools.
var WalletCommand = &cobra.Command{
	Use:              "wallet",
	Short:            "tBTC wallets tools",
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
	Short:            "get deposits",
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

		_, tbtcChain, _, _, _, err := ethereum.Connect(ctx, clientConfig.Ethereum)
		if err != nil {
			return fmt.Errorf(
				"could not connect to Bitcoin difficulty chain: [%v]",
				err,
			)
		}

		return walletcmd.ListDeposits(tbtcChain, wallet, hideSwept, sortByAmount)
	},
}

var sweepCommand = cobra.Command{
	Use:              "sweep",
	Short:            "propose deposits sweep",
	Long:             sweepCommandDescription,
	TraverseChildren: true,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
			return err
		}

		for _, arg := range args {
			if !walletcmd.BitcoinTxRegexp.MatchString(arg) {
				return fmt.Errorf(
					"argument [%s] doesn't match pattern: <unprefixed transaction hash>:<output index>",
					arg,
				)
			}
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		wallet, err := cmd.Flags().GetString(walletFlagName)
		if err != nil {
			return fmt.Errorf("failed to find wallet flag: %v", err)
		}

		fee, err := cmd.Flags().GetInt64(feeFlagName)
		if err != nil {
			return fmt.Errorf("failed to find fee flag: %v", err)
		}

		_, tbtcChain, _, _, _, err := ethereum.Connect(cmd.Context(), clientConfig.Ethereum)
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
		return walletcmd.ProposeSweep(tbtcChain, btcChain, wallet, fee, args)
	},
}

var sweepCommandDescription = `Submits a deposits sweep proposal to the chain.
Expects --wallet and --fee flags along with bitcoin transactions to sweep provided
as arguments.

Bitcoin transactions should be provided in the following format:
<unprefixed transaction hash>:<output index>
e.g. bd99d1d0a61fd104925d9b7ac997958aa8af570418b3fde091f7bfc561608865:1
`

func init() {
	initFlags(
		WalletCommand,
		&configFilePath,
		clientConfig,
		config.General, config.Ethereum, config.BitcoinElectrum,
	)

	// Wallet Command
	WalletCommand.PersistentFlags().String(
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

	WalletCommand.AddCommand(&depositsCommand)

	// Sweep Subcommand
	if err := sweepCommand.MarkFlagRequired(walletFlagName); err != nil {
		logger.Panicf("failed to mark wallet flag as required: %v", err)
	}

	sweepCommand.Flags().Int64(
		feeFlagName,
		0,
		"fee for bitcoin transaction",
	)

	WalletCommand.AddCommand(&sweepCommand)
}
