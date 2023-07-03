package cmd

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/keep-network/keep-core/config"
	"github.com/keep-network/keep-core/internal/hexutils"
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/bitcoin/electrum"
	"github.com/keep-network/keep-core/pkg/chain/ethereum"
	"github.com/keep-network/keep-core/pkg/maintainer/spv"
	walletmtr "github.com/keep-network/keep-core/pkg/maintainer/wallet"
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

// MaintainerCliCommand contains the definition of tools associated with maintainers
// module.
var MaintainerCliCommand = &cobra.Command{
	Use:              "maintainer-cli",
	Short:            "Maintainer CLI Tools",
	Long:             "The tool exposes commands for tools associated with maintainers.",
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

		var walletPublicKeyHash [20]byte
		if len(wallet) > 0 {
			var err error
			walletPublicKeyHash, err = newWalletPublicKeyHash(wallet)
			if err != nil {
				return fmt.Errorf(
					"failed to extract wallet public key hash: %v",
					err,
				)
			}
		}

		deposits, err := walletmtr.FindDeposits(
			tbtcChain,
			btcChain,
			walletPublicKeyHash,
			head,
			hideSwept,
			false,
		)
		if err != nil {
			return fmt.Errorf(
				"failed to get deposits: [%w]",
				err,
			)
		}

		if len(deposits) == 0 {
			return fmt.Errorf("no deposits found")
		}

		if err := printDepositsTable(deposits); err != nil {
			return fmt.Errorf("failed to print deposits table: %v", err)
		}

		return nil
	},
}

func printDepositsTable(deposits []*walletmtr.Deposit) error {
	w := tabwriter.NewWriter(os.Stdout, 2, 4, 1, ' ', tabwriter.AlignRight)
	fmt.Fprintf(w, "index\twallet\tvalue (BTC)\tdeposit key\trevealed deposit data\tconfirmations\tswept\t\n")

	for i, deposit := range deposits {
		fmt.Fprintf(w, "%d\t%s\t%.5f\t%s\t%s\t%d\t%t\t\n",
			i,
			hexutils.Encode(deposit.WalletPublicKeyHash[:]),
			deposit.AmountBtc,
			deposit.DepositKey,
			fmt.Sprintf(
				"%s:%d:%d",
				deposit.FundingTxHash.Hex(bitcoin.ReversedByteOrder),
				deposit.FundingOutputIndex,
				deposit.RevealBlock,
			),
			deposit.Confirmations,
			deposit.IsSwept,
		)
	}

	if err := w.Flush(); err != nil {
		return fmt.Errorf("failed to flush the writer: %v", err)
	}

	return nil
}

