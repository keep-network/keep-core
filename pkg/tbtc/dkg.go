package tbtc

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"sort"
	"time"

	"go.uber.org/zap"
	"golang.org/x/exp/slices"

	"github.com/ipfs/go-log/v2"
	"github.com/keep-network/keep-common/pkg/persistence"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/generator"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/protocol/announcer"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/tecdsa/dkg"
	"github.com/keep-network/keep-core/pkg/tecdsa/retry"
)

// TODO: Revisit those constants, especially dkgResultSubmissionDelayStep
// which should be bigger once the contract integration is ready.
const (
	// dkgAttemptMaxBlockDuration determines the maximum block duration of a
	// single DKG attempt.
	dkgAttemptMaxBlockDuration = 150
	// dkgResultSubmissionDelayStep determines the delay step that is used to
	// calculate the submission delay time that should be respected by the
	// given member to avoid all members submitting the same DKG result at the
	// same time.
	dkgResultSubmissionDelayStep = 10 * time.Second
)

// dkgExecutor is a component responsible for the full execution of ECDSA
// Distributed Key Generation: determining members selected to the signing
// group, executing off-chain protocol, and publishing the result to the chain.
type dkgExecutor struct {
	chain          Chain
	netProvider    net.Provider
	walletRegistry *walletRegistry
	protocolLatch  *generator.ProtocolLatch

	// waitForBlockFn is a function used to wait for the given block.
	waitForBlockFn waitForBlockFn

	tecdsaExecutor *dkg.Executor
}

// newDkgExecutor creates a new instance of dkgExecutor struct. There should
// be only one instance of dkgExecutor.
func newDkgExecutor(
	chain Chain,
	netProvider net.Provider,
	walletRegistry *walletRegistry,
	protocolLatch *generator.ProtocolLatch,
	config Config,
	workPersistence persistence.BasicHandle,
	scheduler *generator.Scheduler,
) *dkgExecutor {
	tecdsaExecutor := dkg.NewExecutor(
		logger,
		scheduler,
		workPersistence,
		config.PreParamsPoolSize,
		config.PreParamsGenerationTimeout,
		config.PreParamsGenerationDelay,
		config.PreParamsGenerationConcurrency,
		config.KeyGenerationConcurrency,
	)

	return &dkgExecutor{
		chain:          chain,
		netProvider:    netProvider,
		walletRegistry: walletRegistry,
		protocolLatch:  protocolLatch,
		tecdsaExecutor: tecdsaExecutor,
	}
}

// preParamsCount returns the current count of the ECDSA DKG pre-parameters.
func (de *dkgExecutor) preParamsCount() int {
	return de.tecdsaExecutor.PreParamsCount()
}

// executeDkgIfEligible is the main function of dkgExecutor. It performs the
// full execution of ECDSA Distributed Key Generation: determining members
// selected to the signing group, executing off-chain protocol, and publishing
// the result to the chain.
func (de *dkgExecutor) executeDkgIfEligible(
	seed *big.Int,
	startBlockNumber uint64,
) {
	dkgLogger := logger.With(
		zap.String("seed", fmt.Sprintf("0x%x", seed)),
	)

	dkgLogger.Info("checking eligibility for DKG")
	memberIndexes, selectedSigningGroupOperators, err := de.checkEligibility(dkgLogger, seed)
	if err != nil {
		dkgLogger.Errorf("could not check eligibility for DKG: [%v]", err)
		return
	}

	if membersCount := len(memberIndexes); membersCount > 0 {
		if preParamsCount := de.tecdsaExecutor.PreParamsCount(); membersCount > preParamsCount {
			dkgLogger.Infof(
				"cannot join DKG as pre-parameters pool size is "+
					"too small; [%v] pre-parameters are required but "+
					"only [%v] available",
				membersCount,
				preParamsCount,
			)
			return
		}

		dkgLogger.Infof(
			"joining DKG and controlling [%v] group members",
			membersCount,
		)

		de.generateSigningGroup(
			dkgLogger,
			seed,
			memberIndexes,
			selectedSigningGroupOperators,
			startBlockNumber,
		)
	} else {
		dkgLogger.Infof("not eligible for DKG")
	}
}

