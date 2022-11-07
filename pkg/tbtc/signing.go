package tbtc

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/ipfs/go-log/v2"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/tecdsa/retry"
	"github.com/keep-network/keep-core/pkg/tecdsa/signing"
	"golang.org/x/exp/slices"
	"math/big"
	"math/rand"
	"sort"
)

// signingAttemptMaxBlockDuration determines the maximum block duration of a
// single signing attempt.
const signingAttemptMaxBlockDuration = 100

// signingRequestID represents the identifier of a signing protocol instance.
type signingRequestID [32]byte

// newSigningRequestID creates a new signing request ID based on the provided
// messages being subject of the given request.
func newSigningRequestID(messages []*big.Int) signingRequestID {
	var buffer bytes.Buffer

	for _, message := range messages {
		buffer.Write(message.Bytes())
	}

	return sha256.Sum256(buffer.Bytes())
}

func (sri signingRequestID) String() string {
	return hex.EncodeToString(sri[:])
}

// signingAnnouncer represents a component responsible for exchanging readiness
// announcements for the given signing attempt of the given message.
type signingAnnouncer interface {
	Announce(
		ctx context.Context,
		memberIndex group.MemberIndex,
		sessionID string,
	) ([]group.MemberIndex, error)
}

// signingRetryLoop is a struct that encapsulates the signing retry logic.
type signingRetryLoop struct {
	logger log.StandardLogger

	requestID signingRequestID

	signingGroupMemberIndex group.MemberIndex
	signingGroupOperators   chain.Addresses

	chainConfig *ChainConfig

	announcer                signingAnnouncer
	announcementDelayBlocks  uint64
	announcementActiveBlocks uint64

	attemptCounter     uint
	attemptStartBlock  uint64
	attemptSeed        int64
	attemptDelayBlocks uint64
}

func newSigningRetryLoop(
	logger log.StandardLogger,
	requestID signingRequestID,
	initialStartBlock uint64,
	signingGroupMemberIndex group.MemberIndex,
	signingGroupOperators chain.Addresses,
	chainConfig *ChainConfig,
	announcer signingAnnouncer,
) *signingRetryLoop {
	// Get the 8-byte seed needed for the random retry algorithm. We take the
	// first 8 bytes of signing request ID. This allows us to not care in this
	// piece of the code about the length of the ID and how this ID is proposed.
	attemptSeed := int64(binary.BigEndian.Uint64(requestID[:8]))

	return &signingRetryLoop{
		logger:                   logger,
		requestID:                requestID,
		signingGroupMemberIndex:  signingGroupMemberIndex,
		signingGroupOperators:    signingGroupOperators,
		chainConfig:              chainConfig,
		announcer:                announcer,
		announcementDelayBlocks:  1,
		announcementActiveBlocks: 5,
		attemptCounter:           0,
		attemptStartBlock:        initialStartBlock,
		attemptSeed:              attemptSeed,
		attemptDelayBlocks:       5,
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

// start begins the signing retry loop using the given signing attempt function.
// The retry loop terminates when the signing result is produced or the ctx
// parameter is done, whatever comes first.
func (srl *signingRetryLoop) start(
	ctx context.Context,
	waitForBlockFn waitForBlockFn,
	signingAttemptFn signingAttemptFn,
) (*signing.Result, error) {
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
		// signingAttemptMaxBlockDuration constant.
		//
		// For example, the attempt may fail at the end of the protocol but the
		// error is returned after some time and more blocks than expected are
		// mined in the meantime.
		if srl.attemptCounter > 1 {
			srl.attemptStartBlock = srl.attemptStartBlock +
				srl.announcementDelayBlocks +
				srl.announcementActiveBlocks +
				signingAttemptMaxBlockDuration +
				srl.attemptDelayBlocks
		}

		announcementStartBlock := srl.attemptStartBlock + srl.announcementDelayBlocks
		err := waitForBlockFn(ctx, announcementStartBlock)
		if err != nil {
			return nil, fmt.Errorf(
				"failed waiting for announcement start block [%v] "+
					"for attempt [%v]: [%v]",
				announcementStartBlock,
				srl.attemptCounter,
				err,
			)
		}

		// Set up the announcement phase stop signal.
		announceCtx, cancelAnnounceCtx := context.WithCancel(ctx)
		announcementEndBlock := announcementStartBlock + srl.announcementActiveBlocks
		go func() {
			defer cancelAnnounceCtx()

			if err := waitForBlockFn(ctx, announcementEndBlock); err != nil {
				srl.logger.Errorf(
					"[member:%v] failed waiting for announcement end "+
						"block [%v] for attempt [%v]: [%v]",
					srl.signingGroupMemberIndex,
					announcementEndBlock,
					srl.attemptCounter,
					err,
				)
			}
		}()

		srl.logger.Infof(
			"[member:%v] starting announcement phase for attempt [%v]",
			srl.signingGroupMemberIndex,
			srl.attemptCounter,
		)

		readyMembersIndexes, err := srl.announcer.Announce(
			announceCtx,
			srl.signingGroupMemberIndex,
			fmt.Sprintf("%s-%v", srl.requestID, srl.attemptCounter),
		)
		if err != nil {
			srl.logger.Warnf(
				"[member:%v] announcement for attempt [%v] "+
					"failed: [%v]; starting next attempt",
				srl.signingGroupMemberIndex,
				srl.attemptCounter,
				err,
			)
			continue
		}

		// Check the loop stop signal.
		if ctx.Err() != nil {
			return nil, nil
		}

		if len(readyMembersIndexes) >= srl.chainConfig.HonestThreshold {
			srl.logger.Infof(
				"[member:%v] completed announcement phase for attempt [%v] "+
					"with honest majority of [%v] members ready to sign",
				srl.signingGroupMemberIndex,
				srl.attemptCounter,
				len(readyMembersIndexes),
			)
		} else {
			srl.logger.Warnf(
				"[member:%v] completed announcement phase for attempt [%v] "+
					"with minority of [%v] members ready to sign; "+
					"starting next attempt",
				srl.signingGroupMemberIndex,
				srl.attemptCounter,
				len(readyMembersIndexes),
			)
			continue
		}

		excludedMembersIndexes, err := srl.performMembersSelection(
			readyMembersIndexes,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"cannot select members for attempt [%v]: [%w]",
				srl.attemptCounter,
				err,
			)
		}

		attemptSkipped := slices.Contains(
			excludedMembersIndexes,
			srl.signingGroupMemberIndex,
		)

		var result *signing.Result
		var attemptErr error

		if !attemptSkipped {
			result, attemptErr = signingAttemptFn(&signingAttemptParams{
				number:                 srl.attemptCounter,
				startBlock:             announcementEndBlock,
				excludedMembersIndexes: excludedMembersIndexes,
			})
		} else {
			srl.logger.Infof(
				"[member:%v] attempt [%v] skipped",
				srl.signingGroupMemberIndex,
				srl.attemptCounter,
			)
		}

		if attemptSkipped || attemptErr != nil {
			continue
		}

		return result, nil
	}
}

