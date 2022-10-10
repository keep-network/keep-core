package tbtc

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"sort"

	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/tecdsa/retry"
	"github.com/keep-network/keep-core/pkg/tecdsa/signing"
	"golang.org/x/exp/slices"
)

// signingRetryLoop is a struct that encapsulates the signing retry logic.
type signingRetryLoop struct {
	signingGroupMemberIndex group.MemberIndex
	signingGroupOperators   chain.Addresses
	inactiveOperatorsSet    map[chain.Address]bool

	chainConfig *ChainConfig

	seed int64

	attemptCounter    uint
	attemptStartBlock uint64

	// We use a separate counter for the random retry algorithm because we
	// try to exclude inactive members in the first attempts and then switch
	// to the random retry mechanism. In result, an attempt is not always
	// the same as one run of the random retry algorithm.
	randomRetryCounter uint

	delayBlocks uint64
}

func newSigningRetryLoop(
	message *big.Int,
	initialStartBlock uint64,
	signingGroupMemberIndex group.MemberIndex,
	signingGroupOperators chain.Addresses,
	chainConfig *ChainConfig,
) *signingRetryLoop {
	// Compute the 8-byte seed needed for the random actions. We take
	// the first 8 bytes of the hash of the signed message. This allows us to
	// not care in this piece of the code about the length of the message and
	// how this message is proposed.
	messageSha256 := sha256.Sum256(message.Bytes())
	seed := int64(binary.BigEndian.Uint64(messageSha256[:8]))

	return &signingRetryLoop{
		signingGroupMemberIndex: signingGroupMemberIndex,
		signingGroupOperators:   signingGroupOperators,
		inactiveOperatorsSet:    make(map[chain.Address]bool),
		chainConfig:             chainConfig,
		attemptCounter:          0,
		attemptStartBlock:       initialStartBlock,
		seed:                    seed,
		randomRetryCounter:      0,
		delayBlocks:             5,
	}
}

// signingAttemptParams represents parameters of a signing attempt.
type signingAttemptParams struct {
	number                 uint
	startBlock             uint64
	excludedMembersIndexes []group.MemberIndex
}

// signingAttemptFn represents a function performing a signing attempt.
type signingAttemptFn func(*signingAttemptParams) (*signing.Result, error)

// waitWithSigningAttemptFn represents a function blocking the attempt execution
// until the given block height.
type waitWithSigningAttemptFn func(context.Context, uint64) error

