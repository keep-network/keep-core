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
	feeFlagName          = "fee"
	dryRunFlagName       = "dry-run"
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

var listDepositsCommand = cobra.Command{
	Use:              "list-deposits",
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

var proposeDepositsSweepCommand = cobra.Command{
	Use:              "propose-deposits-sweep",
	Short:            "propose deposits sweep",
	Long:             proposeDepositsSweepCommandDescription,
	TraverseChildren: true,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
			return err
		}

		for _, arg := range args {
			if !coordinator.BitcoinTxRegexp.MatchString(arg) {
				return fmt.Errorf(
					"argument [%s] doesn't match pattern: <unprefixed transaction hash>:<output index>:<reveal block number>",
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

		dryRun, err := cmd.Flags().GetBool(dryRunFlagName)
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

		return coordinator.ProposeDepositsSweep(tbtcChain, btcChain, wallet, fee, args, dryRun)
	},
}

var proposeDepositsSweepCommandDescription = `Submits a deposits sweep proposal to 
the chain.
Expects --wallet and --fee flags along with deposits to sweep provided
as arguments.

Deposit details should be provided in the following format:
<unprefixed bitcoin transaction hash>:<bitcoin transaction output index>:<ethereum reveal block number>
e.g. bd99d1d0a61fd104925d9b7ac997958aa8af570418b3fde091f7bfc561608865:1:8392394
`

func init() {
	initFlags(
		CoordinatorCommand,
		&configFilePath,
		clientConfig,
		config.General, config.Ethereum, config.BitcoinElectrum,
	)

	// Deposits Subcommand
	listDepositsCommand.Flags().String(
		walletFlagName,
		"",
		"wallet public key hash",
	)

	listDepositsCommand.Flags().Bool(
		hideSweptFlagName,
		false,
		"hide swept deposits",
	)

	listDepositsCommand.Flags().Bool(
		sortByAmountFlagName,
		false,
		"sort by deposit amount",
	)

	listDepositsCommand.Flags().Int(
		headFlagName,
		0,
		"get head of deposits",
	)

	listDepositsCommand.Flags().Int(
		tailFlagName,
		0,
		"get tail of deposits",
	)

	CoordinatorCommand.AddCommand(&listDepositsCommand)

	// Propose Deposits Sweep Subcommand
	proposeDepositsSweepCommand.Flags().String(
		walletFlagName,
		"",
		"wallet public key hash",
	)

	if err := proposeDepositsSweepCommand.MarkFlagRequired(walletFlagName); err != nil {
		logger.Panicf("failed to mark wallet flag as required: %v", err)
	}

	proposeDepositsSweepCommand.Flags().Int64(
		feeFlagName,
		0,
		"fee for bitcoin transaction",
	)

	proposeDepositsSweepCommand.Flags().Bool(
		dryRunFlagName,
		false,
		"don't submit a proposal to the chain",
	)

	CoordinatorCommand.AddCommand(&proposeDepositsSweepCommand)
}