// checkEligibility performs on-chain group selection and returns two pieces
// of information:
// - Indexes of members selected to the signing group and controlled by this
//   operator. The indexes are in range [1, `groupSize`]. The slice is nil if
//   none of the selected signing group members is controlled by this operator.
// - Addresses of all signing group members. There are always `groupSize`
//   elements in this slice.
func (de *dkgExecutor) checkEligibility(
	dkgLogger log.StandardLogger,
	seed *big.Int,
) ([]uint8, chain.Addresses, error) {
	selectedSigningGroupOperators, err := de.chain.SelectGroup(seed)
	if err != nil {
		return nil, nil, fmt.Errorf("selecting group not possible: [%v]", err)
	}

	dkgLogger.Infof("selected group members for DKG = %s", selectedSigningGroupOperators)

	if len(selectedSigningGroupOperators) > de.chain.GetConfig().GroupSize {
		return nil, nil, fmt.Errorf(
			"group size larger than supported: [%v]",
			len(selectedSigningGroupOperators),
		)
	}

	_, operatorPublicKey, err := de.chain.OperatorKeyPair()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get operator public key: [%v]", err)
	}

	operatorAddress, err := de.chain.Signing().PublicKeyToAddress(operatorPublicKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get operator address: [%v]", err)
	}

	indexes := make([]uint8, 0)
	for index, operator := range selectedSigningGroupOperators {
		// See if we are amongst those chosen
		if operator == operatorAddress {
			// The group member index should be in range [1, groupSize] so we
			// need to add 1.
			indexes = append(indexes, uint8(index)+1)
		}
	}

	return indexes, selectedSigningGroupOperators, nil
}

// setupBroadcastChannel creates and initializes broadcast channel for the
// current DKG execution. It is a temporary channel named after the seed and
// the protocol name.
func (de *dkgExecutor) setupBroadcastChannel(
	seed *big.Int,
	membershipValidator *group.MembershipValidator,
) (net.BroadcastChannel, error) {
	// Create temporary broadcast channel name for DKG using the
	// group selection seed with the protocol name as prefix.
	channelName := fmt.Sprintf("%s-%s", ProtocolName, seed.Text(16))

	broadcastChannel, err := de.netProvider.BroadcastChannelFor(channelName)
	if err != nil {
		return nil, fmt.Errorf("failed to get broadcast channel: [%v]", err)
	}

	dkg.RegisterUnmarshallers(broadcastChannel)
	registerStopPillUnmarshaller(broadcastChannel)

	err = broadcastChannel.SetFilter(membershipValidator.IsInGroup)
	if err != nil {
		return nil, fmt.Errorf(
			"could not set filter for channel [%v]: [%v]",
			broadcastChannel.Name(),
			err,
		)
	}

	return broadcastChannel, nil
}