// performMembersSelection runs the member selection process whose result
// is a list of members' indexes that should be excluded by the client
// for the given signing attempt.
//
// The member selection process is done based on the list of ready members
// provided as the readyMembersIndexes argument. This list is used twice:
//
// First, the algorithm determining the qualified operators set uses the
// ready members list to build an input consisting of only active operators.
// This way we guarantee that the qualified operators set contains only
// ready and active operators that will actually take part in the signing
// attempt.
//
// Second, the ready members list is used to determine a list of excluded
// members. The excluded members list is built using the qualified operators
// set. The algorithm that determines the qualified operators set does not
// care about an exact mapping between operators and controlled members but
// relies on the members count solely. That means the information about
// readiness of specific members controlled by the given operators is not
// included in the resulting qualified operators set. In order to properly
// decide about inclusion or exclusion of specific members of a given
// qualified operator, we must take the ready members list into account.
func (srl *signingRetryLoop) performMembersSelection(
	readyMembersIndexes []group.MemberIndex,
) ([]group.MemberIndex, error) {
	qualifiedOperatorsSet, err := srl.qualifiedOperatorsSet(readyMembersIndexes)
	if err != nil {
		return nil, fmt.Errorf("cannot get qualified operators: [%w]", err)
	}

	// Exclude all members controlled by the operators that were not
	// qualified for the current attempt.
	return srl.excludedMembersIndexes(
		qualifiedOperatorsSet,
		readyMembersIndexes,
	), nil
}

// qualifiedOperatorsSet returns a set of operators qualified to participate
// in the given signing attempt. The set of qualified operators is taken
// from the set of active operators who announced readiness through
// their controlled signing group members.
func (srl *signingRetryLoop) qualifiedOperatorsSet(
	readyMembersIndexes []group.MemberIndex,
) (map[chain.Address]bool, error) {
	// The retry algorithm expects that we count retries from 0. Since
	// the first invocation of the algorithm will be for `attemptCounter == 1`
	// we need to subtract one while determining the number of the given retry.
	retryCount := srl.attemptCounter - 1

	var readySigningGroupOperators []chain.Address
	for _, memberIndex := range readyMembersIndexes {
		readySigningGroupOperators = append(
			readySigningGroupOperators,
			srl.signingGroupOperators[memberIndex-1],
		)
	}

	qualifiedOperators, err := retry.EvaluateRetryParticipantsForSigning(
		readySigningGroupOperators,
		srl.attemptSeed,
		retryCount,
		uint(srl.chainConfig.HonestThreshold),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"random operator selection failed: [%w]",
			err,
		)
	}

	return chain.Addresses(qualifiedOperators).Set(), nil
}

// excludedMembersIndexes returns a list of excluded members' indexes for
// the given qualified operators set.
func (srl *signingRetryLoop) excludedMembersIndexes(
	qualifiedOperatorsSet map[chain.Address]bool,
	readyMembersIndexes []group.MemberIndex,
) []group.MemberIndex {
	includedMembersIndexes := make([]group.MemberIndex, 0)
	excludedMembersIndexes := make([]group.MemberIndex, 0)
	for i, operator := range srl.signingGroupOperators {
		memberIndex := group.MemberIndex(i + 1)

		if qualifiedOperatorsSet[operator] &&
			slices.Contains(readyMembersIndexes, memberIndex) {
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
			srl.attemptSeed + int64(srl.attemptCounter),
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

	return excludedMembersIndexes
}
