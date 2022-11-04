package maintainer

import (
	"context"
	"fmt"
	"time"

	"github.com/ipfs/go-log"
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/chain"
)

const (
	// Back-off time which should be applied when the relay is restarted.
	// It helps to avoid being flooded with error logs in case of a permanent error
	// in the relay.
	restartBackoffTime = 10 * time.Second

	// Back-off time applied when the relay has not submitted headers for
	// retargetting, i.e. it is up to date with the Bitcoin blockchain.
	// It helps to slow down the relay so that it does not overload the on-chain
	// client by asking it for its state too often.
	hardRetargetSleepTime = 5 * time.Minute

	// Back-off time applied when the relay has just submitted a new retarget.
	// It helps to speed up the relay when it may be behind the Bitcoin
	// blockchain.
	softRetargetSleepTime = 60 * time.Second

	// The number of blocks in a Bitcoin epoch.
	bitcoinEpochLength = 2016
)

var logger = log.Logger("keep-maintainer-relay")

// RelayChain is an interface that provides the ability to communicate with the
// relay on-chain contract.
type RelayChain interface {
	// Ready checks whether the relay is active (i.e. genesis has been
	// performed).
	Ready() (bool, error)

	// IsAuthorizationRequired checks whether the relay requires the address
	// submitting a retarget to be authorised in advance by governance.
	IsAuthorizationRequired() (bool, error)

	// IsAuthorized checks whether the given address has been authorised to
	// submit a retarget by governance.
	IsAuthorized(address chain.Address) (bool, error)

	// Signing returns the signing associated with the chain.
	Signing() chain.Signing

	// Retarget adds a new epoch to the relay by providing a proof
	// of the difficulty before and after the retarget.
	Retarget(headers []*bitcoin.BlockHeader) error

	// CurrentEpoch returns the number of the latest epoch whose difficulty is
	// proven to the relay. If the genesis epoch's number is set correctly, and
	// retargets along the way have been legitimate, the current epoch equals
	// the height of the block starting the most recent epoch, divided by 2016.
	CurrentEpoch() (uint64, error)

	// ProofLength returns the number of blocks required for each side of a
	// retarget proof: a retarget must provide `proofLength` blocks before
	// the retarget and `proofLength` blocks after it.
	ProofLength() (uint64, error)
}

func newRelay(
	ctx context.Context,
	btcChain bitcoin.Chain,
	chain RelayChain,
) *Relay {
	relay := &Relay{
		btcChain: btcChain,
		chain:    chain,
	}

	go relay.startControlLoop(ctx)

	return relay
}

// Relay is the part of maintainer responsible for maintaining the state of
// the relay on-chain contract.
type Relay struct {
	btcChain bitcoin.Chain
	chain    RelayChain
}

// startControlLoop starts the loop responsible for controlling the relay.
func (r *Relay) startControlLoop(ctx context.Context) {
	logger.Info("starting headers relay")

	defer func() {
		logger.Info("stopping headers relay")
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := r.submitHeaders(ctx)
			if err != nil {
				logger.Errorf("error while submitting headers [%v]", err)
			}
		}

		time.Sleep(restartBackoffTime)
	}
}

// submitHeaders submits block headers for retargets.
func (r *Relay) submitHeaders(ctx context.Context) error {
	if err := r.verifySubmissionEligibility(); err != nil {
		return fmt.Errorf(
			"relay maintainer not eligible to submit retargets [%v]",
			err,
		)
	}

	for {
		currentBlockNumber, err := r.btcChain.GetCurrentBlockNumber()
		if err != nil {
			return fmt.Errorf(
				"failed to get current block number [%v]",
				err,
			)
		}

		currentEpoch, err := r.chain.CurrentEpoch()
		if err != nil {
			return fmt.Errorf(
				"failed to get current epoch [%v]",
				err,
			)
		}

		proofLength, err := r.chain.ProofLength()
		if err != nil {
			return fmt.Errorf(
				"failed to get proof length [%v]",
				err,
			)
		}

		nextEpoch := currentEpoch + 1
		nextEpochHeight := nextEpoch * bitcoinEpochLength

		firstBlockHeaderHeight := uint(nextEpochHeight - proofLength)
		lastBlockHeaderHeight := uint(nextEpochHeight + proofLength - 1)

		retargetSleepTime := hardRetargetSleepTime

		if currentBlockNumber >= lastBlockHeaderHeight {
			headers, err := r.getBlockHeaders(
				firstBlockHeaderHeight,
				lastBlockHeaderHeight,
			)
			if err != nil {
				return fmt.Errorf(
					"failed to get block headers for retarget [%v]",
					err,
				)
			}

			if err := r.chain.Retarget(headers); err != nil {
				return fmt.Errorf(
					"failed to submit block headers [%v] for a retarget [%v]",
					headers,
					err,
				)
			}

			retargetSleepTime = softRetargetSleepTime
		}

		time.Sleep(retargetSleepTime)
	}
}

// verifySubmissionEligibility verifies whether a relay maintainer is eligible
// to submit block headers for retargets.
func (r *Relay) verifySubmissionEligibility() error {
	isReady, err := r.chain.Ready()
	if err != nil {
		return fmt.Errorf(
			"cannot check whether relay genesis has been performed [%v]",
			err,
		)
	}

	if !isReady {
		return fmt.Errorf("relay genesis has not been performed")
	}

	isAuthorizationRequired, err := r.chain.IsAuthorizationRequired()
	if err != nil {
		return fmt.Errorf(
			"cannot check whether authorization is required to submit "+
				"retargets [%v]",
			err,
		)
	}

	if !isAuthorizationRequired {
		return nil
	}

	maintainerAddress := r.chain.Signing().Address()

	isAuthorized, err := r.chain.IsAuthorized(maintainerAddress)
	if err != nil {
		return fmt.Errorf(
			"cannot check whether relay maintainer is authorized to "+
				"submit retargets [%v]",
			err,
		)
	}

	if !isAuthorized {
		return fmt.Errorf(
			"relay maintainer is not authorized to submit retargets",
		)
	}

	return nil
}

// getBlockHeaders returns block headers from the given range.
func (r *Relay) getBlockHeaders(firstHeaderHeight, lastHeaderHeight uint) (
	[]*bitcoin.BlockHeader, error,
) {
	var headers []*bitcoin.BlockHeader

	for height := firstHeaderHeight; height <= lastHeaderHeight; height++ {
		header, err := r.btcChain.GetBlockHeader(height)
		if err != nil {
			return []*bitcoin.BlockHeader{}, fmt.Errorf(
				"failed to get block header at height [%v]: [%v]",
				height,
				err,
			)
		}

		headers = append(headers, header)
	}

	return headers, nil
}
