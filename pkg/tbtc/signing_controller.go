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
	"math/big"
	"time"
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
			loopCtx, cancelLoopCtx := withCancelOnBlock(
				ctx,
				loopTimeoutBlock,
				sgc.waitForBlockFn,
			)

			cancelSigningContextOnStopSignal(
				loopCtx,
				cancelLoopCtx,
				sgc.broadcastChannel,
				message.Text(16),
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

					// Schedule the stop pill to be sent a fixed amount of
					// time after the result is returned. Do not do it
					// immediately as other members can be very close
					// to produce the result as well. This mechanism should
					// be more sophisticated but since it is temporary, we
					// can live with it for now.
					go func() {
						time.Sleep(1 * time.Minute)
						if err := sendSigningStopPill(
							loopCtx,
							sgc.broadcastChannel,
							message.Text(16),
							attempt.number,
						); err != nil {
							signingLogger.Errorf(
								"[member:%v] could not send the stop pill: [%v]",
								signer.signingGroupMemberIndex,
								err,
							)
						}
					}()

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

				signerOutcomesChan <- &signerOutcome{err: err}

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

				signerOutcomesChan <- &signerOutcome{}

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

	// TODO: Try to simplify the following code.

	signerOutcomes := make([]*signerOutcome, 0)

outcomesLoop:
	for {
		select {
		case outcome := <-signerOutcomesChan:
			signerOutcomes = append(signerOutcomes, outcome)

			if len(signerOutcomes) == len(sgc.signers) {
				break outcomesLoop
			}
		case <-ctx.Done():
			return nil, 0, ctx.Err()
		}
	}

	var signature *tecdsa.Signature
	var endBlock uint64

	for _, outcome := range signerOutcomes {
		if outcome.err != nil {
			return nil, 0, fmt.Errorf("failed signers")
		}

		if signature == nil {
			signature = outcome.signature
		}

		if endBlock == 0 {
			endBlock = outcome.endBlock
		}

		if signature.String() != outcome.signature.String() {
			return nil, 0, fmt.Errorf("signers came with different signatures")
		}

		if endBlock != outcome.endBlock {
			return nil, 0, fmt.Errorf("signers came with different end blocks")
		}
	}

	return signature, endBlock, nil
}
