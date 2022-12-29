package maintainer

import (
	"context"
	"fmt"
	"time"

	"github.com/ipfs/go-log"
	"github.com/keep-network/keep-core/pkg/bitcoin"
)

const (
	// Default value for back-off time which should be applied when the Bitcoin
	// difficulty maintainer is restarted. It helps to avoid being flooded with
	// error logs in case of a permanent error in the Bitcoin difficulty
	// maintainer.
	bitcoinDifficultyDefaultRestartBackoffTime = 120 * time.Second

	// Default value for back-off time which should be applied when there are no
	// more Bitcoin epochs to be proven because the difficulty maintainer is
	// up-to-date with the Bitcoin blockchain or there are not enough blocks yet
	// to prove the new epoch.
	bitcoinDifficultyDefaultIdleBackOffTime = 60 * time.Second

	// The number of blocks in a Bitcoin difficulty epoch.
	bitcoinDifficultyEpochLength = 2016
)

var logger = log.Logger("maintainer-btcdiff")

var (
	errNotAuthorized = fmt.Errorf(
		"bitcoin difficulty maintainer has not been authorized to submit " +
			"block headers",
	)
	errNoGenesis = fmt.Errorf(
		"genesis has not been performed in the Bitcoin difficulty chain",
	)
)

func initializeBitcoinDifficultyMaintainer(
	ctx context.Context,
	btcChain bitcoin.Chain,
	chain BitcoinDifficultyChain,
	idleBackOffTime time.Duration,
	restartBackOffTime time.Duration,
) {
	bitcoinDifficultyMaintainer := &bitcoinDifficultyMaintainer{
		btcChain:           btcChain,
		chain:              chain,
		idleBackOffTime:    idleBackOffTime,
		restartBackOffTime: restartBackOffTime,
	}

	go bitcoinDifficultyMaintainer.startControlLoop(ctx)
}

// bitcoinDifficultyMaintainer is the part of maintainer responsible for
// maintaining the state of the Bitcoin difficulty on-chain contract.
type bitcoinDifficultyMaintainer struct {
	btcChain bitcoin.Chain
	chain    BitcoinDifficultyChain

	idleBackOffTime    time.Duration
	restartBackOffTime time.Duration
}

// startControlLoop starts the loop responsible for controlling the Bitcoin
// difficulty maintainer.
func (bdm *bitcoinDifficultyMaintainer) startControlLoop(ctx context.Context) {
	logger.Info("starting Bitcoin difficulty maintainer")

	defer func() {
		logger.Info("stopping Bitcoin difficulty maintainer")
	}()

	for {
		err := bdm.proveEpochs(ctx)
		if err != nil {
			logger.Errorf(
				"error while proving Bitcoin epochs: [%v]; restarting maintainer",
				err,
			)
		}

		select {
		case <-time.After(bdm.restartBackOffTime):
		case <-ctx.Done():
			return
		}
	}
}

