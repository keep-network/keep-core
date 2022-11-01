package tbtc

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"golang.org/x/exp/slices"
	"math/big"
	"sort"

	"github.com/ipfs/go-log/v2"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/tecdsa/dkg"
	"github.com/keep-network/keep-core/pkg/tecdsa/retry"
)

// dkgAnnouncer represents a component responsible for exchanging readiness
// announcements for the given DKG attempt for the given seed.
type dkgAnnouncer interface {
	Announce(
		ctx context.Context,
		memberIndex group.MemberIndex,
		sessionID string,
	) ([]group.MemberIndex, error)
}

// dkgRetryLoop is a struct that encapsulates the DKG retry logic.
type dkgRetryLoop struct {
	logger log.StandardLogger

	// seed is the original seed for DKG.
	// Used for the announcement. It never changes.
	seed *big.Int

	memberIndex       group.MemberIndex
	selectedOperators chain.Addresses

	chainConfig *ChainConfig

	announcer                dkgAnnouncer
	announcementDelayBlocks  uint64
	announcementActiveBlocks uint64

	attemptCounter     uint
	attemptStartBlock  uint64
	// attemptSeed is a 8-byte seed obtained from the original seed.
	// Used for the random operator selection. It never changes.
	attemptSeed        int64
	attemptDelayBlocks uint64
}

