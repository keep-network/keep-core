package spv

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/ipfs/go-log/v2"

	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/maintainer/btcdiff"
	"github.com/keep-network/keep-core/pkg/tbtc"
)

var logger = log.Logger("keep-maintainer-spv")

// The length of the Bitcoin difficulty epoch in blocks.
const difficultyEpochLength = 2016

func Initialize(
	ctx context.Context,
	config Config,
	spvChain Chain,
	btcDiffChain btcdiff.Chain,
	btcChain bitcoin.Chain,
) {
	spvMaintainer := &spvMaintainer{
		config:       config,
		spvChain:     spvChain,
		btcDiffChain: btcDiffChain,
		btcChain:     btcChain,
	}

	go spvMaintainer.startControlLoop(ctx)
}

type spvMaintainer struct {
	config       Config
	spvChain     Chain
	btcDiffChain btcdiff.Chain
	btcChain     bitcoin.Chain
}

func (sm *spvMaintainer) startControlLoop(ctx context.Context) {
	logger.Info("starting SPV maintainer")

	defer func() {
		logger.Info("stopping SPV maintainer")
	}()

	for {
		err := sm.maintainSpv(ctx)
		if err != nil {
			logger.Errorf(
				"error while maintaining SPV: [%v]; restarting maintainer",
				err,
			)
		}

		select {
		case <-time.After(sm.config.RestartBackoffTime):
		case <-ctx.Done():
			return
		}
	}
}

func (sm *spvMaintainer) maintainSpv(ctx context.Context) error {
	for {
		logger.Infof("starting deposit sweep proof task execution...")

		if err := sm.proveDepositSweepTransactions(); err != nil {
			return fmt.Errorf(
				"error while proving deposit sweep transactions: [%v]",
				err,
			)
		}

		logger.Infof(
			"deposit sweep proof task run completed; next run in [%s]",
			sm.config.IdleBackoffTime,
		)

		// TODO: Add proving of other type of SPV transactions: redemption
		// transactions, moving funds transaction, etc.

		select {
		case <-time.After(sm.config.IdleBackoffTime):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (sm *spvMaintainer) proveDepositSweepTransactions() error {
	depositSweepTransactions, err := sm.getUnprovenDepositSweepTransactions()
	if err != nil {
		return fmt.Errorf(
			"failed to get unproven deposit sweep transactions: [%v]",
			err,
		)
	}

	logger.Infof(
		"found [%d] unproven deposit sweep transaction(s)",
		len(depositSweepTransactions),
	)

	for _, transaction := range depositSweepTransactions {
		// Print the transaction in the same endianness as block explorers do.
		transactionHashStr := transaction.Hash().Hex(bitcoin.ReversedByteOrder)

		logger.Infof(
			"proceeding with deposit sweep proof for transaction [%s]",
			transactionHashStr,
		)

		accumulatedConfirmations, err := sm.btcChain.GetTransactionConfirmations(
			transaction.Hash(),
		)
		if err != nil {
			return fmt.Errorf(
				"failed to get transaction confirmations: [%v]",
				err,
			)
		}

		isProofWithinRelayRange, requiredConfirmations, err := sm.getProofInfo(
			transaction.Hash(),
		)
		if err != nil {
			return fmt.Errorf("failed to get proof info: [%v]", err)
		}

		if !isProofWithinRelayRange {
			// The required proof goes outside the previous and current
			// difficulty epochs as seen by the relay. Skip the transaction. It
			// will most likely be proven later.
			logger.Warnf(
				"skipped proving deposit sweep transaction [%s]; the range "+
					"of the required proof goes outside the previous and "+
					"current difficulty epochs as seen by the relay",
				transactionHashStr,
			)
			continue
		}

		if accumulatedConfirmations < uint(requiredConfirmations) {
			// Skip the transaction as it has not accumulated enough
			// confirmations. It will be proven later.
			logger.Infof(
				"skipped proving deposit sweep transaction [%s]; transaction "+
					"has [%v/%v] confirmations",
				transactionHashStr,
				accumulatedConfirmations,
				requiredConfirmations,
			)
			continue
		}

		_, proof, err := bitcoin.AssembleSpvProof(
			transaction.Hash(),
			uint(requiredConfirmations),
			sm.btcChain,
		)
		if err != nil {
			return fmt.Errorf("failed to assemble SPV proof: [%v]", err)
		}

		mainUTXO, vault, err := parseTransactionInputs(
			sm.btcChain,
			sm.spvChain,
			transaction,
		)
		if err != nil {
			return fmt.Errorf(
				"error while parsing transaction inputs: [%v]",
				err,
			)
		}

		if err := sm.spvChain.SubmitDepositSweepProofWithReimbursement(
			transaction,
			proof,
			mainUTXO,
			vault,
		); err != nil {
			return fmt.Errorf(
				"failed to submit deposit sweep proof with reimbursement: [%v]",
				err,
			)
		}

		logger.Infof(
			"successfully submitted deposit sweep proof for transaction [%s]",
			transactionHashStr,
		)
	}

	logger.Infof("finished round of proving deposit sweep transactions")

	return nil
}

func (sm *spvMaintainer) getUnprovenDepositSweepTransactions() (
	[]*bitcoin.Transaction,
	error,
) {
	blockCounter, err := sm.spvChain.BlockCounter()
	if err != nil {
		return nil, fmt.Errorf("failed to get block counter: [%v]", err)
	}

	currentBlock, err := blockCounter.CurrentBlock()
	if err != nil {
		return nil, fmt.Errorf("failed to get current block: [%v]", err)
	}

	// Calculate the starting block of the range in which the events will be
	// searched for.
	startBlock := currentBlock - sm.config.HistoryDepth

	depositSweepTransactionProposals, err :=
		sm.spvChain.PastDepositSweepProposalSubmittedEvents(
			&tbtc.DepositSweepProposalSubmittedEventFilter{
				StartBlock: startBlock,
			},
		)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get past deposit sweep proposal submitted events: [%v]",
			err,
		)
	}

	// There will often be multiple events emitted for a single wallet. Prepare
	// a list of unique wallet public key hashes.
	walletPublicKeyHashes := uniqueWalletPublicKeyHashes(
		depositSweepTransactionProposals,
	)

	unprovenDepositSweepTransactions := []*bitcoin.Transaction{}

	for _, walletPublicKeyHash := range walletPublicKeyHashes {
		wallet, err := sm.spvChain.GetWallet(walletPublicKeyHash)
		if err != nil {
			return nil, fmt.Errorf("failed to get wallet: [%v]", err)
		}

		if wallet.State != tbtc.StateLive &&
			wallet.State != tbtc.StateMovingFunds {
			// The wallet can only submit deposit sweep proofs if it's `Live` or
			// `MovingFunds`. If the state is different skip it.
			logger.Infof(
				"skipped proving deposit sweep transactions for wallet [%x] "+
					"because of wallet state [%v]",
				walletPublicKeyHash,
				wallet.State,
			)
			continue
		}

		walletTransactions, err := sm.btcChain.GetTransactionsForPublicKeyHash(
			walletPublicKeyHash,
			sm.config.TransactionLimit,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to get transactions for wallet: [%v]",
				err,
			)
		}

		for _, transaction := range walletTransactions {
			isUnprovenDepositSweepTransaction, err :=
				sm.isUnprovenDepositSweepTransaction(
					transaction,
					walletPublicKeyHash,
				)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to check if transaction is an unproven deposit sweep "+
						"transaction: [%v]",
					err,
				)
			}

			if isUnprovenDepositSweepTransaction {
				unprovenDepositSweepTransactions = append(
					unprovenDepositSweepTransactions,
					transaction,
				)
			}
		}
	}

	return unprovenDepositSweepTransactions, nil
}

