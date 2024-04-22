package tbtc

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"github.com/ipfs/go-log/v2"
	"go.uber.org/zap"
	"golang.org/x/sync/semaphore"

	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/generator"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/tecdsa/inactivity"
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
	protocolLatch 		*generator.ProtocolLatch

	waitForBlockFn waitForBlockFn
}

// TODO Consider moving all inactivity-related code to pkg/protocol/inactivity.
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

	execLogger := logger.With(
		zap.String("wallet", fmt.Sprintf("0x%x", walletPublicKeyBytes)),
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

	wg := sync.WaitGroup{}
	wg.Add(len(ice.signers))

	for _, currentSigner := range ice.signers {
		ice.protocolLatch.Lock()
		defer ice.protocolLatch.Unlock()

		defer wg.Done()

		inactivityClaimTimeoutBlock := uint64(0) // TODO: Set the value of timeout block

		go func(signer *signer) {
			ctx, cancelCtx := withCancelOnBlock(
				context.Background(),
				inactivityClaimTimeoutBlock,
				ice.waitForBlockFn,
			)
			defer cancelCtx()

			ice.publish(
				ctx,
				execLogger,
				message,
				signer.signingGroupMemberIndex,
				wallet.groupSize(),
				wallet.groupDishonestThreshold(
					ice.groupParameters.HonestThreshold,
				),
				ice.membershipValidator,
				claim,
			)

		}(currentSigner)
	}

	// Wait until all controlled signers complete their routine.
	wg.Wait()

	return nil
}

func (ice *inactivityClaimExecutor) publish(
	ctx context.Context,
	inactivityLogger log.StandardLogger,
	seed *big.Int,
	memberIndex group.MemberIndex,
	groupSize int,
	dishonestThreshold int,
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
		newInactivityClaimSubmitter(),
		inactivityClaim,
	)
}

func (ice *inactivityClaimExecutor) wallet() wallet {
	// All signers belong to one wallet. Take that wallet from the
	// first signer.
	return ice.signers[0].wallet
}