// proveEpochs proves Bitcoin blockchain epochs in the Bitcoin difficulty chain.
func (bdm *bitcoinDifficultyMaintainer) proveEpochs(ctx context.Context) error {
	if err := bdm.verifySubmissionEligibility(); err != nil {
		return fmt.Errorf(
			"cannot verify submission eligibility: [%w]",
			err,
		)
	}

	for {
		epochProven, err := bdm.proveNextEpoch(ctx)
		if err != nil {
			return fmt.Errorf(
				"cannot prove Bitcoin blockchain epoch: [%w]",
				err,
			)
		}

		// Sleep for some time if the Bitcoin epoch was not proven (i.e. Bitcoin
		// difficulty chain is up-to-date or there are not enough block headers
		// in the new epoch). Do not sleep if a Bitcoin epoch was proven as
		// there are likely more Bitcoin epochs to prove.
		if !epochProven {
			select {
			case <-time.After(bdm.idleBackOffTime):
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}

// verifySubmissionEligibility verifies whether a maintainer is eligible to
// submit block headers to the Bitcoin difficulty chain.
func (bdm *bitcoinDifficultyMaintainer) verifySubmissionEligibility() error {
	isReady, err := bdm.chain.Ready()
	if err != nil {
		return fmt.Errorf(
			"cannot check whether genesis has been performed: [%w]",
			err,
		)
	}

	if !isReady {
		return errNoGenesis
	}

	authorizationRequired, err := bdm.chain.AuthorizationRequired()
	if err != nil {
		return fmt.Errorf(
			"cannot check whether authorization is required to submit "+
				"block headers: [%w]",
			err,
		)
	}

	if !authorizationRequired {
		return nil
	}

	maintainerAddress := bdm.chain.Signing().Address()

	isAuthorized, err := bdm.chain.IsAuthorized(maintainerAddress)
	if err != nil {
		return fmt.Errorf(
			"cannot check whether Bitcoin difficulty maintainer is authorized "+
				"to submit block headers: [%w]",
			err,
		)
	}

	if !isAuthorized {
		return errNotAuthorized
	}

	return nil
}

// proveNextEpoch proves a single Bitcoin blockchain epoch in the Bitcoin
// difficulty chain if there is an epoch to be proven. If it was possible to
// prove an epoch, it returns true. If it was not possible (i.e. Bitcoin
// difficulty chain is up-to-date or there are not enough headers in the new
// Bitcoin epoch), it returns false.
func (bdm *bitcoinDifficultyMaintainer) proveNextEpoch(ctx context.Context) (
	bool,
	error,
) {
	// The height of the Bitcoin blockchain.
	currentBlockHeight, err := bdm.btcChain.GetLatestBlockHeight()
	if err != nil {
		return false, fmt.Errorf(
			"failed to get latest block height: [%w]",
			err,
		)
	}

	// The current epoch proven in the Bitcoin difficulty chain.
	currentEpoch, err := bdm.chain.CurrentEpoch()
	if err != nil {
		return false, fmt.Errorf(
			"failed to get current epoch: [%w]",
			err,
		)
	}

	// The number of blocks required for each side of a retarget proof.
	proofLength, err := bdm.chain.ProofLength()
	if err != nil {
		return false, fmt.Errorf(
			"failed to get proof length: [%w]",
			err,
		)
	}

	// The new epoch to be proven in the Bitcoin difficulty chain.
	newEpoch := uint(currentEpoch) + 1

	// Height of the first block of the new epoch.
	newEpochHeight := newEpoch * bitcoinDifficultyEpochLength

	// The range of block headers to be pulled from the Bitcoin chain should
	// start `proofLength` blocks before the first block of the new difficulty
	// epoch and end `proofLength`-1 after it.
	// For example, if the new epoch begins at block 522144 and `proofLength`
	// is 3, then the range should be [522141, 522146]:
	// 522141 <- old difficulty target
	// 522142 <- old difficulty target
	// 522143 <- old difficulty target
	// << difficulty retarget >>
	// 522144 <- new difficulty target (first block of the new epoch)
	// 522145 <- new difficulty target
	// 522146 <- new difficulty target
	firstBlockHeaderHeight := newEpochHeight - uint(proofLength)
	lastBlockHeaderHeight := newEpochHeight + uint(proofLength) - 1

	// The required range of block headers can be pulled from the Bitcoin
	// blockchain only if the blockchain height is equal to or greater than
	// the end of the range.
	if currentBlockHeight >= lastBlockHeaderHeight {
		headers, err := bdm.getBlockHeaders(
			firstBlockHeaderHeight,
			lastBlockHeaderHeight,
		)
		if err != nil {
			return false, fmt.Errorf(
				"failed to get block headers from Bitcoin chain: [%w]",
				err,
			)
		}

		if err := bdm.chain.Retarget(headers); err != nil {
			return false, fmt.Errorf(
				"failed to submit block headers from range [%d:%d] to "+
					"the Bitcoin difficulty chain: [%w]",
				firstBlockHeaderHeight,
				lastBlockHeaderHeight,
				err,
			)
		}

		if err := bdm.waitForCurrentEpochUpdate(ctx, uint64(newEpoch)); err != nil {
			return false, fmt.Errorf(
				"error while waiting for current Bitcoin difficulty epoch "+
					"update: [%w]",
				err,
			)
		}

		logger.Infof(
			"successfully submitted block headers [%d:%d] to the Bitcoin "+
				"difficulty chain; the current proven epoch is [%d]",
			firstBlockHeaderHeight,
			lastBlockHeaderHeight,
			newEpoch,
		)

		return true, nil
	}

	if currentBlockHeight >= newEpochHeight {
		logger.Infof(
			"the Bitcoin difficulty chain has to be synced with the "+
				"Bitcoin blockchain; waiting for [%d] new blocks to "+
				"be mined to form a headers chain for retarget",
			lastBlockHeaderHeight-currentBlockHeight,
		)
	} else {
		logger.Infof(
			"the Bitcoin difficulty chain is up-to-date with the Bitcoin " +
				"blockchain",
		)
	}

	return false, nil
}

// getBlockHeaders returns block headers from the given range.
func (bdm *bitcoinDifficultyMaintainer) getBlockHeaders(
	firstHeaderHeight,
	lastHeaderHeight uint,
) (
	[]*bitcoin.BlockHeader, error,
) {
	var headers []*bitcoin.BlockHeader

	for height := firstHeaderHeight; height <= lastHeaderHeight; height++ {
		header, err := bdm.btcChain.GetBlockHeader(height)
		if err != nil {
			return []*bitcoin.BlockHeader{}, fmt.Errorf(
				"failed to get block header at height %d: [%w]",
				height,
				err,
			)
		}

		headers = append(headers, header)
	}

	return headers, nil
}

// waitForCurrentEpochUpdate waits until the current epoch in the Bitcoin
// difficulty chain is equal to or higher than the provided target epoch.
func (bdm *bitcoinDifficultyMaintainer) waitForCurrentEpochUpdate(
	ctx context.Context,
	targetEpoch uint64,
) error {
	for {
		currentEpoch, err := bdm.chain.CurrentEpoch()
		if err != nil {
			return fmt.Errorf("failed to get current epoch: [%w]", err)
		}

		if currentEpoch >= targetEpoch {
			break
		}

		logger.Infof(
			"waiting for bitcoin difficulty chain to reach epoch %d, "+
				"current proven epoch is %d",
			targetEpoch,
			currentEpoch,
		)

		select {
		case <-time.After(time.Second):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}
