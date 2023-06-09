package cmd

import (
	"fmt"

	"github.com/btcsuite/btcd/txscript"
	"github.com/spf13/cobra"

	"github.com/ethereum/go-ethereum/common"
	"github.com/keep-network/keep-core/config"
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/bitcoin/electrum"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/chain/ethereum"
	"github.com/keep-network/keep-core/pkg/coordinator"
)

var (
	// listDepositsCommand:
	// proposeDepositsSweepCommand:
	walletFlagName = "wallet"

	// listDepositsCommand:
	hideSweptFlagName    = "hide-swept"
	sortByAmountFlagName = "sort-amount"
	headFlagName         = "head"
	tailFlagName         = "tail"

	// proposeDepositsSweepCommand:
	feeFlagName    = "fee"
	dryRunFlagName = "dry-run"

	// estimateDepositsSweepFeeCommand:
	depositsCountFlagName = "deposits-count"

	// submitDepositSweepProofCommand:
	transactionHashFlagName = "transaction-hash"
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

		sortByAmount, err := cmd.Flags().GetBool(sortByAmountFlagName)
		if err != nil {
			return fmt.Errorf("failed to find sort by amount flag: %v", err)
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

		return coordinator.ProposeDepositsSweep(tbtcChain, btcChain, wallet, fee, args, dryRun)
	},
}

var proposeDepositsSweepCommandDescription = `Submits a deposits sweep proposal to
the chain.
Expects --wallet and --fee flags along with deposits to sweep provided
as arguments.

` + coordinator.DepositsFormatDescription

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
	Long:             submitDepositSweepProofCommandDescription,
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
			return fmt.Errorf("could not connect to Electrum chain: %v", err)
		}

		transactionHashFlag, err := cmd.Flags().GetString(transactionHashFlagName)
		if err != nil {
			return fmt.Errorf("failed to find transaction hash flag: %v", err)
		}

		transactionHash, err := bitcoin.NewHashFromString(
			transactionHashFlag,
			bitcoin.ReversedByteOrder,
		)
		if err != nil {
			return fmt.Errorf("failed to parse transaction hash flag: %v", err)
		}

		// TODO: Due to how the Bridge contract checks the accumulated difficulty
		//       it may be necessary to increment the required confirmations
		//       by one.
		requiredConfirmations, err := tbtcChain.TxProofDifficultyFactor()
		if err != nil {
			return fmt.Errorf(
				"failed to get transaction proof difficulty factor: %v",
				err,
			)
		}

		transaction, proof, err := bitcoin.AssembleSpvProof(
			transactionHash,
			uint(requiredConfirmations.Uint64()),
			btcChain,
		)
		if err != nil {
			return fmt.Errorf("failed to assemble transaction spv proof: %v", err)
		}

		mainUTXO, vault, err := parseTransactionInputs(
			btcChain,
			*tbtcChain,
			transaction)
		if err != nil {
			return fmt.Errorf("error while parsing transaction inputs: %v", err)
		}

		if err := tbtcChain.SubmitDepositSweepProof(
			proof,
			transaction,
			mainUTXO,
			vault,
		); err != nil {
			return fmt.Errorf("failed to submit deposit sweep proof: %v", err)
		}

		return nil
	},
}

func convertVaultAddress(vault *chain.Address) common.Address {
	if vault == nil {
		return common.Address{}
	}

	return common.HexToAddress(string(*vault))
}

func parseTransactionInputs(
	btcChain bitcoin.Chain,
	tbtcChain ethereum.TbtcChain,
	transaction *bitcoin.Transaction,
) (
	bitcoin.UnspentTransactionOutput,
	common.Address,
	error,
) {
	// Represents the main UTXO of the deposit sweep transaction. Nil if there
	// was no main UTXO.
	var mainUTXO *bitcoin.UnspentTransactionOutput = nil

	// Stores the vault address of the deposits. Each deposit should have the
	// same value of vault. The zero-filled value indicates there was no vault
	// value set for the deposits.
	var vault = common.Address{}

	// This flag checks if at least one deposit input has been found during
	// deposit processing.
	var depositAlreadyProcessed = false

	for _, input := range transaction.Inputs {
		outpointTransactionHash := input.Outpoint.TransactionHash
		outpointIndex := input.Outpoint.OutputIndex

		previousTransaction, err := btcChain.GetTransaction(
			outpointTransactionHash,
		)
		if err != nil {
			return bitcoin.UnspentTransactionOutput{}, common.Address{}, fmt.Errorf(
				"failed to get previous transaction: %v",
				err,
			)
		}

		publicKeyScript := previousTransaction.Outputs[outpointIndex].PublicKeyScript
		value := previousTransaction.Outputs[outpointIndex].Value
		scriptClass := txscript.GetScriptClass(publicKeyScript)

		if scriptClass == txscript.PubKeyHashTy ||
			scriptClass == txscript.WitnessV0PubKeyHashTy {
			if mainUTXO == nil {
				mainUTXO = &bitcoin.UnspentTransactionOutput{
					Outpoint: &bitcoin.TransactionOutpoint{
						TransactionHash: outpointTransactionHash,
						OutputIndex:     outpointIndex,
					},
					Value: value,
				}
			} else {
				return bitcoin.UnspentTransactionOutput{}, common.Address{}, fmt.Errorf(
					"deposit sweep transaction has incorrect structure",
				)
			}
		} else if scriptClass == txscript.ScriptHashTy ||
			scriptClass == txscript.WitnessV0ScriptHashTy {

			deposit, err := tbtcChain.GetDepositRequest(
				outpointTransactionHash,
				outpointIndex,
			)
			if err != nil {
				return bitcoin.UnspentTransactionOutput{}, common.Address{}, fmt.Errorf(
					"failed to get deposit request: %v",
					err,
				)
			}

			if depositAlreadyProcessed {
				if vault != convertVaultAddress(deposit.Vault) {
					return bitcoin.UnspentTransactionOutput{}, common.Address{}, fmt.Errorf(
						"swept deposits have different vaults",
					)
				}
			} else {
				// The first deposit input has been encountered, save the vault
				// and require subsequent deposit inputs to have the same vault.
				vault = convertVaultAddress(deposit.Vault)
				depositAlreadyProcessed = true
			}

		} else {
			return bitcoin.UnspentTransactionOutput{}, common.Address{}, fmt.Errorf(
				"deposit sweep transaction has incorrect input types",
			)
		}
	}

	// If none of the input was main UTXO, return zero-filled main UTXO.
	if mainUTXO == nil {
		mainUTXO = &bitcoin.UnspentTransactionOutput{
			Outpoint: &bitcoin.TransactionOutpoint{
				TransactionHash: bitcoin.Hash{},
				OutputIndex:     0,
			},
			Value: 0,
		}
	}

	return *mainUTXO, vault, nil
}

var submitDepositSweepProofCommandDescription = "Submits deposit sweep proof " +
	"to the Bridge contract"

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
		"fee for the entire bitcoin transaction (satoshi)",
	)

	proposeDepositsSweepCommand.Flags().Bool(
		dryRunFlagName,
		false,
		"don't submit a proposal to the chain",
	)

	CoordinatorCommand.AddCommand(&proposeDepositsSweepCommand)

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
			"be the same as in Bitcoin explorers)",
	)

	CoordinatorCommand.AddCommand(&submitDepositSweepProofCommand)
}