func newDkgRetryLoop(
	logger log.StandardLogger,
	seed *big.Int,
	initialStartBlock uint64,
	memberIndex group.MemberIndex,
	selectedOperators chain.Addresses,
	chainConfig *ChainConfig,
	announcer dkgAnnouncer,
) *dkgRetryLoop {
	// Compute the 8-byte seed needed for the random retry algorithm. We take
	// the first 8 bytes of the hash of the DKG seed. This allows us to not
	// care in this piece of the code about the length of the seed and how this
	// seed is proposed.
	seedSha256 := sha256.Sum256(seed.Bytes())
	attemptSeed := int64(binary.BigEndian.Uint64(seedSha256[:8]))

	return &dkgRetryLoop{
		logger:                   logger,
		seed:                     seed,
		memberIndex:              memberIndex,
		selectedOperators:        selectedOperators,
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

// dkgAttemptParams represents parameters of a DKG attempt.
type dkgAttemptParams struct {
	number                 uint
	startBlock             uint64
	excludedMembersIndexes []group.MemberIndex
}

// dkgAttemptFn represents a function performing a DKG attempt.
type dkgAttemptFn func(*dkgAttemptParams) (*dkg.Result, uint64, error)

// start begins the DKG retry loop using the given DKG attempt function.
// The retry loop terminates when the DKG result is produced or the ctx
// parameter is done, whatever comes first.
func (drl *dkgRetryLoop) start(
	ctx context.Context,
	waitForBlockFn waitForBlockFn,
	dkgAttemptFn dkgAttemptFn,
) (*dkg.Result, uint64, error) {
	for {
		drl.attemptCounter++

		// In order to start attempts >1 in the right place, we need to
		// determine how many blocks were taken by previous attempts. We assume
		// the worst case that each attempt failed at the end of the DKG
		// protocol.
		//
		// That said, we need to increment the previous attempt start
		// block by the number of blocks equal to the protocol duration and
		// by some additional delay blocks. We need a small fixed delay in
		// order to mitigate all corner cases where the actual attempt duration
		// was slightly longer than the expected duration determined by the
		// dkg.ProtocolBlocks function.
		//
		// For example, the attempt may fail at the end of the protocol but the
		// error is returned after some time and more blocks than expected are
		// mined in the meantime.
		if drl.attemptCounter > 1 {
			drl.attemptStartBlock = drl.attemptStartBlock +
				drl.announcementDelayBlocks +
				drl.announcementActiveBlocks +
				dkg.ProtocolBlocks() +
				drl.attemptDelayBlocks
		}

		announcementStartBlock := drl.attemptStartBlock + drl.announcementDelayBlocks
		err := waitForBlockFn(ctx, announcementStartBlock)
		if err != nil {
			return nil, 0, fmt.Errorf(
				"failed waiting for announcement start block [%v] "+
					"for attempt [%v]: [%v]",
				announcementStartBlock,
				drl.attemptCounter,
				err,
			)
		}

		// Set up the announcement phase stop signal.
		announceCtx, cancelAnnounceCtx := context.WithCancel(ctx)
		announcementEndBlock := announcementStartBlock + drl.announcementActiveBlocks
		go func() {
			defer cancelAnnounceCtx()

			if err := waitForBlockFn(ctx, announcementEndBlock); err != nil {
				drl.logger.Errorf(
					"[member:%v] failed waiting for announcement end "+
						"block [%v] for attempt [%v]: [%v]",
					drl.memberIndex,
					announcementEndBlock,
					drl.attemptCounter,
					err,
				)
			}
		}()

		drl.logger.Infof(
			"[member:%v] starting announcement phase for attempt [%v]",
			drl.memberIndex,
			drl.attemptCounter,
		)

		readyMembersIndexes, err := drl.announcer.Announce(
			announceCtx,
			drl.memberIndex,
			fmt.Sprintf("%v-%v", drl.seed, drl.attemptCounter),
		)
		if err != nil {
			drl.logger.Warnf(
				"[member:%v] announcement for attempt [%v] "+
					"failed: [%v]; starting next attempt",
				drl.memberIndex,
				drl.attemptCounter,
				err,
			)
			continue
		}

		// Check the loop stop signal.
		if ctx.Err() != nil {
			return nil, 0, nil
		}

		if len(readyMembersIndexes) >= drl.chainConfig.GroupQuorum {
			drl.logger.Infof(
				"[member:%v] completed announcement phase for attempt [%v] "+
					"with quorum of [%v] members ready to perform DKG",
				drl.memberIndex,
				drl.attemptCounter,
				len(readyMembersIndexes),
			)
		} else {
			drl.logger.Warnf(
				"[member:%v] completed announcement phase for attempt [%v] "+
					"with non-quorum of [%v] members ready to perform DKG; "+
					"starting next attempt",
				drl.memberIndex,
				drl.attemptCounter,
				len(readyMembersIndexes),
			)
			continue
		}

		excludedMembersIndexes, err := drl.performMembersSelection(
			readyMembersIndexes,
		)
		if err != nil {
			return nil, 0, fmt.Errorf(
				"cannot select members for attempt [%v]: [%w]",
				drl.attemptCounter,
				err,
			)
		}

		attemptSkipped := slices.Contains(
			excludedMembersIndexes,
			drl.memberIndex,
		)

		var result *dkg.Result
		var executionEndBlock uint64
		var attemptErr error

		if !attemptSkipped {
			result, executionEndBlock, attemptErr = dkgAttemptFn(&dkgAttemptParams{
				number:                 drl.attemptCounter,
				startBlock:             announcementEndBlock,
				excludedMembersIndexes: excludedMembersIndexes,
			})
		} else {
			drl.logger.Infof(
				"[member:%v] attempt [%v] skipped",
				drl.memberIndex,
				drl.attemptCounter,
			)
		}

		if attemptSkipped || attemptErr != nil {
			continue
		}

		return result, executionEndBlock, nil
	}
}

// performMembersSelection runs the member selection process whose result
// is a list of members' indexes that should be excluded by the client
// for the given DKG attempt.
//
// The member selection process is done based on the list of ready members
// provided as the readyMembersIndexes argument. This list is used twice:
//
// First, the algorithm determining the qualified operators set uses the
// ready members list to build an input consisting of only active operators.
// This way we guarantee that the qualified operators set contains only
// ready and active operators that will actually take part in the DKG
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
func (drl *dkgRetryLoop) performMembersSelection(
	readyMembersIndexes []group.MemberIndex,
) ([]group.MemberIndex, error) {
	qualifiedOperatorsSet, err := drl.qualifiedOperatorsSet(readyMembersIndexes)
	if err != nil {
		return nil, fmt.Errorf("cannot get qualified operators: [%w]", err)
	}

	excludedMembersIndexes := make([]group.MemberIndex, 0)
	for i, operator := range drl.selectedOperators {
		memberIndex := group.MemberIndex(i + 1)

		included := qualifiedOperatorsSet[operator] &&
			slices.Contains(readyMembersIndexes, memberIndex)

		if !included {
			excludedMembersIndexes = append(
				excludedMembersIndexes,
				memberIndex,
			)
		}
	}

	return excludedMembersIndexes, nil
}

// qualifiedOperatorsSet returns a set of operators qualified to participate
// in the given DKG attempt. The set of qualified operators is taken from the
// set of active operators who announced readiness through their controlled DKG
// group members.
func (drl *dkgRetryLoop) qualifiedOperatorsSet(
	readyMembersIndexes []group.MemberIndex,
) (map[chain.Address]bool, error) {
	var readyOperators chain.Addresses
	for _, memberIndex := range readyMembersIndexes {
		readyOperators = append(
			readyOperators,
			drl.selectedOperators[memberIndex-1],
		)
	}

	// For the first attempt, just return the operators who announced readiness.
	// Otherwise, randomly exclude operators from the ready operators set.
	if drl.attemptCounter == 1 {
		return readyOperators.Set(), nil
	}

	// The retry algorithm expects that we count retries from 0. Since
	// the first invocation of the algorithm will be for `attemptCounter == 1`
	// we need to subtract one while determining the number of the given retry.
	retryCount := drl.attemptCounter - 1

	qualifiedOperators, err := retry.EvaluateRetryParticipantsForKeyGeneration(
		readyOperators,
		drl.attemptSeed,
		retryCount,
		uint(drl.chainConfig.GroupQuorum),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"random operator selection failed: [%w]",
			err,
		)
	}

	return chain.Addresses(qualifiedOperators).Set(), nil
}

// decideSigningGroupMemberFate decides what the member will do in case it
// failed to publish its DKG result. Member can stay in the group if it supports
// the same group public key as the one registered on-chain and the member is
// not considered as misbehaving by the group.
func decideSigningGroupMemberFate(
	memberIndex group.MemberIndex,
	dkgResultChannel chan *DKGResultSubmittedEvent,
	publicationStartBlock uint64,
	result *dkg.Result,
	chainConfig *ChainConfig,
	blockCounter chain.BlockCounter,
) ([]group.MemberIndex, error) {
	dkgResultEvent, err := waitForDkgResultEvent(
		dkgResultChannel,
		publicationStartBlock,
		chainConfig,
		blockCounter,
	)
	if err != nil {
		return nil, err
	}

	groupPublicKeyBytes, err := result.GroupPublicKeyBytes()
	if err != nil {
		return nil, err
	}

	// If member doesn't support the same group public key, it could not stay
	// in the group.
	if !bytes.Equal(groupPublicKeyBytes, dkgResultEvent.GroupPublicKeyBytes) {
		return nil, fmt.Errorf(
			"[member:%v] could not stay in the group because "+
				"the member does not support the same group public key",
			memberIndex,
		)
	}

	misbehavedSet := make(map[group.MemberIndex]struct{})
	for _, misbehavedID := range dkgResultEvent.Misbehaved {
		misbehavedSet[misbehavedID] = struct{}{}
	}

	// If member is considered as misbehaved, it could not stay in the group.
	if _, isMisbehaved := misbehavedSet[memberIndex]; isMisbehaved {
		return nil, fmt.Errorf(
			"[member:%v] could not stay in the group because "+
				"the member is considered as misbehaving",
			memberIndex,
		)
	}

	// Construct a new view of the operating members according to the accepted
	// DKG result.
	operatingMemberIndexes := make([]group.MemberIndex, 0)
	for _, memberID := range result.Group.MemberIDs() {
		if _, isMisbehaved := misbehavedSet[memberID]; !isMisbehaved {
			operatingMemberIndexes = append(operatingMemberIndexes, memberID)
		}
	}

	return operatingMemberIndexes, nil
}

// waitForDkgResultEvent waits for the DKG result submission event. It times out
// and returns error if the DKG result event is not emitted on time.
func waitForDkgResultEvent(
	dkgResultChannel chan *DKGResultSubmittedEvent,
	publicationStartBlock uint64,
	chainConfig *ChainConfig,
	blockCounter chain.BlockCounter,
) (*DKGResultSubmittedEvent, error) {
	timeoutBlock := publicationStartBlock + dkg.PrePublicationBlocks() +
		(uint64(chainConfig.GroupSize) * chainConfig.ResultPublicationBlockStep)

	timeoutBlockChannel, err := blockCounter.BlockHeightWaiter(timeoutBlock)
	if err != nil {
		return nil, err
	}

	select {
	case dkgResultEvent := <-dkgResultChannel:
		return dkgResultEvent, nil
	case <-timeoutBlockChannel:
		return nil, fmt.Errorf("ECDSA DKG result publication timed out")
	}
}

// finalSigningGroup takes three parameters:
//   - selectedOperators: Contains addresses of all selected operators. Slice
//     length equals to the groupSize. Each element with index N corresponds
//     to the group member with ID N+1.
//   - operatingMembersIndexes: Contains group members indexes that were neither
//     disqualified nor marked as inactive. Slice length is lesser than or equal
//     to the groupSize.
//   - chainConfig: The tBTC chain's configuration
//
// Using those parameters, this function transforms the selectedOperators
// slice into another slice that contains addresses of all operators
// that were neither disqualified nor marked as inactive. This way, the
// resulting slice has only addresses of properly operating operators
// who form the resulting group.
//
// Apart from that, this function returns a map that holds the final signing
// group members indexes that should be used by particular members who behaved
// correctly during the DKG protocol execution. The key of this map is the
// member index used during DKG protocol and the value is the new member
// index that should be used in the context of the final signing group.
//
// Example:
// selectedOperators: [0xAA, 0xBB, 0xCC, 0xDD, 0xEE]
// operatingMembersIndexes: [5, 1, 3]
// finalOperators: [0xAA, 0xCC, 0xEE]
// finalMembersIndexes: [1:1, 3:2, 5:3]
func finalSigningGroup(
	selectedOperators []chain.Address,
	operatingMembersIndexes []group.MemberIndex,
	chainConfig *ChainConfig,
) (
	[]chain.Address,
	map[group.MemberIndex]group.MemberIndex,
	error,
) {
	if len(selectedOperators) != chainConfig.GroupSize ||
		len(operatingMembersIndexes) < chainConfig.GroupQuorum {
		return nil, nil, fmt.Errorf("invalid input parameters")
	}

	sort.Slice(operatingMembersIndexes, func(i, j int) bool {
		return operatingMembersIndexes[i] < operatingMembersIndexes[j]
	})

	finalOperators := make(
		[]chain.Address,
		len(operatingMembersIndexes),
	)
	finalMembersIndexes := make(
		map[group.MemberIndex]group.MemberIndex,
		len(operatingMembersIndexes),
	)

	for i, operatingMemberID := range operatingMembersIndexes {
		finalOperators[i] = selectedOperators[operatingMemberID-1]
		finalMembersIndexes[operatingMemberID] = group.MemberIndex(i + 1)
	}

	return finalOperators, finalMembersIndexes, nil
}

// dkgResultSigner is responsible for signing the DKG result and verification of
// signatures generated by other group members.
type dkgResultSigner struct {
	chain Chain
}

func newDkgResultSigner(chain Chain) *dkgResultSigner {
	return &dkgResultSigner{
		chain: chain,
	}
}

// SignResult signs the provided DKG result. It returns the information
// pertaining to the signing process: public key, signature, result hash.
func (drs *dkgResultSigner) SignResult(result *dkg.Result) (*dkg.SignedResult, error) {
	resultHash, err := drs.chain.CalculateDKGResultHash(result)
	if err != nil {
		return nil, fmt.Errorf(
			"dkg result hash calculation failed [%w]",
			err,
		)
	}

	signing := drs.chain.Signing()

	signature, err := signing.Sign(resultHash[:])
	if err != nil {
		return nil, fmt.Errorf(
			"dkg result hash signing failed [%w]",
			err,
		)
	}

	return &dkg.SignedResult{
		PublicKey:  signing.PublicKey(),
		Signature:  signature,
		ResultHash: resultHash,
	}, nil
}

// VerifySignature verifies if the signature was generated from the provided
// DKG result has using the provided public key.
func (drs *dkgResultSigner) VerifySignature(signedResult *dkg.SignedResult) (bool, error) {
	return drs.chain.Signing().VerifyWithPublicKey(
		signedResult.ResultHash[:],
		signedResult.Signature,
		signedResult.PublicKey,
	)
}

// dkgResultSubmitter is responsible for submitting the DKG result to the chain.
type dkgResultSubmitter struct {
	dkgLogger log.StandardLogger
	chain     Chain
}

func newDkgResultSubmitter(
	dkgLogger log.StandardLogger,
	chain Chain,
) *dkgResultSubmitter {
	return &dkgResultSubmitter{
		dkgLogger: dkgLogger,
		chain:     chain,
	}
}

// SubmitResult submits the DKG result along with submitting signatures to the
// chain. In the process, it checks if the number of signatures is above
// the required threshold, whether the result was already submitted and waits
// until the member is eligible for DKG result submission.
func (drs *dkgResultSubmitter) SubmitResult(
	memberIndex group.MemberIndex,
	result *dkg.Result,
	signatures map[group.MemberIndex][]byte,
	startBlockNumber uint64,
) error {
	config := drs.chain.GetConfig()

	if len(signatures) < config.GroupQuorum {
		return fmt.Errorf(
			"could not submit result with [%v] signatures for group quorum [%v]",
			len(signatures),
			config.GroupQuorum,
		)
	}

	resultSubmittedChan := make(chan uint64)

	subscription := drs.chain.OnDKGResultSubmitted(
		func(event *DKGResultSubmittedEvent) {
			resultSubmittedChan <- event.BlockNumber
		},
	)
	defer subscription.Unsubscribe()

	dkgState, err := drs.chain.GetDKGState()
	if err != nil {
		return fmt.Errorf("could not check DKG state: [%w]", err)
	}

	if dkgState != AwaitingResult {
		// Someone who was ahead of us in the queue submitted the result. Giving up.
		drs.dkgLogger.Infof(
			"[member:%v] DKG is no longer awaiting the result; "+
				"aborting DKG result submission",
			memberIndex,
		)
		return nil
	}

	// Wait until the current member is eligible to submit the result.
	submitterEligibleChan, err := drs.setupEligibilityQueue(
		startBlockNumber,
		memberIndex,
	)
	if err != nil {
		return fmt.Errorf("cannot set up eligibility queue: [%w]", err)
	}

	for {
		select {
		case blockNumber := <-submitterEligibleChan:
			// Member becomes eligible to submit the result. Result submission
			// would trigger the sender side of the result submission event
			// listener but also cause the receiver side (this select)
			// termination that will result with a dangling goroutine blocked
			// forever on the `onSubmittedResultChan` channel. This would
			// cause a resource leak. In order to avoid that, we should
			// unsubscribe from the result submission event listener before
			// submitting the result.
			subscription.Unsubscribe()

			publicKeyBytes, err := result.GroupPublicKeyBytes()
			if err != nil {
				return fmt.Errorf("cannot get public key bytes [%w]", err)
			}

			drs.dkgLogger.Infof(
				"[member:%v] submitting DKG result with public key [0x%x] and "+
					"[%v] supporting member signatures at block [%v]",
				memberIndex,
				publicKeyBytes,
				len(signatures),
				blockNumber,
			)

			return drs.chain.SubmitDKGResult(
				memberIndex,
				result,
				signatures,
			)
		case blockNumber := <-resultSubmittedChan:
			drs.dkgLogger.Infof(
				"[member:%v] leaving; DKG result submitted by other member "+
					"at block [%v]",
				memberIndex,
				blockNumber,
			)
			// A result has been submitted by other member. Leave without
			// publishing the result.
			return nil
		}
	}
}

// setupEligibilityQueue waits until the current member is eligible to
// submit a result to the blockchain. First member is eligible to submit straight
// away, each following member is eligible after pre-defined block step.
//
// TODO: Revisit the setupEligibilityQueue function. The RFC mentions we should
// start submitting from a random member, not the first one.
func (drs *dkgResultSubmitter) setupEligibilityQueue(
	startBlockNumber uint64,
	memberIndex group.MemberIndex,
) (<-chan uint64, error) {
	blockWaitTime := (uint64(memberIndex) - 1) *
		drs.chain.GetConfig().ResultPublicationBlockStep

	eligibleBlockHeight := startBlockNumber + blockWaitTime

	drs.dkgLogger.Infof(
		"[member:%v] waiting for block [%v] to submit",
		memberIndex,
		eligibleBlockHeight,
	)

	blockCounter, err := drs.chain.BlockCounter()
	if err != nil {
		return nil, fmt.Errorf("could not get block counter [%w]", err)
	}

	waiter, err := blockCounter.BlockHeightWaiter(eligibleBlockHeight)
	if err != nil {
		return nil, fmt.Errorf("block height waiter failure [%w]", err)
	}

	return waiter, err
}
