package tbtc

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"time"

	"go.uber.org/zap"

	"github.com/ipfs/go-log/v2"
	"github.com/keep-network/keep-common/pkg/persistence"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/generator"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/protocol/announcer"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/tecdsa/dkg"
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
	// dkgResultApprovalDelayStepBlocks determines the delay step in blocks
	// that is used to calculate the approval delay period that should be
	// respected by the given member to avoid all members approving the same
	// DKG result at the same time.
	dkgResultApprovalDelayStepBlocks = 5
)

// dkgExecutor is a component responsible for the full execution of ECDSA
// Distributed Key Generation: determining members selected to the signing
// group, executing off-chain protocol, and publishing the result to the chain.
type dkgExecutor struct {
	operatorID      chain.OperatorID
	operatorAddress chain.Address

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
	operatorID chain.OperatorID,
	operatorAddress chain.Address,
	chain Chain,
	netProvider net.Provider,
	walletRegistry *walletRegistry,
	protocolLatch *generator.ProtocolLatch,
	config Config,
	workPersistence persistence.BasicHandle,
	scheduler *generator.Scheduler,
	waitForBlockFn waitForBlockFn,
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
		operatorID:      operatorID,
		operatorAddress: operatorAddress,
		chain:           chain,
		netProvider:     netProvider,
		walletRegistry:  walletRegistry,
		protocolLatch:   protocolLatch,
		tecdsaExecutor:  tecdsaExecutor,
		waitForBlockFn:  waitForBlockFn,
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
	startBlock uint64,
) {
	dkgLogger := logger.With(
		zap.String("seed", fmt.Sprintf("0x%x", seed)),
	)

	dkgLogger.Info("checking eligibility for DKG")
	memberIndexes, groupSelectionResult, err := de.checkEligibility(
		dkgLogger,
	)
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
			groupSelectionResult,
			startBlock,
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
// - Group selection result holding chain.OperatorID and chain.Address for
//   operators selected to the signing group. There are always `groupSize`
//   selected operators.
func (de *dkgExecutor) checkEligibility(
	dkgLogger log.StandardLogger,
) ([]uint8, *GroupSelectionResult, error) {
	groupSelectionResult, err := de.chain.SelectGroup()
	if err != nil {
		return nil, nil, fmt.Errorf("selecting group not possible: [%v]", err)
	}

	dkgLogger.Infof(
		"selected group members for DKG = %s",
		groupSelectionResult.OperatorsAddresses,
	)

	if len(groupSelectionResult.OperatorsAddresses) > de.chain.GetConfig().GroupSize {
		return nil, nil, fmt.Errorf(
			"group size larger than supported: [%v]",
			len(groupSelectionResult.OperatorsAddresses),
		)
	}

	indexes := make([]uint8, 0)
	for index, operator := range groupSelectionResult.OperatorsAddresses {
		// See if we are amongst those chosen
		if operator == de.operatorAddress {
			// The group member index should be in range [1, groupSize] so we
			// need to add 1.
			indexes = append(indexes, uint8(index)+1)
		}
	}

	return indexes, groupSelectionResult, nil
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
	groupSelectionResult *GroupSelectionResult,
	startBlock uint64,
) {
	membershipValidator := group.NewMembershipValidator(
		dkgLogger,
		groupSelectionResult.OperatorsAddresses,
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
				startBlock,
				memberIndex,
				groupSelectionResult.OperatorsAddresses,
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

			signer, err := de.registerSigner(
				result,
				memberIndex,
				groupSelectionResult.OperatorsAddresses,
			)
			if err != nil {
				dkgLogger.Errorf(
					"[member:%v] failed to register signing group member: [%v]",
					memberIndex,
					err,
				)
			}

			dkgLogger.Infof("registered %s", signer)

			err = de.submitDkgResult(
				dkgLogger,
				seed,
				memberIndex,
				broadcastChannel,
				membershipValidator,
				result,
				groupSelectionResult,
				startBlock,
			)
			if err != nil {
				dkgLogger.Errorf(
					"[member:%v] DKG result publication process failed [%v]",
					memberIndex,
					err,
				)
			}
		}()
	}
}

// registerSigner determines the final signing group shape and persists the
// generated signer with a unique key share. Note that the final group members
// may differ from the ones returned by the sortition pool if there was any
// misbehavior or inactivities during the key generation.
func (de *dkgExecutor) registerSigner(
	result *dkg.Result,
	memberIndex group.MemberIndex,
	selectedSigningGroupOperators chain.Addresses,
) (*signer, error) {
	chainConfig := de.chain.GetConfig()
	// Final signing group may differ from the original DKG
	// group outputted by the sortition protocol. One need to
	// determine the final signing group based on the selected
	// group members who behaved correctly during DKG protocol.
	operatingMemberIndexes := result.Group.OperatingMemberIndexes()
	finalSigningGroupOperators, finalSigningGroupMembersIndexes, err :=
		finalSigningGroup(
			selectedSigningGroupOperators,
			operatingMemberIndexes,
			chainConfig,
		)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve final signing group members")
	}

	// Just like the final and original group may differ, the
	// member index used during the DKG protocol may differ
	// from the final signing group member index as well.
	// We need to remap it.
	finalSigningGroupMemberIndex, ok :=
		finalSigningGroupMembersIndexes[memberIndex]
	if !ok {
		return nil, fmt.Errorf("failed to resolve final signing " +
			"group member index",
		)
	}

	signer := newSigner(
		result.PrivateKeyShare.PublicKey(),
		finalSigningGroupOperators,
		finalSigningGroupMemberIndex,
		result.PrivateKeyShare,
	)

	err = de.walletRegistry.registerSigner(signer)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to register %s: [%v]",
			signer,
			err,
		)
	}

	return signer, nil
}