// generateSigningGroup executes off-chain protocol for each member controlled
// by the current operator and upon successful execution of the protocol
// publishes the result to the chain.
func (de *dkgExecutor) generateSigningGroup(
	dkgLogger *zap.SugaredLogger,
	seed *big.Int,
	memberIndexes []uint8,
	selectedSigningGroupOperators chain.Addresses,
	startBlockNumber uint64,
) {
	membershipValidator := group.NewMembershipValidator(
		dkgLogger,
		selectedSigningGroupOperators,
		de.chain.Signing(),
	)

	broadcastChannel, err := de.setupBroadcastChannel(seed, membershipValidator)
	if err != nil {
		dkgLogger.Errorf("could not set up a broadcast channel: [%v]", err)
		return
	}

	chainConfig := de.chain.GetConfig()

	for _, index := range memberIndexes {
		// Capture the member index for the goroutine.
		memberIndex := index

		go func() {
			de.protocolLatch.Lock()
			defer de.protocolLatch.Unlock()

			announcer := announcer.New(
				fmt.Sprintf("%v-%v", ProtocolName, "dkg"),
				broadcastChannel,
				membershipValidator,
			)

			retryLoop := newDkgRetryLoop(
				dkgLogger,
				seed,
				startBlockNumber,
				memberIndex,
				selectedSigningGroupOperators,
				chainConfig,
				announcer,
			)

			// TODO: For this client iteration, the retry loop is started
			//       with a 168h timeout and a stop pill sent by any group
			//       member. Once the WalletRegistry is integrated, the stop
			//       signal should be generated by observing the DKG result
			//       submission or timeout.
			loopCtx, cancelLoopCtx := context.WithTimeout(
				context.Background(),
				7*24*time.Hour,
			)
			defer cancelLoopCtx()
			cancelDkgContextOnStopSignal(
				loopCtx,
				cancelLoopCtx,
				broadcastChannel,
				seed.Text(16),
			)

			result, err := retryLoop.start(
				loopCtx,
				de.waitForBlockFn,
				func(attempt *dkgAttemptParams) (*dkg.Result, error) {
					dkgAttemptLogger := dkgLogger.With(
						zap.Uint("attempt", attempt.number),
						zap.Uint64("attemptStartBlock", attempt.startBlock),
					)

					dkgAttemptLogger.Infof(
						"[member:%v] scheduled dkg attempt "+
							"with [%v] group members (excluded: [%v])",
						memberIndex,
						chainConfig.GroupSize-len(attempt.excludedMembersIndexes),
						attempt.excludedMembersIndexes,
					)

					// Set up the attempt timeout signal.
					attemptCtx, cancelAttemptCtx := context.WithCancel(
						loopCtx,
					)
					go func() {
						defer cancelAttemptCtx()

						err := de.waitForBlockFn(
							loopCtx,
							attempt.startBlock+dkgAttemptMaxBlockDuration,
						)
						if err != nil {
							dkgAttemptLogger.Warnf(
								"[member:%v] failed waiting for "+
									"attempt stop signal: [%v]",
								memberIndex,
								err,
							)
						}
					}()

					// sessionID must be different for each attempt.
					sessionID := fmt.Sprintf(
						"%v-%v",
						seed.Text(16),
						attempt.number,
					)

					result, err := de.tecdsaExecutor.Execute(
						attemptCtx,
						dkgAttemptLogger,
						seed,
						sessionID,
						memberIndex,
						chainConfig.GroupSize,
						chainConfig.DishonestThreshold(),
						attempt.excludedMembersIndexes,
						broadcastChannel,
						membershipValidator,
					)
					if err != nil {
						dkgAttemptLogger.Errorf(
							"[member:%v] dkg attempt failed: [%v]",
							memberIndex,
							err,
						)

						return nil, err
					}

					// Schedule the stop pill to be sent a fixed amount of
					// time after the result is returned. Do not do it
					// immediately as other members can be very close
					// to produce the result as well. This mechanism should
					// be more sophisticated but since it is temporary, we
					// can live with it for now.
					go func() {
						time.Sleep(1 * time.Minute)
						if err := sendDkgStopPill(
							loopCtx,
							broadcastChannel,
							seed.Text(16),
							attempt.number,
						); err != nil {
							dkgLogger.Errorf(
								"[member:%v] could not send the stop pill: [%v]",
								memberIndex,
								err,
							)
						}
					}()

					return result, nil
				},
			)
			if err != nil {
				dkgLogger.Errorf(
					"[member:%v] failed to execute dkg: [%v]",
					memberIndex,
					err,
				)
				return
			}
			// TODO: This condition should go away once we integrate
			// WalletRegistry contract. In this scenario, member received
			// a StopPill from some other group member and it means that
			// the result has been produced but the current member did not
			// participate in the work so they do not know the result.
			if result == nil {
				dkgLogger.Infof(
					"[member:%v] dkg retry loop received stop signal",
					memberIndex,
				)
				return
			}

			// TODO: Snapshot the key material before doing on-chain result
			//       submission.

			operatingMemberIndexes := result.Group.OperatingMemberIDs()
			dkgResultChannel := make(chan *DKGResultSubmittedEvent)

			dkgResultSubscription := de.chain.OnDKGResultSubmitted(
				func(event *DKGResultSubmittedEvent) {
					dkgResultChannel <- event
				},
			)
			defer dkgResultSubscription.Unsubscribe()

			// Set up the publication stop signal that should allow to
			// perform all the result-signing-related actions and
			// handle the worst case when the result is submitted by the
			// last group member.
			publicationTimeout := time.Duration(chainConfig.GroupSize) *
				dkgResultSubmissionDelayStep
			publicationCtx, cancelPublicationCtx := context.WithTimeout(
				context.Background(),
				publicationTimeout,
			)
			// TODO: Call cancelPublicationCtx() when the result is
			//       available and published and remove this goroutine.
			//       This goroutine is duplicating context.WithTimeout work
			//       right now but is here to emphasize the need of manual
			//       context cancellation.
			go func() {
				defer cancelPublicationCtx()
				time.Sleep(publicationTimeout)
			}()

			err = dkg.Publish(
				publicationCtx,
				dkgLogger,
				seed.Text(16),
				memberIndex,
				broadcastChannel,
				membershipValidator,
				newDkgResultSigner(de.chain),
				newDkgResultSubmitter(dkgLogger, de.chain),
				result,
			)
			if err != nil {
				// Result publication failed. It means that either the result
				// this member proposed is not supported by the majority of
				// group members or that the chain interaction failed.
				// In either case, we observe the chain for the result
				// published by any other group member and based on that,
				// we decide whether we should stay in the final group or
				// drop our membership.
				dkgLogger.Warnf(
					"[member:%v] DKG result publication process failed [%v]",
					memberIndex,
					err,
				)

				if operatingMemberIndexes, err = decideSigningGroupMemberFate(
					publicationCtx,
					memberIndex,
					dkgResultChannel,
					result,
				); err != nil {
					dkgLogger.Errorf(
						"[member:%v] failed to handle DKG result "+
							"publishing failure: [%v]",
						memberIndex,
						err,
					)
					return
				}
			}

			// Final signing group may differ from the original DKG
			// group outputted by the sortition protocol. One need to
			// determine the final signing group based on the selected
			// group members who behaved correctly during DKG protocol.
			finalSigningGroupOperators, finalSigningGroupMembersIndexes, err :=
				finalSigningGroup(
					selectedSigningGroupOperators,
					operatingMemberIndexes,
					chainConfig,
				)
			if err != nil {
				dkgLogger.Errorf(
					"[member:%v] failed to resolve final signing "+
						"group: [%v]",
					memberIndex,
					err,
				)
				return
			}

			// Just like the final and original group may differ, the
			// member index used during the DKG protocol may differ
			// from the final signing group member index as well.
			// We need to remap it.
			finalSigningGroupMemberIndex, ok :=
				finalSigningGroupMembersIndexes[memberIndex]
			if !ok {
				dkgLogger.Errorf(
					"[member:%v] failed to resolve final signing "+
						"group member index",
					memberIndex,
				)
				return
			}

			signer := newSigner(
				result.PrivateKeyShare.PublicKey(),
				finalSigningGroupOperators,
				finalSigningGroupMemberIndex,
				result.PrivateKeyShare,
			)

			err = de.walletRegistry.registerSigner(signer)
			if err != nil {
				dkgLogger.Errorf(
					"failed to register %s: [%v]",
					signer,
					err,
				)
				return
			}

			dkgLogger.Infof("registered %s", signer)
		}()
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
	ctx context.Context,
	memberIndex group.MemberIndex,
	dkgResultChannel chan *DKGResultSubmittedEvent,
	result *dkg.Result,
) ([]group.MemberIndex, error) {
	select {
	case dkgResultEvent := <-dkgResultChannel:
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
	case <-ctx.Done():
		return nil, fmt.Errorf("result publication timed out")
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
	ctx context.Context,
	memberIndex group.MemberIndex,
	result *dkg.Result,
	signatures map[group.MemberIndex][]byte,
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

	submissionDelay := time.Duration(memberIndex-1) * dkgResultSubmissionDelayStep

	drs.dkgLogger.Infof(
		"[member:%v] waiting [%v] to submit",
		memberIndex,
		submissionDelay,
	)

	submissionTimer := time.NewTimer(submissionDelay)
	defer submissionTimer.Stop()

	for {
		select {
		case <-submissionTimer.C:
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
					"[%v] supporting member signatures",
				memberIndex,
				publicKeyBytes,
				len(signatures),
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
		case <-ctx.Done():
			return fmt.Errorf("result publication timed out")
		}
	}
}
