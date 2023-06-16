package cmd

import (
	"fmt"

	"github.com/btcsuite/btcd/txscript"
	"github.com/spf13/cobra"

	"github.com/ethereum/go-ethereum/common"
	"github.com/keep-network/keep-core/config"
	"github.com/keep-network/keep-core/internal/hexutils"
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
	hideSweptFlagName = "hide-swept"
	headFlagName      = "head"

	// proposeDepositsSweepCommand:
	feeFlagName                 = "fee"
	depositSweepMaxSizeFlagName = "deposit-sweep-max-size"
	dryRunFlagName              = "dry-run"

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

		logger.Infof(
			"Preparing deposit sweep proof for transaction: %s",
			transactionHashFlag,
		)

		transactionHash, err := bitcoin.NewHashFromString(
			transactionHashFlag,
			bitcoin.ReversedByteOrder,
		)
		if err != nil {
			return fmt.Errorf("failed to parse transaction hash flag: %v", err)
		}

		txProofDifficulty, err := tbtcChain.TxProofDifficultyFactor()
		if err != nil {
			return fmt.Errorf(
				"failed to get transaction proof difficulty factor: %v",
				err,
			)
		}

		// Increase the required confirmations by one. The Bridge calculates
		// the required difficulty of a chain of block headers by multiplying
		// the difficulty of the first block header by the difficulty factor.
		// If the block headers happen to span the Bitcoin epoch difficulty
		// change and there is a drop of difficulty between the epochs, the sum
		// of difficulties from the headers chain may be too low. Adding one
		// more block header will ensure the sum of difficulties is high enough.
		requiredConfirmations := uint(txProofDifficulty.Uint64()) + 1

		transaction, proof, err := bitcoin.AssembleSpvProof(
			transactionHash,
			requiredConfirmations,
			btcChain,
		)
		if err != nil {
			return fmt.Errorf("failed to assemble transaction spv proof: %v", err)
		}

		mainUTXO, vault, err := parseTransactionInputs(
			btcChain,
			*tbtcChain,
			transaction,
		)
		if err != nil {
			return fmt.Errorf("error while parsing transaction inputs: %v", err)
		}

		if err := tbtcChain.SubmitDepositSweepProof(
			transaction,
			proof,
			mainUTXO,
			vault,
		); err != nil {
			return fmt.Errorf("failed to submit deposit sweep proof: %v", err)
		}

		logger.Infof(
			"Successfully submitted deposit sweep proof for transaction: %s",
			transactionHashFlag,
		)

		return nil
	},
}

// convertVaultAddress converts the vault's address from the chain's address
// to the common address. If the chain's address is nil, the returned value
// is zero-filled.
func convertVaultAddress(vault *chain.Address) common.Address {
	if vault == nil {
		return common.Address{}
	}

	return common.HexToAddress(string(*vault))
}

// parseTransactionInputs parses the transaction's inputs and returns the main
// UTXO and the vault.
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
			// The input is P2PKH or P2WPKH, so we found main UTXO. There should
			// be at most one main UTXO. If any input of this kind has already
			// been found, report an error.
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
			// The input is P2SH or P2WSH, so we found a deposit input. All
			// the deposits should have the same vault set or no vault at all.
			// If the vault if different than the vault from any previous
			// deposit input, report an error.
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
			// The type of the input is neither P2PKH, P2WPKH, P2SH or P2WSH.
			// Report an error.
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
