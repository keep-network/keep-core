package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/keep-network/keep-core/config"
	"github.com/keep-network/keep-core/internal/hexutils"
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/bitcoin/electrum"
	"github.com/keep-network/keep-core/pkg/chain/ethereum"
	"github.com/keep-network/keep-core/pkg/coordinator"
	"github.com/keep-network/keep-core/pkg/maintainer/spv"
)

var (
	// listDepositsCommand:
	// proposeDepositsSweepCommand:
	walletFlagName = "wallet"

	// listDepositsCommand:
	hideSweptFlagName = "hide-swept"
	headFlagName      = "head"

	// proposeDepositsSweepCommand:
	// proposeRedemptionsCommand:
	feeFlagName    = "fee"
	dryRunFlagName = "dry-run"

	// proposeDepositsSweepCommand:
	depositSweepMaxSizeFlagName = "deposit-sweep-max-size"

	// proposeRedemptionsCommand:
	redemptionMaxSizeFlagName = "redemption-max-size"

	// estimateDepositsSweepFeeCommand:
	depositsCountFlagName = "deposits-count"

	// submitDepositSweepProofCommand:
	transactionHashFlagName = "transaction-hash"
	confirmationsFlagName   = "confirmations"
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
			return fmt.Errorf("failed to find hide swept flag: %v", err)
		}

		head, err := cmd.Flags().GetInt(headFlagName)
		if err != nil {
			return fmt.Errorf("failed to find head flag: %v", err)
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

		var walletPublicKeyHash [20]byte
		if len(wallet) > 0 {
			var err error
			walletPublicKeyHash, err = newWalletPublicKeyHash(wallet)
			if err != nil {
				return fmt.Errorf("failed to extract wallet public key hash: %v", err)
			}
		}

		return coordinator.ListDeposits(
			tbtcChain,
			btcChain,
			walletPublicKeyHash,
			head,
			hideSwept,
		)
	},
}

