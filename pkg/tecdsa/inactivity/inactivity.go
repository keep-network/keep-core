package inactivity

import (
	"context"
	"fmt"

	"github.com/ipfs/go-log/v2"

	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/protocol/state"
)

// SignedClaim represents information pertaining to the process of signing
// an inactivity claim: the public key used during signing, the resulting
// signature and the hash of the inactivity claim that was used during signing.
type SignedClaim struct {
	PublicKey []byte
	Signature []byte
	ClaimHash ClaimSignatureHash
}

type ClaimSigner interface {
	SignClaim(claim *Claim) (*SignedClaim, error)
	VerifySignature(signedClaim *SignedClaim) (bool, error)
}

type ClaimSubmitter interface {
	SubmitClaim(
		ctx context.Context,
		memberIndex group.MemberIndex,
		claim *Claim,
		signatures map[group.MemberIndex][]byte,
	) error
}

func Publish(
	ctx context.Context,
	logger log.StandardLogger,
	sessionID string,
	memberIndex group.MemberIndex,
	channel net.BroadcastChannel,
	groupSize int,
	dishonestThreshold int,
	membershipValidator *group.MembershipValidator,
	claimSigner ClaimSigner,
	claimSubmitter ClaimSubmitter,
	claim *Claim,
) error {
	initialState := &claimSigningState{
		BaseAsyncState: state.NewBaseAsyncState(),
		channel:        channel,
		claimSigner:    claimSigner,
		claimSubmitter: claimSubmitter,
		member: newSigningMember(
			logger,
			memberIndex,
			groupSize,
			dishonestThreshold,
			membershipValidator,
			sessionID,
		),
		claim: claim,
	}

	stateMachine := state.NewAsyncMachine(logger, ctx, channel, initialState)

	lastState, err := stateMachine.Execute()
	if err != nil {
		return err
	}

	_, ok := lastState.(*claimSubmissionState)
	if !ok {
		return fmt.Errorf("execution ended on state %T", lastState)
	}

	return nil
}
