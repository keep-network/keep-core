package maintainer

import (
	"context"
	"fmt"
	"time"

	"github.com/ipfs/go-log"
	"github.com/keep-network/keep-core/pkg/bitcoin"
)

const (
	// Default value for back-off time which should be applied when the bitcoin
	// difficulty maintainer is restarted. It helps to avoid being flooded with
	// error logs in case of a permanent error  in the relay.
	defaultRestartBackoffTime = 10 * time.Second

	// Default value for backoff time which should be applied after each attempt
	// to prove a single Bitcoin epoch.
	defaultEpochProvenBackOffTime = 60 * time.Second

	// The number of blocks in a Bitcoin epoch.
	bitcoinEpochLength = 2016
)

var logger = log.Logger("keep-maintainer-bitcoin-difficulty")

func initializeBitcoinDifficultyMaintainer(
	ctx context.Context,
	btcChain bitcoin.Chain,
	chain BitcoinDifficultyChain,
	epochProvenBackOffTime time.Duration,
	restartBackOffTime time.Duration,
) {
	bitcoinDifficultyMaintainer := &BitcoinDifficultyMaintainer{
		btcChain:               btcChain,
		chain:                  chain,
		epochProvenBackOffTime: epochProvenBackOffTime,
		restartBackOffTime:     restartBackOffTime,
	}

	go bitcoinDifficultyMaintainer.startControlLoop(ctx)
}

// BitcoinDifficultyMaintainer is the part of maintainer responsible for
// maintaining the state of the Bitcoin difficulty on-chain contract.
type BitcoinDifficultyMaintainer struct {
	btcChain bitcoin.Chain
	chain    BitcoinDifficultyChain

	epochProvenBackOffTime time.Duration
	restartBackOffTime     time.Duration
}

// startControlLoop starts the loop responsible for controlling the Bitcoin
// difficulty maintainer.
func (bdm *BitcoinDifficultyMaintainer) startControlLoop(ctx context.Context) {
	logger.Info("starting Bitcoin difficulty maintainer")

	defer func() {
		logger.Info("stopping Bitcoin difficulty maintainer")
	}()

	for {
		err := bdm.proveEpochs(ctx)
		if err != nil {
			logger.Errorf(
				"restarting relay maintainer due to error while proving "+
					"Bitcoin blockchain epochs [%v]",
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

// proveEpochs proves Bitcoin blockchain epochs in the relay chain.
func (bdm *BitcoinDifficultyMaintainer) proveEpochs(ctx context.Context) error {
	if err := bdm.verifySubmissionEligibility(); err != nil {
		return fmt.Errorf(
			"cannot proceed with proving Bitcoin blockchain epochs in the "+
				"relay chain [%v]",
			err,
		)
	}

	for {
		if err := bdm.proveSingleEpoch(); err != nil {
			return fmt.Errorf(
				"cannot prove Bitcoin blockchain epoch to the relay chain [%v]",
				err,
			)
		}

		select {
		case <-time.After(bdm.epochProvenBackOffTime):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// verifySubmissionEligibility verifies whether a relay maintainer is eligible
// to submit block headers to the relay chain.
func (bdm *BitcoinDifficultyMaintainer) verifySubmissionEligibility() error {
	isReady, err := bdm.chain.Ready()
	if err != nil {
		return fmt.Errorf(
			"cannot check whether relay genesis has been performed [%v]",
			err,
		)
	}

	if !isReady {
		return fmt.Errorf("relay genesis has not been performed")
	}

	isAuthorizationRequired, err := bdm.chain.IsAuthorizationRequired()
	if err != nil {
		return fmt.Errorf(
			"cannot check whether authorization is required to submit "+
				"block headers [%v]",
			err,
		)
	}

	if !isAuthorizationRequired {
		return nil
	}

	maintainerAddress := bdm.chain.Signing().Address()

	isAuthorized, err := bdm.chain.IsAuthorized(maintainerAddress)
	if err != nil {
		return fmt.Errorf(
			"cannot check whether relay maintainer is authorized to "+
				"submit block headers [%v]",
			err,
		)
	}

	if !isAuthorized {
		return fmt.Errorf(
			"relay maintainer has not been authorized to submit block headers",
		)
	}

	return nil
}

// proveSingleEpoch proves a single Bitcoin blockchain epoch in the relay chain
// if there is a Bitcoin blockchain epoch to be proven.
func (bdm *BitcoinDifficultyMaintainer) proveSingleEpoch() error {
	// The height of the Bitcoin blockchain.
	currentBlockNumber, err := bdm.btcChain.GetCurrentBlockNumber()
	if err != nil {
		return fmt.Errorf(
			"failed to get current block number [%v]",
			err,
		)
	}

	// The current epoch proven in the relay chain.
	currentEpoch, err := bdm.chain.CurrentEpoch()
	if err != nil {
		return fmt.Errorf(
			"failed to get current epoch [%v]",
			err,
		)
	}

	// The number of blocks required for each side of a retarget proof.
	proofLength, err := bdm.chain.ProofLength()
	if err != nil {
		return fmt.Errorf(
			"failed to get proof length [%v]",
			err,
		)
	}

	// The new epoch to be proven in the relay chain.
	newEpoch := currentEpoch + 1

	// Height of the first block of the new epoch.
	newEpochHeight := newEpoch * bitcoinEpochLength

	// The range of block headers to be pull from the Bitcoin chain should
	// start `proofLength` blocks before the first block of the new epoch
	// and end `proofLength`-1 after it.
	// For example, if the new epoch begins at block 522144 and `proofLength`
	// is 3, then the range should be [522141, 522146]:
	// 522141 <- old difficulty target
	// 522142 <- old difficulty target
	// 522143 <- old difficulty target
	// << difficulty retarget >>
	// 522144 <- new difficulty target (first block of the new epoch)
	// 522145 <- new difficulty target
	// 522146 <- new difficulty target
	firstBlockHeaderHeight := uint(newEpochHeight - proofLength)
	lastBlockHeaderHeight := uint(newEpochHeight + proofLength - 1)

	// The required range of block headers can be pulled from the Bitcoin
	// blockchain only if the blockchain height is equal to or greater than
	// the end of the range.
	if currentBlockNumber >= lastBlockHeaderHeight {
		headers, err := bdm.getBlockHeaders(
			firstBlockHeaderHeight,
			lastBlockHeaderHeight,
		)
		if err != nil {
			return fmt.Errorf(
				"failed to get block headers from Bitcoin chain [%v]",
				err,
			)
		}

		if err := bdm.chain.Retarget(headers); err != nil {
			return fmt.Errorf(
				"failed to submit block headers from range [%v, %v] to "+
					"the relay chain [%v]",
				firstBlockHeaderHeight,
				lastBlockHeaderHeight,
				err,
			)
		}

		logger.Infof(
			"Successfully submitted block headers [%v:%v] to the relay "+
				"chain. The current proven epoch is %v.",
			firstBlockHeaderHeight,
			lastBlockHeaderHeight,
			newEpoch,
		)
	} else {
		logger.Infof("The relay is up-to-date with the Bitcoin blockchain")
	}

	return nil
}

// getBlockHeaders returns block headers from the given range.
func (bdm *BitcoinDifficultyMaintainer) getBlockHeaders(
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
				"failed to get block header at height %v: [%v]",
				height,
				err,
			)
		}

		headers = append(headers, header)
	}

	return headers, nil
}