// submitDkgResult submits the DKG result to the chain.
func (de *dkgExecutor) submitDkgResult(
	dkgLogger log.StandardLogger,
	seed *big.Int,
	memberIndex group.MemberIndex,
	broadcastChannel net.BroadcastChannel,
	membershipValidator *group.MembershipValidator,
	dkgResult *dkg.Result,
	groupSelectionResult *GroupSelectionResult,
	startBlock uint64,
) error {
	// Set up the publication stop signal that should allow to
	// perform all the result-signing-related actions and
	// handle the worst case when the result is submitted by the
	// last group member.
	chainConfig := de.chain.GetConfig()
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

	return dkg.Publish(
		publicationCtx,
		dkgLogger,
		seed.Text(16),
		memberIndex,
		broadcastChannel,
		membershipValidator,
		newDkgResultSigner(de.chain, startBlock),
		newDkgResultSubmitter(dkgLogger, de.chain, groupSelectionResult),
		dkgResult,
	)
}

// executeDkgValidation performs the submitted DKG result validation process.
// If the result is not valid, this function submits an on-chain result
// challenge. If the result is valid and the given node was involved in the DKG,
// this function schedules an on-chain approve that is submitted once the
// challenge period elapses.
func (de *dkgExecutor) executeDkgValidation(
	seed *big.Int,
	submissionBlock uint64,
	result *DKGChainResult,
	resultHash [32]byte,
) {
	dkgLogger := logger.With(
		zap.String("seed", fmt.Sprintf("0x%x", seed)),
		zap.String("resultHash", fmt.Sprintf("0x%x", resultHash)),
	)

	isValid, err := de.chain.IsDKGResultValid(result)
	if err != nil {
		dkgLogger.Errorf("cannot validate DKG result: [%v]", err)
		return
	}

	if !isValid {
		dkgLogger.Infof("DKG result is invalid")

		err = de.chain.ChallengeDKGResult(result)
		if err != nil {
			dkgLogger.Errorf("cannot challenge invalid DKG result: [%v]", err)
			return
		}

		dkgLogger.Infof("challenged invalid DKG result")

		return
	}

	dkgLogger.Infof("DKG result is valid")

	// Determine the member indexes controlled by this node's operator.
	memberIndexes := make([]group.MemberIndex, 0)
	for index, operatorID := range result.Members {
		if operatorID == de.operatorID {
			// The group member index should be in range [1, groupSize] so we
			// need to add 1.
			memberIndexes = append(memberIndexes, group.MemberIndex(index+1))
		}
	}

	if len(memberIndexes) == 0 {
		dkgLogger.Infof("not eligible for DKG result approval")
		return
	}

	dkgLogger.Infof("scheduling DKG result approval")

	parameters, err := de.chain.DKGParameters()
	if err != nil {
		dkgLogger.Errorf("cannot get current DKG parameters: [%v]", err)
		return
	}

	// The challenge period starts at the result submission block and lasts
	// for challengePeriodBlocks.
	challengePeriodEndBlock := submissionBlock + parameters.ChallengePeriodBlocks
	// The approval is possible one block after the challenge period end.
	// The result submitter has precedence for approvePrecedencePeriodBlocks.
	approvePrecedencePeriodStartBlock := challengePeriodEndBlock + 1
	// Everyone else can approve once the precedence period ends.
	approvePeriodStartBlock := approvePrecedencePeriodStartBlock +
		parameters.ApprovePrecedencePeriodBlocks

	for _, currentMemberIndex := range memberIndexes {
		go func(memberIndex group.MemberIndex) {
			var approveBlock uint64

			if memberIndex == result.SubmitterMemberIndex {
				// The submitter can approve earlier, during the precedence
				// period.
				approveBlock = approvePrecedencePeriodStartBlock
			} else {
				// Everyone else must approve after the precedence period ends.
				// Each member preserves a delay according to their index
				// to avoid simultaneous approval.
				delayBlocks := uint64(memberIndex-1) * dkgResultApprovalDelayStepBlocks
				approveBlock = approvePeriodStartBlock + delayBlocks
			}

			dkgLogger.Infof(
				"[member:%v] waiting for block [%v] to approve DKG result",
				memberIndex,
				approveBlock,
			)

			ctx, cancelCtx := context.WithCancel(context.Background())
			defer cancelCtx()

			subscription := de.chain.OnDKGResultApproved(
				func(event *DKGResultApprovedEvent) {
					cancelCtx()
				},
			)
			defer subscription.Unsubscribe()

			err := de.waitForBlockFn(ctx, approveBlock)
			if err != nil {
				dkgLogger.Errorf(
					"[member:%v] error while waiting for DKG result "+
						"approve block: [%v]",
					memberIndex,
					err,
				)
				return
			}

			// If the context got cancelled that means the result was approved
			// by someone else.
			if ctx.Err() != nil {
				dkgLogger.Infof(
					"[member:%v] DKG result approved by someone else",
					memberIndex,
				)
				return
			}

			err = de.chain.ApproveDKGResult(result)
			if err != nil {
				dkgLogger.Errorf(
					"[member:%v] cannot approve DKG result: [%v]",
					memberIndex,
					err,
				)
				return
			}

			dkgLogger.Errorf("[member:%v] approved DKG result", memberIndex)
		}(currentMemberIndex)
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