var proposeDepositsSweepCommand = cobra.Command{
	Use:              "propose-deposits-sweep",
	Short:            "propose deposits sweep",
	Long:             proposeDepositsSweepCommandDescription,
	TraverseChildren: true,
	Args: func(cmd *cobra.Command, args []string) error {
		for i, arg := range args {
			if err := validateDepositReferenceString(arg); err != nil {
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
			return fmt.Errorf("failed to find deposit sweep max size flag: %v", err)
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

		var deposits []*walletmtr.DepositReference
		if len(args) > 0 {
			deposits, err = parseDepositsReferences(args)
			if err != nil {
				return fmt.Errorf("failed extract wallet public key hash: %v", err)
			}
		} else {
			walletPublicKeyHash, deposits, err = walletmtr.FindDepositsToSweep(
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

		return walletmtr.ProposeDepositsSweep(
			tbtcChain,
			btcChain,
			walletPublicKeyHash,
			fee,
			deposits,
			dryRun,
		)
	},
}

var (
	depositsFormatDescription = "Deposits details should be provided as strings containing: \n" +
		" - bitcoin transaction hash (unprefixed bitcoin transaction hash in reverse (RPC) order), \n" +
		" - bitcoin transaction output index, \n" +
		" - ethereum block number when the deposit was revealed to the chain. \n" +
		"The properties should be separated by semicolons, in the following format: \n" +
		depositReferenceFormatPattern + "\n" +
		"e.g. bd99d1d0a61fd104925d9b7ac997958aa8af570418b3fde091f7bfc561608865:1:8392394"

	depositReferenceFormatPattern = "<unprefixed bitcoin transaction hash>:" +
		"<bitcoin transaction output index>:" +
		"<ethereum reveal block number>"

	depositReferenceFormatRegexp = regexp.MustCompile(`^([[:xdigit:]]+):(\d+):(\d+)$`)

	proposeDepositsSweepCommandDescription = "Submits a deposits sweep proposal " +
		"to the chain. Expects --wallet and --fee flags along with deposits to " +
		"sweep provided as arguments.\n" +
		depositsFormatDescription
)

// parseDepositsReferences decodes a list of deposits references.
func parseDepositsReferences(
	depositsRefsStings []string,
) ([]*walletmtr.DepositReference, error) {
	depositsRefs := make([]*walletmtr.DepositReference, len(depositsRefsStings))

	for i, depositRefString := range depositsRefsStings {
		matched := depositReferenceFormatRegexp.FindStringSubmatch(depositRefString)
		// Check if number of resolved entries match expected number of groups
		// for the given regexp.
		if len(matched) != 4 {
			return nil, fmt.Errorf(
				"failed to parse deposit: [%s]",
				depositRefString,
			)
		}

		txHash, err := bitcoin.NewHashFromString(matched[1], bitcoin.ReversedByteOrder)
		if err != nil {
			return nil, fmt.Errorf(
				"invalid bitcoin transaction hash [%s]: %v",
				matched[1],
				err,
			)
		}

		outputIndex, err := strconv.ParseInt(matched[2], 10, 32)
		if err != nil {
			return nil, fmt.Errorf(
				"invalid bitcoin transaction output index [%s]: %v",
				matched[2],
				err,
			)
		}

		revealBlock, err := strconv.ParseUint(matched[3], 10, 32)
		if err != nil {
			return nil, fmt.Errorf(
				"invalid reveal block number [%s]: %v",
				matched[3],
				err,
			)
		}

		depositsRefs[i] = &walletmtr.DepositReference{
			FundingTxHash:      txHash,
			FundingOutputIndex: uint32(outputIndex),
			RevealBlock:        revealBlock,
		}
	}

	return depositsRefs, nil
}

// validateDepositReferenceString validates format of the string containing a
// deposit reference.
func validateDepositReferenceString(depositRefString string) error {
	if !depositReferenceFormatRegexp.MatchString(depositRefString) {
		return fmt.Errorf(
			"[%s] doesn't match pattern: %s",
			depositRefString,
			depositReferenceFormatPattern,
		)
	}
	return nil
}

var proposeRedemptionCommand = cobra.Command{
	Use:              "propose-redemption",
	Short:            "propose redemption",
	Long:             "Submits a redemption proposal to the chain.",
	TraverseChildren: true,
	RunE: func(cmd *cobra.Command, args []string) error {
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
			return fmt.Errorf("failed to find redemption max size flag: %v", err)
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

		btcChain, err := electrum.Connect(cmd.Context(), clientConfig.Bitcoin.Electrum)
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

		if redemptionMaxSize == 0 {
			redemptionMaxSize, err = tbtcChain.GetRedemptionMaxSize()
			if err != nil {
				return fmt.Errorf("failed to get redemption max size: [%v]", err)
			}
		}

		walletPublicKeyHash, redemptions, err := walletmtr.FindPendingRedemptions(
			tbtcChain,
			walletPublicKeyHash,
			redemptionMaxSize,
		)
		if err != nil {
			return fmt.Errorf("failed to prepare redemption proposal: %v", err)
		}

		if len(redemptions) > int(redemptionMaxSize) {
			return fmt.Errorf(
				"redemptions number [%d] is greater than redemptions max size [%d]",
				len(redemptions),
				redemptionMaxSize,
			)
		}

		return walletmtr.ProposeRedemption(
			tbtcChain,
			btcChain,
			walletPublicKeyHash,
			fee,
			redemptions,
			dryRun,
		)
	},
}

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

		fees, err := walletmtr.EstimateDepositsSweepFee(
			tbtcChain,
			btcChain,
			depositsCount,
		)
		if err != nil {
			return fmt.Errorf("cannot estimate deposits sweep fee: [%v]", err)
		}

		err = printDepositsSweepFeeTable(fees)
		if err != nil {
			return fmt.Errorf("cannot print fees table: [%v]", err)
		}

		return nil
	},
}

// printDepositsSweepFeeTable prints estimated fees for specific deposits counts
// to the standard output. For example:
//
// ---------------------------------------------
// deposits count total fee (satoshis) sat/vbyte
//	            1                  201         1
//	            2                  292         1
//	            3                  384         1
// ---------------------------------------------
func printDepositsSweepFeeTable(
	fees map[int]struct {
		TotalFee       int64
		SatPerVByteFee int64
	},
) error {
	writer := tabwriter.NewWriter(
		os.Stdout,
		2,
		4,
		1,
		' ',
		tabwriter.AlignRight,
	)

	_, err := fmt.Fprintf(writer, "deposits count\ttotal fee (satoshis)\tsat/vbyte\t\n")
	if err != nil {
		return err
	}

	var depositsCountKeys []int
	for depositsCountKey := range fees {
		depositsCountKeys = append(depositsCountKeys, depositsCountKey)
	}

	sort.Slice(depositsCountKeys, func(i, j int) bool {
		return depositsCountKeys[i] < depositsCountKeys[j]
	})

	for _, depositsCountKey := range depositsCountKeys {
		_, err := fmt.Fprintf(
			writer,
			"%v\t%v\t%v\t\n",
			depositsCountKey,
			fees[depositsCountKey].TotalFee,
			fees[depositsCountKey].SatPerVByteFee,
		)
		if err != nil {
			return err
		}
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush the writer: %v", err)
	}

	return nil
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
		MaintainerCliCommand,
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

	MaintainerCliCommand.AddCommand(&listDepositsCommand)

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

	MaintainerCliCommand.AddCommand(&proposeDepositsSweepCommand)

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

	MaintainerCliCommand.AddCommand(&proposeRedemptionCommand)

	// Estimate Deposits Sweep Fee Subcommand.
	estimateDepositsSweepFeeCommand.Flags().Int(
		depositsCountFlagName,
		0,
		"get estimation for a specific count of input deposits",
	)

	MaintainerCliCommand.AddCommand(&estimateDepositsSweepFeeCommand)

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

	MaintainerCliCommand.AddCommand(&submitDepositSweepProofCommand)
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
