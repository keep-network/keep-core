package tbtc

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"go.uber.org/zap"

	"github.com/keep-network/keep-common/pkg/persistence"
	"github.com/keep-network/keep-core/pkg/generator"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/tecdsa/dkg"
	"github.com/keep-network/keep-core/pkg/tecdsa/signing"
)

// TODO: Unit tests for `node.go`.

// node represents the current state of an ECDSA node.
type node struct {
	chain          Chain
	netProvider    net.Provider
	walletRegistry *walletRegistry
	dkgExecutor    *dkg.Executor
	protocolLatch  *generator.ProtocolLatch
}

func newNode(
	chain Chain,
	netProvider net.Provider,
	keyStorePersistance persistence.ProtectedHandle,
	workPersistence persistence.BasicHandle,
	scheduler *generator.Scheduler,
	config Config,
) *node {
	walletRegistry := newWalletRegistry(keyStorePersistance)

	dkgExecutor := dkg.NewExecutor(
		logger,
		scheduler,
		workPersistence,
		config.PreParamsPoolSize,
		config.PreParamsGenerationTimeout,
		config.PreParamsGenerationDelay,
		config.PreParamsGenerationConcurrency,
		config.KeyGenerationConcurrency,
	)

	latch := generator.NewProtocolLatch()
	scheduler.RegisterProtocol(latch)

	return &node{
		chain:          chain,
		netProvider:    netProvider,
		walletRegistry: walletRegistry,
		dkgExecutor:    dkgExecutor,
		protocolLatch:  latch,
	}
}