func (sm *spvMaintainer) isUnprovenDepositSweepTransaction(
	transaction *bitcoin.Transaction,
	walletPublicKeyHash [20]byte,
) (bool, error) {
	// If the transaction does not have exactly one output, it cannot be a
	// deposit sweep transaction.
	if len(transaction.Outputs) != 1 {
		return false, nil
	}

	hasDepositInputs := false

	// Look at the transaction's inputs. All the inputs must be deposit inputs,
	// except for one input which can be the main UTXO.
	for _, input := range transaction.Inputs {
		fundingTransactionHash := input.Outpoint.TransactionHash
		fundingOutpointIndex := input.Outpoint.OutputIndex

		// Check if the input is a deposit input.
		deposit, found, err := sm.spvChain.GetDepositRequest(
			fundingTransactionHash,
			fundingOutpointIndex,
		)
		if err != nil {
			return false, fmt.Errorf("failed to get deposit request: [%v]", err)
		}

		if !found {
			// The input is not a deposit input. The transaction can still be
			// a deposit sweep transaction, since the input may be the main UTXO.

			// Check if the input represents the current main UTXO of the wallet.
			// Notice that we don't have to verify if there is only one main
			// UTXO among the transaction's inputs since only one input may have
			// such a structure that the calculated hash will match the wallet's
			// main UTXO hash stored on-chain.
			isMainUtxo, err := sm.isInputCurrentWalletsMainUTXO(
				fundingTransactionHash,
				fundingOutpointIndex,
				walletPublicKeyHash,
			)
			if err != nil {
				return false, fmt.Errorf(
					"failed to check if input is the main UTXO",
				)
			}

			// The input is not the current main UTXO of the wallet. The
			// transaction is either a deposit sweep transaction that is already
			// proven or it's not a deposit sweep transaction at all.
			if !isMainUtxo {
				return false, nil
			}

			// The input is the current main UTXO of the wallet. Proceed with
			// checking other inputs.
		} else {
			// The input is a deposit input. Check if it swept or not.
			if deposit.SweptAt.Equal(time.Unix(0, 0)) {
				// The input is a deposit and it's unswept.
				hasDepositInputs = true
			} else {
				// The input is a deposit, but it's already swept.
				// The transaction must a deposit sweep transaction, but it's
				// already proven.
				return false, nil
			}
		}
	}

	// All the inputs represent either unswept deposits or the current main UTXO.
	// As the final check verify if at least one of them was a deposit input.
	// This will distinguish a deposit sweep transaction from a different
	// transaction type that may have the main UTXO as input, e.g. redemption.
	return hasDepositInputs, nil
}

