package tbtc

import (
	"bytes"
	"fmt"
	"math/big"
	"sort"

	"github.com/keep-network/keep-common/pkg/persistence"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/internal/testutils"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/tecdsa/dkg"
)

// TODO: Unit tests for `node.go`.

// node represents the current state of an ECDSA node.
type node struct {
	chain          Chain
	netProvider    net.Provider
	walletRegistry *walletRegistry
	dkgExecutor    *dkg.Executor
}

func newNode(
	chain Chain,
	netProvider net.Provider,
	persistence persistence.Handle,
	config Config,
) *node {
	walletRegistry := newWalletRegistry(persistence)

	dkgExecutor := dkg.NewExecutor(
		logger,
		config.PreParamsPoolSize,
		config.PreParamsGenerationTimeout,
		config.PreParamsGenerationConcurrency,
	)

	return &node{
		chain:          chain,
		netProvider:    netProvider,
		walletRegistry: walletRegistry,
		dkgExecutor:    dkgExecutor,
	}
}

// joinDKGIfEligible takes a seed value and undergoes the process of the
// distributed key generation if this node's operator proves to be eligible for
// the group generated by that seed. This is an interactive on-chain process,
// and joinDKGIfEligible can block for an extended period of time while it
// completes the on-chain operation.
func (n *node) joinDKGIfEligible(seed *big.Int, startBlockNumber uint64) {
	logger.Infof(
		"checking eligibility for DKG with seed [0x%x]",
		seed,
	)

	selectedSigningGroupOperators, err := n.chain.SelectGroup(seed)
	if err != nil {
		logger.Errorf(
			"failed to select group with seed [0x%x]: [%v]",
			seed,
			err,
		)
		return
	}

	chainConfig := n.chain.GetConfig()

	if len(selectedSigningGroupOperators) > chainConfig.GroupSize {
		logger.Errorf(
			"group size larger than supported: [%v]",
			len(selectedSigningGroupOperators),
		)
		return
	}

	signing := n.chain.Signing()

	_, operatorPublicKey, err := n.chain.OperatorKeyPair()
	if err != nil {
		logger.Errorf("failed to get operator public key: [%v]", err)
		return
	}

	operatorAddress, err := signing.PublicKeyToAddress(operatorPublicKey)
	if err != nil {
		logger.Errorf("failed to get operator address: [%v]", err)
		return
	}

	indexes := make([]uint8, 0)
	for index, operator := range selectedSigningGroupOperators {
		// See if we are amongst those chosen
		if operator == operatorAddress {
			indexes = append(indexes, uint8(index))
		}
	}

	// Create temporary broadcast channel name for DKG using the
	// group selection seed with the protocol name as prefix.
	channelName := fmt.Sprintf("%s-%s", ProtocolName, seed.Text(16))

	if len(indexes) > 0 {
		logger.Infof(
			"joining DKG with seed [0x%x] and controlling [%v] group members",
			seed,
			len(indexes),
		)

		broadcastChannel, err := n.netProvider.BroadcastChannelFor(channelName)
		if err != nil {
			logger.Errorf("failed to get broadcast channel: [%v]", err)
			return
		}

		membershipValidator := group.NewMembershipValidator(
			&testutils.MockLogger{},
			selectedSigningGroupOperators,
			signing,
		)

		err = broadcastChannel.SetFilter(membershipValidator.IsInGroup)
		if err != nil {
			logger.Errorf(
				"could not set filter for channel [%v]: [%v]",
				broadcastChannel.Name(),
				err,
			)
		}

		blockCounter, err := n.chain.BlockCounter()
		if err != nil {
			logger.Errorf("failed to get block counter: [%v]", err)
			return
		}

		for _, index := range indexes {
			// Capture the member index for the goroutine. The group member
			// index should be in range [1, groupSize] so we need to add 1.
			memberIndex := index + 1

			go func() {
				result, endBlock, err := n.dkgExecutor.Execute(
					seed,
					startBlockNumber,
					memberIndex,
					chainConfig.GroupSize,
					chainConfig.DishonestThreshold(),
					blockCounter,
					broadcastChannel,
					membershipValidator,
				)
				if err != nil {
					// TODO: Add retries into the mix.
					logger.Errorf(
						"[member:%v] failed to execute dkg: [%v]",
						memberIndex,
						err,
					)
					return
				}

				publicationStartBlock := endBlock
				operatingMemberIDs := result.Group.OperatingMemberIDs()
				dkgResultChannel := make(chan *dkg.ResultSubmissionEvent)

				dkgResultSubscription := n.chain.OnDKGResultSubmitted(
					func(event *dkg.ResultSubmissionEvent) {
						dkgResultChannel <- event
					},
				)
				defer dkgResultSubscription.Unsubscribe()

				err = dkg.Publish(
					logger,
					publicationStartBlock,
					memberIndex,
					blockCounter,
					broadcastChannel,
					membershipValidator,
					n.chain,
					n.chain,
					result,
				)
				if err != nil {
					// Result publication failed. It means that either the result this
					// member proposed is not supported by the majority of group members or
					// that the chain interaction failed. In either case, we observe the
					// chain for the result published by any other group member and based
					// on that, we decide whether we should stay in the final group
					// or drop our membership.
					logger.Warningf(
						"[member:%v] DKG result publication process failed [%v]",
						memberIndex,
						err,
					)

					if operatingMemberIDs, err = n.decideMemberFate(
						memberIndex,
						result,
						dkgResultChannel,
						publicationStartBlock,
					); err != nil {
						logger.Errorf(
							"failed to handle DKG result publishing failure: [%v]",
							err,
						)
						return
					}
				}

				signingGroupOperators, err := n.resolveSigningGroupOperators(
					selectedSigningGroupOperators,
					operatingMemberIDs,
				)
				if err != nil {
					logger.Errorf(
						"failed to resolve group operators: [%v]",
						err,
					)
					return
				}

				// TODO: Snapshot the key material before doing on-chain result
				//       submission.

				// TODO: The final `signingGroupOperators` may differ from
				//       the original `selectedSigningGroupOperators`.
				//       Consider that when integrating the retry algorithm.
				signer := newSigner(
					result.PrivateKeyShare.PublicKey(),
					signingGroupOperators,
					memberIndex,
					result.PrivateKeyShare,
				)

				err = n.walletRegistry.registerSigner(signer)
				if err != nil {
					logger.Errorf(
						"failed to register %s: [%v]",
						signer,
						err,
					)
					return
				}

				logger.Infof("registered %s", signer)
			}()
		}
	} else {
		logger.Infof("not eligible for DKG with seed [0x%x]", seed)
	}
}

