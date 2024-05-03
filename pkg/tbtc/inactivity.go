package tbtc

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"

	"github.com/ipfs/go-log/v2"
	"go.uber.org/zap"
	"golang.org/x/sync/semaphore"

	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/generator"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/protocol/inactivity"
)

const (
	// inactivityClaimSubmissionDelayStepBlocks determines the delay step in blocks
	// that is used to calculate the submission delay period that should be respected
	// by the given member to avoid all members submitting the same inactivity claim
	// at the same time.
	inactivityClaimSubmissionDelayStepBlocks = 3
	// inactivityClaimMaximumSubmissionBlocks determines the maximum block
	// duration of inactivity claim submission procedure.
	inactivityClaimMaximumSubmissionBlocks = 60
)

// errInactivityClaimExecutorBusy is an error returned when the inactivity claim
// executor cannot execute the inactivity claim due to another inactivity claim
// execution in progress.
var errInactivityClaimExecutorBusy = fmt.Errorf("inactivity claim executor is busy")

type inactivityClaimExecutor struct {
	lock *semaphore.Weighted

	chain               Chain
	signers             []*signer
	broadcastChannel    net.BroadcastChannel
	membershipValidator *group.MembershipValidator
	groupParameters     *GroupParameters
	protocolLatch       *generator.ProtocolLatch

	waitForBlockFn waitForBlockFn
}

func newInactivityClaimExecutor(
	chain Chain,
	signers []*signer,
	broadcastChannel net.BroadcastChannel,
	membershipValidator *group.MembershipValidator,
	groupParameters *GroupParameters,
	protocolLatch *generator.ProtocolLatch,
	waitForBlockFn waitForBlockFn,
) *inactivityClaimExecutor {
	return &inactivityClaimExecutor{
		lock:                semaphore.NewWeighted(1),
		chain:               chain,
		signers:             signers,
		broadcastChannel:    broadcastChannel,
		membershipValidator: membershipValidator,
		groupParameters:     groupParameters,
		protocolLatch:       protocolLatch,
		waitForBlockFn:      waitForBlockFn,
	}
}

func (ice *inactivityClaimExecutor) publishClaim(
	inactiveMembersIndexes []group.MemberIndex,
	heartbeatFailed bool,
	message *big.Int,
	startBlock uint64,
) error {
	if lockAcquired := ice.lock.TryAcquire(1); !lockAcquired {
		return errInactivityClaimExecutorBusy
	}
	defer ice.lock.Release(1)

	wallet := ice.wallet()

	walletPublicKeyHash := bitcoin.PublicKeyHash(wallet.publicKey)
	walletPublicKeyBytes, err := marshalPublicKey(wallet.publicKey)
	if err != nil {
		return fmt.Errorf("cannot marshal wallet public key: [%v]", err)
	}

	timeoutBlock := startBlock + inactivityClaimMaximumSubmissionBlocks

	execLogger := logger.With(
		zap.String("wallet", fmt.Sprintf("0x%x", walletPublicKeyBytes)),
		zap.Uint64("inactivityClaimStartBlock", startBlock),
		zap.Uint64("inactivityClaimTimeoutBlock", timeoutBlock),
	)

	walletRegistryData, err := ice.chain.GetWallet(walletPublicKeyHash)
	if err != nil {
		return fmt.Errorf("could not get registry data on wallet: [%v]", err)
	}

	nonce, err := ice.chain.GetInactivityClaimNonce(
		walletRegistryData.EcdsaWalletID,
	)
	if err != nil {
		return fmt.Errorf("could not get nonce for wallet: [%v]", err)
	}

	claim := &inactivity.Claim{
		Nonce:                  nonce,
		WalletPublicKey:        wallet.publicKey,
		InactiveMembersIndexes: inactiveMembersIndexes,
		HeartbeatFailed:        heartbeatFailed,
	}

	groupMembers, err := ice.getWalletMembersInfo()
	if err != nil {
		return fmt.Errorf("could not get wallet members info: [%v]", err)
	}

	wg := sync.WaitGroup{}
	wg.Add(len(ice.signers))

	for _, currentSigner := range ice.signers {
		ice.protocolLatch.Lock()
		defer ice.protocolLatch.Unlock()

		defer wg.Done()

		go func(signer *signer) {
			execLogger.Info(
				"[member:%v] starting inactivity claim publishing",
				signer.signingGroupMemberIndex,
			)

			ctx, cancelCtx := withCancelOnBlock(
				context.Background(),
				timeoutBlock,
				ice.waitForBlockFn,
			)
			defer cancelCtx()

			subscription := ice.chain.OnInactivityClaimed(
				func(event *InactivityClaimedEvent) {
					defer cancelCtx()

					execLogger.Infof(
						"[member:%v] Inactivity claim submitted for wallet "+
							"with ID [0x%x] and nonce [%v] by notifier [%v] "+
							"at block [%v]",
						signer.signingGroupMemberIndex,
						event.WalletID,
						event.Nonce,
						event.Notifier,
						event.BlockNumber,
					)
				})
			defer subscription.Unsubscribe()

			err := ice.publish(
				ctx,
				execLogger,
				message,
				signer.signingGroupMemberIndex,
				wallet.groupSize(),
				wallet.groupDishonestThreshold(
					ice.groupParameters.HonestThreshold,
				),
				groupMembers,
				ice.membershipValidator,
				claim,
			)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					execLogger.Infof(
						"[member:%v] inactivity claim is no longer awaiting "+
							"publishing; aborting inactivity claim publishing",
						signer.signingGroupMemberIndex,
					)
					return
				}

				execLogger.Errorf(
					"[member:%v] inactivity claim publishing failed [%v]",
					signer.signingGroupMemberIndex,
					err,
				)
				return
			}
		}(currentSigner)
	}

	// Wait until all controlled signers complete their routine.
	wg.Wait()

	return nil
}