func (sm *spvMaintainer) isInputCurrentWalletsMainUTXO(
	fundingTxHash bitcoin.Hash,
	fundingOutputIndex uint32,
	walletPublicKeyHash [20]byte,
) (bool, error) {
	// Get the transaction the input originated from to calculate the input value.
	previousTransaction, err := sm.btcChain.GetTransaction(fundingTxHash)
	if err != nil {
		return false, fmt.Errorf("failed to get previous transaction: [%v]", err)
	}
	fundingOutputValue := previousTransaction.Outputs[fundingOutputIndex].Value

	// Assume the input is the main UTXO and calculate hash.
	mainUtxoHash := sm.spvChain.ComputeMainUtxoHash(&bitcoin.UnspentTransactionOutput{
		Outpoint: &bitcoin.TransactionOutpoint{
			TransactionHash: fundingTxHash,
			OutputIndex:     fundingOutputIndex,
		},
		Value: fundingOutputValue,
	})

	// Get the wallet and check if its main UTXO matches the calculated hash.
	wallet, err := sm.spvChain.GetWallet(walletPublicKeyHash)
	if err != nil {
		return false, fmt.Errorf("failed to get wallet: [%v]", err)
	}

	return bytes.Equal(mainUtxoHash[:], wallet.MainUtxoHash[:]), nil
}

