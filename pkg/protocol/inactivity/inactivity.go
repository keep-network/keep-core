package inactivity

import (
	"context"
	"fmt"

	"github.com/ipfs/go-log/v2"

	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/protocol/group"
	"github.com/keep-network/keep-core/pkg/protocol/state"
)

// SignedClaimHash represents information pertaining to the process of signing
// an inactivity claim: the public key used during signing, the resulting
// signature and the hash of the inactivity claim that was used during signing.
type SignedClaimHash struct {
	PublicKey []byte
	Signature []byte
	ClaimHash ClaimHash
}

type ClaimSigner interface {
	SignClaim(claim *ClaimPreimage) (*SignedClaimHash, error)
	VerifySignature(signedClaim *SignedClaimHash) (bool, error)
}

type ClaimSubmitter interface {
	SubmitClaim(
		ctx context.Context,
		memberIndex group.MemberIndex,
		claim *ClaimPreimage,
		signatures map[group.MemberIndex][]byte,
	) error
}

func PublishClaim(
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
	claim *ClaimPreimage,
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

// RegisterUnmarshallers initializes the given broadcast channel to be able to
// perform inactivity claim interactions by registering all the required
// protocol message unmarshallers.
func RegisterUnmarshallers(channel net.BroadcastChannel) {
	channel.SetUnmarshaler(func() net.TaggedUnmarshaler {
		return &claimSignatureMessage{}
	})
}
