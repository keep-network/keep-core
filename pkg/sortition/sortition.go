package sortition

import (
	"context"
	"errors"
	"math/big"
	"time"

	"github.com/ipfs/go-log"

	corechain "github.com/keep-network/keep-core/pkg/chain"
)

// retryDelay defines the delay between retries related to the registration logic
// that do not have their own specific values (like for example `eligibilityRetryDelay`
// for sortition pool join eligibility checks).
const retryDelay = 1 * time.Second

// operatorRegistrationRetryDelay defines the delay between checks whether the operator
// is registered by a staking provider.
var operatorRegistrationRetryDelay = 5 * time.Minute

// eligibilityRetryDelay defines the delay between checks whether the operator
// is eligible to join the sortition pool.
const eligibilityRetryDelay = 20 * time.Minute

var logger = log.Logger("keep-sortition")

// RegisterAndMonitorStatus checks whether the operator is registered by a staking
// provider and joins the sortition pool if the operator is eligible.
func RegisterAndMonitorStatus(
	ctx context.Context,
	blockCounter corechain.BlockCounter,
	chainSortitionHandle Handle,
) {
	go func() {
		operatorRegisteredChan := make(chan string)
		go waitUntilRegistered(ctx, blockCounter, chainSortitionHandle, operatorRegisteredChan)

		for {
			select {
			case <-ctx.Done():
				return
			case stakingProvider := <-operatorRegisteredChan:
				logger.Infof(
					"operator is registered for a staking provider [%s]",
					stakingProvider,
				)

				isInPool, err := chainSortitionHandle.IsOperatorInPool()
				if err != nil {
					logger.Errorf(
						"failed to verify if the operator is in the sortition pool: [%v]",
						err,
					)
					time.Sleep(retryDelay) // TODO: #413 Replace with backoff.
					continue
				}

				if !isInPool {
					// if the operator is not in the sortition pool, we need to
					// join the sortition pool
					joinSortitionPoolWhenEligible(ctx, stakingProvider, blockCounter, chainSortitionHandle)
				}

				return
			}
		}
	}()

	go func() {
		joinedPoolChan := make(chan struct{})
		go waitUntilJoined(ctx, blockCounter, chainSortitionHandle, joinedPoolChan)

		for {
			select {
			case <-ctx.Done():
				return
			case <-joinedPoolChan:
				logger.Infof("operator is in the sortition pool; starting monitoring...")

				// TODO: Monitor status
				// TODO: Rejoin sortition pool if removed

				monitorSignerPoolStatus()
			}
		}
	}()
}

// joinSortitionPoolWhenEligible checks current operator's eligibility to join
// the sortition pool and if it is positive, joins the pool.
// If the operator is not eligible, it executes the check for each new mined
// block until the operator is finally eligible and can join the pool.
func joinSortitionPoolWhenEligible(
	parentCtx context.Context,
	stakingProvider string,
	blockCounter corechain.BlockCounter,
	chainSortitionHandle Handle,
) {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	newBlockChan := blockCounter.WatchBlocks(ctx) // TODO: Check every X blocks

	for {
		select {
		case <-newBlockChan:
			eligibleStake, err := chainSortitionHandle.EligibleStake(stakingProvider)
			if err != nil {
				logger.Errorf(
					"failed to verify if the operator [%s] is eligible to join the sortition pool: [%v]",
					stakingProvider,
					err,
				)
				time.Sleep(retryDelay) // TODO: #413 Replace with backoff.
				continue
			}

			if eligibleStake.Cmp(big.NewInt(0)) == 0 {
				logger.Warnf("operator is not eligible to join the sortition pool")
				time.Sleep(eligibilityRetryDelay)
				continue
			}

			// TODO: Check if the sortition pool is unlocked

			logger.Infof("joining the sortition pool...")

			if err := chainSortitionHandle.JoinSortitionPool(); err != nil {
				logger.Errorf("failed to join the sortition pool: [%v]", err)
				time.Sleep(retryDelay) // TODO: #413 Replace with backoff.
				continue
			}

			return
		case <-ctx.Done():
			return
		}
	}
}

// waitUntilJoined blocks until the operator is registered by a staking provider.
func waitUntilRegistered(
	parentCtx context.Context,
	blockCounter corechain.BlockCounter,
	chainSortitionHandle Handle,
	operatorRegisteredChan chan string,
) {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	newBlockChan := blockCounter.WatchBlocks(ctx) // TODO: Check every X blocks

	for {
		select {
		case <-newBlockChan:
			stakingProvider, err := chainSortitionHandle.OperatorToStakingProvider()
			if err != nil {
				if errors.Is(err, ErrOperatorNotRegistered) {
					logger.Warn(
						"operator is not registered; please make sure a staking provider registered the operator",
					)
					time.Sleep(operatorRegistrationRetryDelay)
				} else {
					logger.Errorf(
						"failed to check if the operator is registered for a staking provider: [%v]",
						err,
					)
					time.Sleep(retryDelay) // TODO: #413 Replace with backoff.
				}
				continue
			}

			operatorRegisteredChan <- stakingProvider
			close(operatorRegisteredChan)
			return
		case <-ctx.Done():
			return
		}
	}
}

// waitUntilJoined blocks until the operator joins the sortition pool.
func waitUntilJoined(
	parentCtx context.Context,
	blockCounter corechain.BlockCounter,
	chainSortitionHandle Handle,
	outChan chan struct{},
) {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	newBlockChan := blockCounter.WatchBlocks(ctx) // TODO: Check every X blocks

	for {
		select {
		case <-newBlockChan:
			isInPool, err := chainSortitionHandle.IsOperatorInPool()
			if err != nil {
				logger.Errorf(
					"failed to verify if the operator is in the sortition pool: [%w]",
					err,
				)
				time.Sleep(retryDelay) // TODO: #413 Replace with backoff.
				continue
			}
			if !isInPool {
				logger.Debugf("operator is not yet in the sortition pool, waiting...")
				continue
			}

			close(outChan)
			return
		case <-ctx.Done():
			return
		}
	}
}

func monitorSignerPoolStatus() {
	for {
		// TODO: Implement
	}
}