// getProofInfo returns information about the SPV proof. It includes the
// information whether the transaction proof range is within the previous and
// current difficulty epochs as seen by the relay and the required number of
// confirmations.
func (sm *spvMaintainer) getProofInfo(transactionHash bitcoin.Hash) (
	bool, uint, error,
) {
	latestBlockHeight, err := sm.btcChain.GetLatestBlockHeight()
	if err != nil {
		return false, 0, fmt.Errorf(
			"failed to get latest block height: [%v]",
			err,
		)
	}

	accumulatedConfirmations, err := sm.btcChain.GetTransactionConfirmations(
		transactionHash,
	)
	if err != nil {
		return false, 0, fmt.Errorf(
			"failed to get transaction confirmations: [%v]",
			err,
		)
	}

	txProofDifficultyFactor, err := sm.spvChain.TxProofDifficultyFactor()
	if err != nil {
		return false, 0, fmt.Errorf(
			"failed to get transaction proof difficulty factor: [%v]",
			err,
		)
	}

	// Calculate the starting block of the proof and the difficulty epoch number
	// it belongs to.
	proofStartBlock := uint64(latestBlockHeight - accumulatedConfirmations + 1)
	proofStartEpoch := proofStartBlock / difficultyEpochLength

	// Calculate the ending block of the proof and the difficulty epoch number
	// it belongs to.
	proofEndBlock := proofStartBlock + txProofDifficultyFactor.Uint64() - 1
	proofEndEpoch := proofEndBlock / difficultyEpochLength

	// Get the current difficulty epoch number as seen by the relay. Subtract
	// one to get the previous epoch number.
	currentEpoch, err := sm.btcDiffChain.CurrentEpoch()
	if err != nil {
		return false, 0, fmt.Errorf("failed to get current epoch: [%v]", err)
	}
	previousEpoch := currentEpoch - 1

	// There are only three possible valid combinations of the proof's block
	// headers range: the proof must either be entirely in the previous epoch,
	// must be entirely in the current epoch or must span the previous and
	// current epochs.

	// If the proof is entirely within the current epoch, required confirmations
	// does not need to be adjusted.
	if proofStartEpoch == currentEpoch &&
		proofEndEpoch == currentEpoch {
		return true, uint(txProofDifficultyFactor.Uint64()), nil
	}

	// If the proof is entirely within the previous epoch, required confirmations
	// does not need to be adjusted.
	if proofStartEpoch == previousEpoch &&
		proofEndEpoch == previousEpoch {
		return true, uint(txProofDifficultyFactor.Uint64()), nil
	}

	// If the proof spans the previous and current difficulty epochs, the
	// required confirmations may have to be adjusted. The reason for this is
	// that there may be a drop in the value of difficulty between the current
	// and the previous epochs. Example:
	// Let's assume the transaction was done near the end of an epoch, so that
	// part of the proof (let's say two block headers) is in the previous epoch
	// and part of it is in the current epoch.
	// If the previous epoch difficulty is 50 and the current epoch difficulty
	// is 30, the total required difficulty of the proof will be transaction
	// difficulty factor times previous difficulty: 6 * 50 = 300.
	// However, if we simply use transaction difficulty factor to get the number
	// of blocks we will end up with the difficulty sum that is too low:
	// 50 + 50 + 30 + 30 + 30 + 30 = 220. To calculate the correct number of
	// block headers needed we need to find how much difficulty needs to come
	// from from the current epoch block headers: 300 - 2*50 = 200 and divide
	// it by the current difficulty: 200 / 30 = 6 and add 1, because there
	// was a remainder. So the number of block headers from the current epoch
	// would be 7. The total number of block headers would be 9 and the sum
	// of their difficulties would be: 50 + 50 + 30 + 30 + 30 + 30 + 30 + 30 +
	// 30 = 310 which is enough to prove the transaction.
	if proofStartEpoch == previousEpoch &&
		proofEndEpoch == currentEpoch {
		currentEpochDifficulty, previousEpochDifficulty, err :=
			sm.btcDiffChain.GetCurrentAndPrevEpochDifficulty()
		if err != nil {
			return false, 0, fmt.Errorf(
				"failed to get Bitcoin epoch difficulties: [%v]",
				err,
			)
		}

		// Calculate the total difficulty that is required for the proof. The
		// proof begins in the previous difficulty epoch, therefore the total
		// required difficulty will be the previous epoch difficulty times
		// transaction proof difficulty factor.
		totalDifficultyRequired := new(big.Int).Mul(
			previousEpochDifficulty,
			txProofDifficultyFactor,
		)

		// Calculate the number of block headers in the proof that will come
		// from the previous difficulty epoch.
		numberOfBlocksPreviousEpoch :=
			uint64(difficultyEpochLength - proofStartBlock%difficultyEpochLength)

		// Calculate how much difficulty the blocks from the previous epoch part
		// of the proof have in total.
		totalDifficultyPreviousEpoch := new(big.Int).Mul(
			big.NewInt(int64(numberOfBlocksPreviousEpoch)),
			previousEpochDifficulty,
		)

		// Calculate how much difficulty must come from the current epoch.
		totalDifficultyCurrentEpoch := new(big.Int).Sub(
			totalDifficultyRequired,
			totalDifficultyPreviousEpoch,
		)

		// Calculate how many blocks from the current epoch we need.
		remainder := new(big.Int)
		numberOfBlocksCurrentEpoch, remainder := new(big.Int).DivMod(
			totalDifficultyCurrentEpoch,
			currentEpochDifficulty,
			remainder,
		)
		// If there is a remainder, it means there is still some amount of
		// difficulty missing that is less than one block difficulty. We need to
		// account for that by adding one additional block.
		if remainder.Cmp(big.NewInt(0)) > 0 {
			numberOfBlocksCurrentEpoch.Add(
				numberOfBlocksCurrentEpoch,
				big.NewInt(1),
			)
		}

		// The total required number of confirmations is the sum of blocks from
		// the previous and current epochs.
		requiredConfirmations := numberOfBlocksPreviousEpoch +
			numberOfBlocksCurrentEpoch.Uint64()

		return true, uint(requiredConfirmations), nil
	}

	// If we entered here, it means that the proof's block headers range goes
	// outside the previous or current difficulty epochs as seen by the relay.
	// The reason for this is most likely that transaction entered the Bitcoin
	// blockchain within the very new difficulty epoch that is not yet proven in
	// the relay. In that case the transaction will be proven in the future.
	// The other case could be that the transaction is older than the last two
	// Bitcoin difficulty epochs. In that case the transaction will soon leave
	// the sliding window of recent transactions.
	return false, 0, nil
}

// uniqueWalletPublicKeyHashes parses the list of events and returns a list of
// unique wallet public key hashes.
func uniqueWalletPublicKeyHashes(
	events []*tbtc.DepositSweepProposalSubmittedEvent,
) [][20]byte {
	cache := make(map[string]struct{})
	var publicKeyHashes [][20]byte

	for _, event := range events {
		key := event.Proposal.WalletPublicKeyHash
		strKey := hex.EncodeToString(key[:])

		// Check for uniqueness
		if _, exists := cache[strKey]; !exists {
			cache[strKey] = struct{}{}
			publicKeyHashes = append(publicKeyHashes, key)
		}
	}

	return publicKeyHashes
}
