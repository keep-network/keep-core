package tbtc

import (
	"context"
	"fmt"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/protocol/announcer"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/tecdsa"
	"github.com/keep-network/keep-core/pkg/tecdsa/signing"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"
	"math/big"
	"reflect"
)

// TODO: Documentation.
type signingGroupController struct {
	signers             []*signer
	broadcastChannel    net.BroadcastChannel
	membershipValidator *group.MembershipValidator
	chainConfig         *ChainConfig

	// waitForBlockFn is a function used to wait for the given block.
	waitForBlockFn waitForBlockFn
	// onSignerStartFn is a callback function invoked when a single signer
	// starts the execution of the signing protocol.
	onSignerStartFn func()
	// onSignerEndFn is a callback function invoked when a single signer
	// end the execution of the signing protocol, regardless of the outcome.
	onSignerEndFn func()

	// signingAttemptsLimit determines the maximum attempts count that will
	// be made by a single signer for the given message. Once the attempts
	// limit is hit the signer gives up.
	signingAttemptsLimit uint
}

// TODO: Documentation.
func (sgc *signingGroupController) signBatch(
	messages []*big.Int,
	startBlockNumber uint64,
) ([]*tecdsa.Signature, error) {
	// TODO: Implementation.
	return nil, nil
}

// TODO: Documentation.
func (sgc *signingGroupController) sign(
	ctx context.Context,
	message *big.Int,
	startBlockNumber uint64,
) (*tecdsa.Signature, uint64, error) {
	// All signers belong to one wallet. Take that wallet from the
	// first signer.
	wallet := sgc.signers[0].wallet
	// Actual wallet signing group size may be different from the
	// `GroupSize` parameter of the chain config.
	signingGroupSize := len(wallet.signingGroupOperators)
	// The dishonest threshold for the wallet signing group must be
	// also calculated using the actual wallet signing group size.
	signingGroupDishonestThreshold := signingGroupSize -
		sgc.chainConfig.HonestThreshold

	walletPublicKeyBytes, err := marshalPublicKey(wallet.publicKey)
	if err != nil {
		return nil, 0, fmt.Errorf("cannot marshal wallet public key: [%v]", err)
	}

	signingLogger := logger.With(
		zap.String("wallet", fmt.Sprintf("0x%x", walletPublicKeyBytes)),
		zap.String("message", fmt.Sprintf("0x%x", message)),
	)

	type signerOutcome struct {
		signature *tecdsa.Signature
		endBlock  uint64
		err       error
	}

	signerOutcomesChan := make(chan *signerOutcome, len(sgc.signers))

	for _, currentSigner := range sgc.signers {
		go func(signer *signer) {
			sgc.onSignerStartFn()
			defer sgc.onSignerEndFn()

			announcer := announcer.New(
				fmt.Sprintf("%v-%v", ProtocolName, "signing"),
				sgc.chainConfig.GroupSize,
				sgc.broadcastChannel,
				sgc.membershipValidator,
			)

			syncer := newSigningSyncer(
				sgc.chainConfig.GroupSize,
				sgc.broadcastChannel,
				sgc.membershipValidator,
			)

			retryLoop := newSigningRetryLoop(
				signingLogger,
				message,
				startBlockNumber,
				signer.signingGroupMemberIndex,
				wallet.signingGroupOperators,
				sgc.chainConfig,
				announcer,
				syncer,
			)

			// Set up the loop timeout signal.
			loopTimeoutBlock := startBlockNumber +
				uint64(sgc.signingAttemptsLimit*signingAttemptMaximumBlocks())
			// Do not cancel the loopCtx upon function exit and continue
			// to broadcast sync messages until the attempt timeout. This way
			// we maximize the chance that other members, especially the
			// ones not participating in the successful signature get synced
			// as well.
			loopCtx, _ := withCancelOnBlock(
				ctx,
				loopTimeoutBlock,
				sgc.waitForBlockFn,
			)

			result, endBlock, err := retryLoop.start(
				loopCtx,
				sgc.waitForBlockFn,
				func(attempt *signingAttemptParams) (*signing.Result, uint64, error) {
					signingAttemptLogger := signingLogger.With(
						zap.Uint("attempt", attempt.number),
						zap.Uint64("attemptStartBlock", attempt.startBlock),
						zap.Uint64("attemptTimeoutBlock", attempt.timeoutBlock),
					)

					signingAttemptLogger.Infof(
						"[member:%v] starting signing protocol "+
							"with [%v] group members (excluded: [%v])",
						signer.signingGroupMemberIndex,
						signingGroupSize-len(attempt.excludedMembersIndexes),
						attempt.excludedMembersIndexes,
					)

					// Set up the attempt timeout signal.
					attemptCtx, _ := withCancelOnBlock(
						loopCtx,
						attempt.timeoutBlock,
						sgc.waitForBlockFn,
					)

					sessionID := fmt.Sprintf(
						"%v-%v",
						message.Text(16),
						attempt.number,
					)

					result, err := signing.Execute(
						attemptCtx,
						signingAttemptLogger,
						message,
						sessionID,
						signer.signingGroupMemberIndex,
						signer.privateKeyShare,
						signingGroupSize,
						signingGroupDishonestThreshold,
						attempt.excludedMembersIndexes,
						sgc.broadcastChannel,
						sgc.membershipValidator,
					)
					if err != nil {
						return nil, 0, err
					}

					return result, 0, nil
				},
			)
			if err != nil {
				signingLogger.Errorf(
					"[member:%v] all retries for the signing failed; "+
						"giving up: [%v]",
					signer.signingGroupMemberIndex,
					err,
				)
				signerOutcomesChan <- &signerOutcome{
					signature: nil,
					endBlock:  0,
					err:       err,
				}
				return
			}

			signingLogger.Infof(
				"[member:%v] generated signature [%v]",
				signer.signingGroupMemberIndex,
				result.Signature,
			)

			signerOutcomesChan <- &signerOutcome{
				signature: result.Signature,
				endBlock:  endBlock,
				err:       nil,
			}
		}(currentSigner)
	}

	signerOutcomes := make([]*signerOutcome, 0)

	for {
		select {
		case outcome := <-signerOutcomesChan:
			signerOutcomes = append(signerOutcomes, outcome)

			if len(signerOutcomes) == len(sgc.signers) {
				outcomes := slices.CompactFunc(
					signerOutcomes,
					func(o1, o2 *signerOutcome) bool {
						return o1.signature.Equals(o2.signature) &&
							o1.endBlock == o2.endBlock &&
							reflect.DeepEqual(o1.err, o2.err)
					},
				)

				if len(outcomes) != 1 {
					return nil, 0, fmt.Errorf("signers came to different outcomes")
				}

				commonOutcome := outcomes[0]

				if err := commonOutcome.err; err != nil {
					return nil, 0, fmt.Errorf("signers failed: [%v]", err)
				}

				return commonOutcome.signature, commonOutcome.endBlock, nil
			}
		case <-ctx.Done():
			return nil, 0, ctx.Err()
		}
	}
}