var proposeDepositsSweepCommand = cobra.Command{
	Use:              "propose-deposits-sweep",
	Short:            "propose deposits sweep",
	Long:             proposeDepositsSweepCommandDescription,
	TraverseChildren: true,
	Args: func(cmd *cobra.Command, args []string) error {
		for i, arg := range args {
			if err := coordinator.ValidateDepositString(arg); err != nil {
				return fmt.Errorf(
					"argument [%d] failed validation: %v",
					i,
					err,
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

		depositSweepMaxSize, err := cmd.Flags().GetUint16(depositSweepMaxSizeFlagName)
		if err != nil {
			return fmt.Errorf("failed to find fee flag: %v", err)
		}

		dryRun, err := cmd.Flags().GetBool(dryRunFlagName)
		if err != nil {
			return fmt.Errorf("failed to find dry run flag: %v", err)
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

		var walletPublicKeyHash [20]byte
		if len(wallet) > 0 {
			var err error
			walletPublicKeyHash, err = newWalletPublicKeyHash(wallet)
			if err != nil {
				return fmt.Errorf("failed extract wallet public key hash: %v", err)
			}
		}

		if depositSweepMaxSize == 0 {
			depositSweepMaxSize, err = tbtcChain.GetDepositSweepMaxSize()
			if err != nil {
				return fmt.Errorf("failed to get deposit sweep max size: [%v]", err)
			}
		}

		var deposits []*coordinator.DepositSweepDetails
		if len(args) > 0 {
			deposits, err = coordinator.ParseDepositsToSweep(args)
			if err != nil {
				return fmt.Errorf("failed extract wallet public key hash: %v", err)
			}
		} else {
			walletPublicKeyHash, deposits, err = coordinator.FindDepositsToSweep(
				tbtcChain,
				btcChain,
				walletPublicKeyHash,
				depositSweepMaxSize,
			)
			if err != nil {
				return fmt.Errorf("failed to prepare deposits sweep proposal: %v", err)
			}
		}

		if len(deposits) > int(depositSweepMaxSize) {
			return fmt.Errorf(
				"deposits number [%d] is greater than deposit sweep max size [%d]",
				len(deposits),
				depositSweepMaxSize,
			)
		}

		return coordinator.ProposeDepositsSweep(
			tbtcChain,
			btcChain,
			walletPublicKeyHash,
			fee,
			deposits,
			dryRun,
		)
	},
}

var proposeDepositsSweepCommandDescription = `Submits a deposits sweep proposal to
the chain.
Expects --wallet and --fee flags along with deposits to sweep provided
as arguments.

` + coordinator.DepositsFormatDescription

var proposeRedemptionCommand = cobra.Command{
	Use:              "propose-redemption",
	Short:            "propose redemption",
	Long:             proposeRedemptionCommandDescription,
	TraverseChildren: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		// ctx := cmd.Context()

		wallet, err := cmd.Flags().GetString(walletFlagName)
		if err != nil {
			return fmt.Errorf("failed to find wallet flag: %v", err)
		}

		fee, err := cmd.Flags().GetInt64(feeFlagName)
		if err != nil {
			return fmt.Errorf("failed to find fee flag: %v", err)
		}

		redemptionMaxSize, err := cmd.Flags().GetUint16(redemptionMaxSizeFlagName)
		if err != nil {
			return fmt.Errorf("failed to find fee flag: %v", err)
		}

		dryRun, err := cmd.Flags().GetBool(dryRunFlagName)
		if err != nil {
			return fmt.Errorf("failed to find dry run flag: %v", err)
		}

		_, tbtcChain, _, _, _, err := ethereum.Connect(cmd.Context(), clientConfig.Ethereum)
		if err != nil {
			return fmt.Errorf(
				"could not connect to Ethereum chain: [%v]",
				err,
			)
		}

		var walletPublicKeyHash [20]byte
		if len(wallet) > 0 {
			var err error
			walletPublicKeyHash, err = newWalletPublicKeyHash(wallet)
			if err != nil {
				return fmt.Errorf("failed extract wallet public key hash: %v", err)
			}
		}

		if redemptionMaxSize == 0 {
			redemptionMaxSize, err = tbtcChain.GetRedemptionMaxSize()
			if err != nil {
				return fmt.Errorf("failed to get deposit sweep max size: [%v]", err)
			}
		}

		walletPublicKeyHash, redemptions, err := coordinator.FindPendingRedemptions(
			tbtcChain,
			walletPublicKeyHash,
			redemptionMaxSize,
		)
		if err != nil {
			return fmt.Errorf("failed to prepare deposits sweep proposal: %v", err)
		}

		if len(redemptions) > int(redemptionMaxSize) {
			return fmt.Errorf(
				"redemptions number [%d] is greater than redemptions max size [%d]",
				len(redemptions),
				redemptionMaxSize,
			)
		}

		return coordinator.ProposeRedemption(
			tbtcChain,
			walletPublicKeyHash,
			fee,
			redemptions,
			dryRun,
		)
	},
}

var proposeRedemptionCommandDescription = `Submits a redemption proposal to the chain.`

var estimateDepositsSweepFeeCommand = cobra.Command{
	Use:              "estimate-deposits-sweep-fee",
	Short:            "estimates deposits sweep fee",
	Long:             estimateDepositsSweepFeeCommandDescription,
	TraverseChildren: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		depositsCount, err := cmd.Flags().GetInt(depositsCountFlagName)
		if err != nil {
			return fmt.Errorf("failed to find deposits count flag: %v", err)
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

		return coordinator.EstimateDepositsSweepFee(
			tbtcChain,
			btcChain,
			depositsCount,
		)
	},
}

var estimateDepositsSweepFeeCommandDescription = "Estimates the satoshi " +
	"fee for the entire Bitcoin deposits sweep transaction, based on " +
	"the number of input deposits. By default, provides estimations for " +
	"transactions containing a various number of input deposits, from 1 up " +
	"to the maximum count allowed by the WalletCoordinator contract. " +
	"The --deposits-count flag can be used to obtain a fee estimation for " +
	"a Bitcoin sweep transaction containing a specific count of input " +
	"deposits. All estimations assume the wallet main UTXO is used as one " +
	"of the transaction's input so the estimation may be overpriced for " +
	"the very first sweep transaction of each wallet. Estimations also " +
	"assume only P2WSH deposits are part of the transaction so the " +
	"estimation may be underpriced if the actual transaction contains " +
	"legacy P2SH deposits. If the estimated fee exceeds the maximum fee " +
	"allowed by the Bridge contract, the maximum fee is returned as result"

var submitDepositSweepProofCommand = cobra.Command{
	Use:              "submit-deposit-sweep-proof",
	Short:            "submit deposit sweep proof",
	Long:             "Submits deposit sweep proof to the Bridge contract",
	TraverseChildren: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		_, tbtcChain, _, _, _, err := ethereum.Connect(
			ctx,
			clientConfig.Ethereum,
		)
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

		transactionHashFlag, err := cmd.Flags().GetString(transactionHashFlagName)
		if err != nil {
			return fmt.Errorf("failed to find transaction hash flag: [%v]", err)
		}

		transactionHash, err := bitcoin.NewHashFromString(
			transactionHashFlag,
			bitcoin.ReversedByteOrder,
		)
		if err != nil {
			return fmt.Errorf(
				"failed to parse transaction hash flag: [%v]",
				err,
			)
		}

		// Allow the caller to request a specific number of confirmations.
		// The Bridge calculates the required difficulty of a chain of block
		// headers by multiplying the difficulty of the first block header by
		// the difficulty factor. If the block headers happen to span the
		// Bitcoin epoch difficulty change and there is a drop of difficulty
		// between the epochs, the sum of difficulties from the headers chain
		// may be too low. Allowing the caller to specify a greater number of
		// confirmations will ensure the transaction can be proven.
		requiredConfirmations, err := cmd.Flags().GetUint(confirmationsFlagName)
		if err != nil {
			return fmt.Errorf("failed to get confirmations flag: [%v]", err)
		}

		logger.Infof(
			"Submitting deposit sweep proof for transaction [%s]",
			transactionHashFlag,
		)

		if err = spv.SubmitDepositSweepProof(
			transactionHash,
			requiredConfirmations,
			btcChain,
			tbtcChain,
		); err != nil {
			return fmt.Errorf("failed to submit deposit sweep proof [%v]", err)
		}

		logger.Infof(
			"Successfully submitted deposit sweep proof for transaction: [%s]",
			transactionHashFlag,
		)

		return nil
	},
}

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

	listDepositsCommand.Flags().Int(
		headFlagName,
		0,
		"get head of deposits",
	)

	CoordinatorCommand.AddCommand(&listDepositsCommand)

	// Propose Deposits Sweep Subcommand
	proposeDepositsSweepCommand.Flags().String(
		walletFlagName,
		"",
		"wallet public key hash",
	)

	proposeDepositsSweepCommand.Flags().Int64(
		feeFlagName,
		0,
		"fee for the entire bitcoin transaction (satoshi)",
	)

	proposeDepositsSweepCommand.Flags().Uint16(
		depositSweepMaxSizeFlagName,
		0,
		"maximum count of deposits that can be swept within a single sweep",
	)

	proposeDepositsSweepCommand.Flags().Bool(
		dryRunFlagName,
		false,
		"don't submit a proposal to the chain",
	)

	CoordinatorCommand.AddCommand(&proposeDepositsSweepCommand)

	// Propose Redemptions Subcommand
	proposeRedemptionCommand.Flags().String(
		walletFlagName,
		"",
		"wallet public key hash",
	)

	proposeRedemptionCommand.Flags().Int64(
		feeFlagName,
		0,
		"fee for the entire bitcoin transaction (satoshi)",
	)

	proposeRedemptionCommand.Flags().Uint16(
		redemptionMaxSizeFlagName,
		0,
		"maximum count of deposits that can be redeemed within a single redemption",
	)

	proposeRedemptionCommand.Flags().Bool(
		dryRunFlagName,
		false,
		"don't submit a proposal to the chain",
	)

	CoordinatorCommand.AddCommand(&proposeRedemptionCommand)

	// Estimate Deposits Sweep Fee Subcommand.
	estimateDepositsSweepFeeCommand.Flags().Int(
		depositsCountFlagName,
		0,
		"get estimation for a specific count of input deposits",
	)

	CoordinatorCommand.AddCommand(&estimateDepositsSweepFeeCommand)

	// Submit Deposit Sweep Proof Subcommand.

	submitDepositSweepProofCommand.Flags().String(
		transactionHashFlagName,
		"",
		"transaction hash the proof will be prepared for (the format should "+
			"be the same as in Bitcoin explorers).",
	)

	if err := submitDepositSweepProofCommand.MarkFlagRequired(
		transactionHashFlagName,
	); err != nil {
		logger.Fatalf("failed to mark flag required: [%v]", err)
	}

	submitDepositSweepProofCommand.Flags().Uint(
		confirmationsFlagName,
		0,
		"(optional) number of confirmations that will be provided in the proof. "+
			"This is an optional parameter that can be used in a rare event when "+
			"more confirmations are required to perform a successful proof "+
			"validation. If this parameter is not provided, the default value, "+
			"retrieved from the Bridge will be used.",
	)

	CoordinatorCommand.AddCommand(&submitDepositSweepProofCommand)
}

func newWalletPublicKeyHash(str string) ([20]byte, error) {
	var result [20]byte

	walletHex, err := hexutils.Decode(str)
	if err != nil {
		return result, err
	}

	if len(walletHex) != 20 {
		return result, fmt.Errorf("invalid bytes length: [%d], expected: [%d]", len(walletHex), 20)
	}

	copy(result[:], walletHex)

	return result, nil
}