// start begins the signing retry loop using the given signing attempt function.
// The retry loop terminates when the signing result is produced or the ctx
// parameter is done, whatever comes first.
func (srl *signingRetryLoop) start(
	ctx context.Context,
	waitWithSigningAttemptFn waitWithSigningAttemptFn,
	signingAttemptFn signingAttemptFn,
) (*signing.Result, error) {
	// All signing group operators should be qualified for the first attempt.
	qualifiedOperatorsSet := srl.signingGroupOperators.Set()

	for {
		srl.attemptCounter++

		// In order to start attempts >1 in the right place, we need to
		// determine how many blocks were taken by previous attempts. We assume
		// the worst case that each attempt failed at the end of the signing
		// protocol.
		//
		// That said, we need to increment the previous attempt start
		// block by the number of blocks equal to the protocol duration and
		// by some additional delay blocks. We need a small fixed delay in
		// order to mitigate all corner cases where the actual attempt duration
		// was slightly longer than the expected duration determined by the
		// signing.ProtocolBlocks function.
		//
		// For example, the attempt may fail at the end of the protocol but the
		// error is returned after some time and more blocks than expected are
		// mined in the meantime.
		if srl.attemptCounter > 1 {
			srl.attemptStartBlock = srl.attemptStartBlock +
				signing.ProtocolBlocks() +
				srl.delayBlocks
		}

		// Exclude all members controlled by the operators that were not
		// qualified for the current attempt.
		includedMembersIndexes := make([]group.MemberIndex, 0)
		excludedMembersIndexes := make([]group.MemberIndex, 0)
		for i, operator := range srl.signingGroupOperators {
			memberIndex := group.MemberIndex(i + 1)

			if qualifiedOperatorsSet[operator] {
				includedMembersIndexes = append(
					includedMembersIndexes,
					memberIndex,
				)
			} else {
				excludedMembersIndexes = append(
					excludedMembersIndexes,
					memberIndex,
				)
			}
		}

		// Make sure we always use just the smallest required count of
		// signing members for performance reasons
		if len(includedMembersIndexes) > srl.chainConfig.HonestThreshold {
			// #nosec G404 (insecure random number source (rand))
			// Shuffling does not require secure randomness.
			rng := rand.New(rand.NewSource(
				srl.seed + int64(srl.attemptCounter),
			))
			// Sort in ascending order just in case.
			sort.Slice(includedMembersIndexes, func(i, j int) bool {
				return includedMembersIndexes[i] < includedMembersIndexes[j]
			})
			// Shuffle the included members slice to randomize the
			// selection of additionally excluded members.
			rng.Shuffle(len(includedMembersIndexes), func(i, j int) {
				includedMembersIndexes[i], includedMembersIndexes[j] =
					includedMembersIndexes[j], includedMembersIndexes[i]
			})
			// Get the surplus of included members and add them to
			// the excluded members list.
			excludedMembersIndexes = append(
				excludedMembersIndexes,
				includedMembersIndexes[srl.chainConfig.HonestThreshold:]...,
			)
			// Sort the resulting excluded members list in ascending order.
			sort.Slice(excludedMembersIndexes, func(i, j int) bool {
				return excludedMembersIndexes[i] < excludedMembersIndexes[j]
			})
		}

		attemptSkipped := slices.Contains(
			excludedMembersIndexes,
			srl.signingGroupMemberIndex,
		)

		// Wait for the right moment to execute the signingAttemptFn, as
		// calculated in srl.attemptStartBlock.
		err := waitWithSigningAttemptFn(ctx, srl.attemptStartBlock)
		if err != nil {
			return nil, fmt.Errorf(
				"failed waiting for block [%v] for attempt [%v]: [%v]",
				srl.attemptStartBlock,
				srl.attemptCounter,
				err,
			)
		}

		// Check the loop stop signal.
		if ctx.Err() != nil {
			return nil, nil
		}

		var result *signing.Result
		var attemptErr error

		if !attemptSkipped {
			result, attemptErr = signingAttemptFn(&signingAttemptParams{
				number:                 srl.attemptCounter,
				startBlock:             srl.attemptStartBlock,
				excludedMembersIndexes: excludedMembersIndexes,
			})
			if attemptErr != nil {
				var imErr *signing.InactiveMembersError
				if errors.As(attemptErr, &imErr) {
					for _, memberIndex := range imErr.InactiveMembersIndexes {
						operator := srl.signingGroupOperators[memberIndex-1]
						srl.inactiveOperatorsSet[operator] = true
					}
				}
			}
		}

		if attemptSkipped || attemptErr != nil {
			var err error
			qualifiedOperatorsSet, err = srl.qualifiedOperatorsSet()
			if err != nil {
				return nil, fmt.Errorf(
					"cannot get qualified operators for attempt [%v]: [%w]",
					srl.attemptCounter+1,
					err,
				)
			}

			continue
		}

		return result, nil
	}
}

// qualifiedOperatorsSet returns a set of operators qualified to participate
// in the given signing attempt.
func (srl *signingRetryLoop) qualifiedOperatorsSet() (
	map[chain.Address]bool,
	error,
) {
	// If this is one of the first attempts and random retries were not started
	// yet, check if there are known inactive operators. If the honest threshold
	// can be maintained, just exclude the members controlled by the inactive
	// operators from the qualified set.
	if srl.attemptCounter <= 5 &&
		srl.randomRetryCounter == 0 &&
		len(srl.inactiveOperatorsSet) > 0 {
		qualifiedOperators := make(chain.Addresses, 0)
		for _, operator := range srl.signingGroupOperators {
			if !srl.inactiveOperatorsSet[operator] {
				qualifiedOperators = append(qualifiedOperators, operator)
			}
		}

		// If this attempt pushes us below the honest threshold we are falling
		// back to the random retry algorithm that excludes specific members
		// from the original signing group.
		if len(qualifiedOperators) >= srl.chainConfig.HonestThreshold {
			return qualifiedOperators.Set(), nil
		}
	}

	// In any other case, try to make a random retry.
	qualifiedOperators, err := retry.EvaluateRetryParticipantsForSigning(
		srl.signingGroupOperators,
		srl.seed,
		srl.randomRetryCounter,
		uint(srl.chainConfig.HonestThreshold),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"random operator selection failed: [%w]",
			err,
		)
	}

	srl.randomRetryCounter++
	return chain.Addresses(qualifiedOperators).Set(), nil
}