// joinDKGIfEligible takes a seed value and undergoes the process of the
// distributed key generation if this node's operator proves to be eligible for
// the group generated by that seed. This is an interactive on-chain process,
// and joinDKGIfEligible can block for an extended period of time while it
// completes the on-chain operation.
func (n *node) joinDKGIfEligible(seed *big.Int, startBlockNumber uint64) {
	dkgLogger := logger.With(
		zap.String("seed", fmt.Sprintf("0x%x", seed)),
	)

	dkgLogger.Info("checking eligibility for DKG")

	selectedSigningGroupOperators, err := n.chain.SelectGroup(seed)
	if err != nil {
		// TODO: We should consider switching this log to Errorf when the
		// Chaosnet 0 phase is completed and results are submitted to the chain.
		// To let the operators join the pool, the dev team will keep unlocking
		// it via DKG timeout from time to time. During the period pool is
		// unlocked, selecting the group will not work.
		//dkgLogger.Errorf("failed to select group: [%v]", err)
		dkgLogger.Warnf("selecting group not possible: [%v]", err)
		return
	}

	dkgLogger.Infof("selected group members for DKG = %s", selectedSigningGroupOperators)

	chainConfig := n.chain.GetConfig()

	if len(selectedSigningGroupOperators) > chainConfig.GroupSize {
		dkgLogger.Errorf(
			"group size larger than supported: [%v]",
			len(selectedSigningGroupOperators),
		)
		return
	}

	signing := n.chain.Signing()

	_, operatorPublicKey, err := n.chain.OperatorKeyPair()
	if err != nil {
		dkgLogger.Errorf("failed to get operator public key: [%v]", err)
		return
	}

	operatorAddress, err := signing.PublicKeyToAddress(operatorPublicKey)
	if err != nil {
		dkgLogger.Errorf("failed to get operator address: [%v]", err)
		return
	}

	indexes := make([]uint8, 0)
	for index, operator := range selectedSigningGroupOperators {
		// See if we are amongst those chosen
		if operator == operatorAddress {
			indexes = append(indexes, uint8(index))
		}
	}

	if membersCount := len(indexes); membersCount > 0 {
		if preParamsCount := n.dkgExecutor.PreParamsCount(); membersCount > preParamsCount {
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

		// Create temporary broadcast channel name for DKG using the
		// group selection seed with the protocol name as prefix.
		channelName := fmt.Sprintf("%s-%s", ProtocolName, seed.Text(16))

		broadcastChannel, err := n.netProvider.BroadcastChannelFor(channelName)
		if err != nil {
			dkgLogger.Errorf("failed to get broadcast channel: [%v]", err)
			return
		}

		dkg.RegisterUnmarshallers(broadcastChannel)
		registerStopPillUnmarshaller(broadcastChannel)

		membershipValidator := group.NewMembershipValidator(
			dkgLogger,
			selectedSigningGroupOperators,
			signing,
		)

		err = broadcastChannel.SetFilter(membershipValidator.IsInGroup)
		if err != nil {
			dkgLogger.Errorf(
				"could not set filter for channel [%v]: [%v]",
				broadcastChannel.Name(),
				err,
			)
		}

		blockCounter, err := n.chain.BlockCounter()
		if err != nil {
			dkgLogger.Errorf("failed to get block counter: [%v]", err)
			return
		}

		for _, index := range indexes {
			// Capture the member index for the goroutine. The group member
			// index should be in range [1, groupSize] so we need to add 1.
			memberIndex := index + 1

			go func() {
				n.protocolLatch.Lock()
				defer n.protocolLatch.Unlock()

				retryLoop := newDkgRetryLoop(
					seed,
					startBlockNumber,
					memberIndex,
					selectedSigningGroupOperators,
					chainConfig,
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
				cancelContextOnStopSignal(loopCtx, cancelLoopCtx, broadcastChannel)

				result, executionEndBlock, err := retryLoop.start(
					loopCtx,
					n.waitForBlockHeight,
					func(attempt *dkgAttemptParams) (*dkg.Result, uint64, error) {
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

						// sessionID must be different for each attempt.
						sessionID := fmt.Sprintf(
							"%v-%v",
							seed.Text(16),
							attempt.number,
						)

						result, executionEndBlock, err := n.dkgExecutor.Execute(
							dkgAttemptLogger,
							seed,
							sessionID,
							attempt.startBlock,
							memberIndex,
							chainConfig.GroupSize,
							chainConfig.DishonestThreshold(),
							attempt.excludedMembersIndexes,
							blockCounter,
							broadcastChannel,
							membershipValidator,
						)
						if err != nil {
							dkgAttemptLogger.Errorf(
								"[member:%v] dkg attempt failed: [%v]",
								memberIndex,
								err,
							)

							return nil, 0, err
						}

						if err := sendStopPill(loopCtx, broadcastChannel, attempt.number); err != nil {
							dkgLogger.Errorf(
								"[member:%v] could not send the stop pill: [%v]",
								memberIndex,
								err,
							)
						}

						return result, executionEndBlock, nil
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

				publicationStartBlock := executionEndBlock
				operatingMemberIndexes := result.Group.OperatingMemberIDs()
				dkgResultChannel := make(chan *DKGResultSubmittedEvent)

				dkgResultSubscription := n.chain.OnDKGResultSubmitted(
					func(event *DKGResultSubmittedEvent) {
						dkgResultChannel <- event
					},
				)
				defer dkgResultSubscription.Unsubscribe()

				err = dkg.Publish(
					dkgLogger,
					seed.Text(16),
					publicationStartBlock,
					memberIndex,
					blockCounter,
					broadcastChannel,
					membershipValidator,
					newDkgResultSigner(n.chain),
					newDkgResultSubmitter(dkgLogger, n.chain),
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
						memberIndex,
						dkgResultChannel,
						publicationStartBlock,
						result,
						chainConfig,
						blockCounter,
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

				err = n.walletRegistry.registerSigner(signer)
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
	} else {
		dkgLogger.Infof("not eligible for DKG")
	}
}

// joinSigningIfEligible takes a message and undergoes the process of tECDSA
// signing if this node's operator proves to control some signers of the
// requested wallet. This is an interactive process, and joinSigningIfEligible
// can block for an extended period of time while it completes the operation.
//
// TODO: Ultimately, one client will handle multiple signature requests for
//
//	multiple wallets. Because of that, logging within that function
//	must use the context of the given signature request.
func (n *node) joinSigningIfEligible(
	message *big.Int,
	walletPublicKey *ecdsa.PublicKey,
	startBlockNumber uint64,
) {
	signingLogger := logger.With(
		zap.String("message", fmt.Sprintf("0x%x", message)),
	)

	walletPublicKeyBytes, err := marshalPublicKey(walletPublicKey)
	if err != nil {
		signingLogger.Errorf("cannot marshal wallet public key: [%v]", err)
		return
	}

	signingLogger = signingLogger.With(
		zap.String("wallet", fmt.Sprintf("0x%x", walletPublicKeyBytes)),
	)

	signingLogger.Info("checking eligibility for signing")

	if signers := n.walletRegistry.getSigners(
		walletPublicKey,
	); len(signers) > 0 {
		signingLogger.Infof(
			"joining signing; controlling [%v] signers",
			len(signers),
		)

		// All signers belong to one wallet. Take that wallet from the
		// first signer.
		wallet := signers[0].wallet
		// Actual wallet signing group size may be different from the
		// `GroupSize` parameter of the chain config.
		signingGroupSize := len(wallet.signingGroupOperators)
		// The dishonest threshold for the wallet signing group must be
		// also calculated using the actual wallet signing group size.
		signingGroupDishonestThreshold := signingGroupSize -
			n.chain.GetConfig().HonestThreshold

		channelName := fmt.Sprintf(
			"%s-%s",
			ProtocolName,
			hex.EncodeToString(walletPublicKeyBytes),
		)

		broadcastChannel, err := n.netProvider.BroadcastChannelFor(channelName)
		if err != nil {
			signingLogger.Errorf("failed to get broadcast channel: [%v]", err)
			return
		}

		signing.RegisterUnmarshallers(broadcastChannel)
		registerStopPillUnmarshaller(broadcastChannel)

		membershipValidator := group.NewMembershipValidator(
			signingLogger,
			wallet.signingGroupOperators,
			n.chain.Signing(),
		)

		err = broadcastChannel.SetFilter(membershipValidator.IsInGroup)
		if err != nil {
			signingLogger.Errorf(
				"could not set filter for channel [%v]: [%v]",
				broadcastChannel.Name(),
				err,
			)
		}

		blockCounter, err := n.chain.BlockCounter()
		if err != nil {
			signingLogger.Errorf("failed to get block counter: [%v]", err)
			return
		}

		for _, currentSigner := range signers {
			go func(signer *signer) {
				n.protocolLatch.Lock()
				defer n.protocolLatch.Unlock()

				retryLoop := newSigningRetryLoop(
					message,
					startBlockNumber,
					signer.signingGroupMemberIndex,
					wallet.signingGroupOperators,
					n.chain.GetConfig(),
				)

				// TODO: For this client iteration, the signing loop is started
				//       with a 24h timeout and a stop pill sent by any group
				//       member. Once the RandomBeacon is integrated, the stop
				//       signal should be generated by observing the DKG result
				//       submission or timeout.
				loopCtx, cancelLoopCtx := context.WithTimeout(
					context.Background(),
					24*time.Hour,
				)
				defer cancelLoopCtx()
				cancelContextOnStopSignal(loopCtx, cancelLoopCtx, broadcastChannel)

				result, err := retryLoop.start(
					loopCtx,
					n.waitForBlockHeight,
					func(attempt *signingAttemptParams) (*signing.Result, error) {
						signingAttemptLogger := signingLogger.With(
							zap.Uint("attempt", attempt.number),
							zap.Uint64("attemptStartBlock", attempt.startBlock),
						)

						signingAttemptLogger.Infof(
							"[member:%v] scheduled signing attempt "+
								"with [%v] group members (excluded: [%v])",
							signer.signingGroupMemberIndex,
							signingGroupSize-len(attempt.excludedMembersIndexes),
							attempt.excludedMembersIndexes,
						)

						sessionID := fmt.Sprintf(
							"%v-%v",
							message.Text(16),
							attempt.number,
						)

						result, err := signing.Execute(
							signingAttemptLogger,
							message,
							sessionID,
							attempt.startBlock,
							signer.signingGroupMemberIndex,
							signer.privateKeyShare,
							signingGroupSize,
							signingGroupDishonestThreshold,
							attempt.excludedMembersIndexes,
							blockCounter,
							broadcastChannel,
							membershipValidator,
						)
						if err != nil {
							signingAttemptLogger.Errorf(
								"[member:%v] signing attempt failed: [%v]",
								signer.signingGroupMemberIndex,
								err,
							)

							return nil, err
						}

						if err := sendStopPill(loopCtx, broadcastChannel, attempt.number); err != nil {
							signingLogger.Errorf(
								"[member:%v] could not send the stop pill: [%v]",
								signer.signingGroupMemberIndex,
								err,
							)
						}

						return result, nil
					},
				)
				if err != nil {
					signingLogger.Errorf(
						"[member:%v] all retries for the signing failed; "+
							"giving up: [%v]",
						signer.signingGroupMemberIndex,
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
					signingLogger.Infof(
						"[member:%v] signing retry loop received stop signal",
						signer.signingGroupMemberIndex,
					)
					return
				}

				signingLogger.Infof(
					"[member:%v] generated signature [%v]",
					signer.signingGroupMemberIndex,
					result.Signature,
				)
			}(currentSigner)
		}
	} else {
		signingLogger.Info("not eligible for signing")
	}
}

// TODO: this should become a part of BlockHeightWaiter interface.
func (n *node) waitForBlockHeight(ctx context.Context, blockHeight uint64) error {
	blockCounter, err := n.chain.BlockCounter()
	if err != nil {
		return err
	}

	wait, err := blockCounter.BlockHeightWaiter(blockHeight)
	if err != nil {
		return err
	}

	select {
	case <-wait:
	case <-ctx.Done():
	}

	return nil
}
