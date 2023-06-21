package spv

import (
	"context"
	"fmt"
	"time"

	"github.com/ipfs/go-log/v2"

	"github.com/keep-network/keep-core/pkg/bitcoin"
)

var logger = log.Logger("keep-maintainer-spv")

func Initialize(
	ctx context.Context,
	config Config,
	chain Chain,
	btcChain bitcoin.Chain,
) {
	spvMaintainer := &spvMaintainer{
		config:   config,
		chain:    chain,
		btcChain: btcChain,
	}

	go spvMaintainer.startControlLoop(ctx)
}

type spvMaintainer struct {
	config   Config
	chain    Chain
	btcChain bitcoin.Chain
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
		case <-time.After(sm.config.RestartBackOffTime):
		case <-ctx.Done():
			return
		}
	}
}

func (sm *spvMaintainer) maintainSpv(ctx context.Context) error {
	logger.Infof("Maintaining SPV proof...")

	for {
		if err := sm.proveDepositSweepTransactions(); err != nil {
			return fmt.Errorf(
				"error while proving deposit sweep transactions: [%v]",
				err,
			)
		}

		// TODO: Add proving of other type of SPV transactions: redemption
		// transactions, moving funds transaction, etc.

		select {
		case <-time.After(sm.config.IdleBackOffTime):
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

	fmt.Println("depositSweepTransactions: ", depositSweepTransactions)

	// TODO: Assemble the proof and submit to the Bridge
	return nil
}

func (sm *spvMaintainer) getUnprovenDepositSweepTransactions() (
	[]*bitcoin.Transaction,
	error,
) {
	// TODO: Limit how far in the past we are looking for the events.
	//       Possibly store latest checked height in memory or file.
	depositSweepTransactionProposals, err :=
		sm.chain.PastDepositSweepProposalSubmittedEvents(nil)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get past deposit sweep proposal submitted events: [%v]",
			err,
		)
	}

	unprovenDepositSweepTransactions := []*bitcoin.Transaction{}

	for _, proposal := range depositSweepTransactionProposals {
		// TODO: Think what the limit of transactions should be.
		walletTransactions, err := sm.btcChain.GetTransactionsForPublicKeyHash(
			proposal.Proposal.WalletPublicKeyHash,
			5,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to get transactions for wallet: [%v]",
				err,
			)
		}

		for _, transaction := range walletTransactions {
			isUnprovenDepositSweepTransaction, err :=
				sm.isUnprovenDepositSweepTransaction(transaction)
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
) (bool, error) {
	// TODO: Look at transaction's inputs. One of the inputs may be the main UTXO,
	//       all the other inputs must be deposits;
	//       Use outpoint data of each input to build a deposit key and see if
	//       such a deposit exists in the Bridge;
	//       The one input that is not a deposit input should be the main UTXO;
	//       Verify this by checking the current main UTXO of the wallet as seen
	//       in the Bridge.

	return false, nil
}