func (ice *inactivityClaimExecutor) getWalletMembersInfo() ([]uint32, error) {
	// Cache mapping operator addresses to their wallet member IDs. It helps to
	// limit the number of calls to the ETH client if some operator addresses
	// occur on the list multiple times.
	operatorIDCache := make(map[chain.Address]uint32)

	walletMemberIDs := make([]uint32, 0)

	for _, operatorAddress := range ice.wallet().signingGroupOperators {
		// Search for the operator address in the cache. Store the operator
		// address in the cache if it's not there.
		if operatorID, found := operatorIDCache[operatorAddress]; !found {
			fetchedOperatorID, err := ice.chain.GetOperatorID(operatorAddress)
			if err != nil {
				return nil, fmt.Errorf("could not get operator ID: [%w]", err)
			}
			operatorIDCache[operatorAddress] = fetchedOperatorID
			walletMemberIDs = append(walletMemberIDs, fetchedOperatorID)
		} else {
			walletMemberIDs = append(walletMemberIDs, operatorID)
		}
	}

	return walletMemberIDs, nil
}

func (ice *inactivityClaimExecutor) publish(
	ctx context.Context,
	inactivityLogger log.StandardLogger,
	seed *big.Int,
	memberIndex group.MemberIndex,
	groupSize int,
	dishonestThreshold int,
	groupMembers []uint32,
	membershipValidator *group.MembershipValidator,
	inactivityClaim *inactivity.Claim,
) error {
	return inactivity.Publish(
		ctx,
		inactivityLogger,
		seed.Text(16),
		memberIndex,
		ice.broadcastChannel,
		groupSize,
		dishonestThreshold,
		membershipValidator,
		newInactivityClaimSigner(ice.chain),
		newInactivityClaimSubmitter(
			inactivityLogger,
			ice.chain,
			ice.groupParameters,
			groupMembers,
			ice.waitForBlockFn,
		),
		inactivityClaim,
	)
}

func (ice *inactivityClaimExecutor) wallet() wallet {
	// All signers belong to one wallet. Take that wallet from the
	// first signer.
	return ice.signers[0].wallet
}