// decideMemberFate decides what the member will do in case it failed to publish
// its DKG result. Member can stay in the group if it supports the same group
// public key as the one registered on-chain and the member is not considered as
// misbehaving by the group.
func (n *node) decideMemberFate(
	memberIndex group.MemberIndex,
	result *dkg.Result,
	dkgResultChannel chan *dkg.ResultSubmissionEvent,
	publicationStartBlock uint64,
) ([]group.MemberIndex, error) {
	dkgResultEvent, err := n.waitForDkgResultEvent(
		dkgResultChannel,
		publicationStartBlock,
	)
	if err != nil {
		return nil, err
	}

	groupPublicKeyBytes, err := result.GetGroupPublicKeyBytes()
	if err != nil {
		return nil, err
	}

	// If member doesn't support the same group public key, it could not stay
	// in the group.
	if !bytes.Equal(groupPublicKeyBytes, dkgResultEvent.GroupPublicKeyBytes) {
		return nil, fmt.Errorf(
			"[member:%v] could not stay in the group because "+
				"the member do not support the same group public key",
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
	operatingMemberIDs := make([]group.MemberIndex, 0)
	for _, memberID := range result.Group.MemberIDs() {
		if _, isMisbehaved := misbehavedSet[memberID]; !isMisbehaved {
			operatingMemberIDs = append(operatingMemberIDs, memberID)
		}
	}

	return operatingMemberIDs, nil
}

// waitForDkgResultEvent waits for the DKG result submission event. It times out
// and returns error if the DKG result event is not emitted on time.
func (n *node) waitForDkgResultEvent(
	dkgResultChannel chan *dkg.ResultSubmissionEvent,
	publicationStartBlock uint64,
) (*dkg.ResultSubmissionEvent, error) {
	config := n.chain.GetConfig()

	timeoutBlock := publicationStartBlock + dkg.PrePublicationBlocks() +
		(uint64(config.GroupSize) * config.ResultPublicationBlockStep)

	blockCounter, err := n.chain.BlockCounter()
	if err != nil {
		return nil, err
	}

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

// resolveSigningGroupOperators takes two parameters:
// - selectedOperators: Contains addresses of all selected operators. Slice
//   length equals to the groupSize. Each element with index N corresponds
//   to the group member with ID N+1.
// - operatingGroupMembersIDs: Contains group members IDs that were neither
//   disqualified nor marked as inactive. Slice length is lesser than or equal
//   to the groupSize.
//
// Using those parameters, this function transforms the selectedOperators
// slice into another slice that contains addresses of all operators
// that were neither disqualified nor marked as inactive. This way, the
// resulting slice has only addresses of properly operating operators
// who form the resulting group.
//
// Example:
// selectedOperators: [member1, member2, member3, member4, member5]
// operatingGroupMembersIDs: [5, 1, 3]
// signingGroupOperators: [member1, member3, member5]
func (n *node) resolveSigningGroupOperators(
	selectedOperators []chain.Address,
	operatingGroupMembersIDs []group.MemberIndex,
) ([]chain.Address, error) {
	config := n.chain.GetConfig()

	if len(selectedOperators) != config.GroupSize ||
		len(operatingGroupMembersIDs) < config.HonestThreshold {
		return nil, fmt.Errorf("invalid input parameters")
	}

	sort.Slice(operatingGroupMembersIDs, func(i, j int) bool {
		return operatingGroupMembersIDs[i] < operatingGroupMembersIDs[j]
	})

	signingGroupOperators := make(
		[]chain.Address,
		len(operatingGroupMembersIDs),
	)

	for i, operatingMemberID := range operatingGroupMembersIDs {
		signingGroupOperators[i] = selectedOperators[operatingMemberID-1]
	}

	return signingGroupOperators, nil
}
