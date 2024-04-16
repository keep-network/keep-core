package inactivity

import (
	"context"
	"fmt"

	"github.com/ipfs/go-log/v2"

	"github.com/keep-network/keep-core/pkg/protocol/group"
)

type signingMember struct {
	logger log.StandardLogger
	// Index of this group member.
	memberIndex group.MemberIndex
	// Group to which this member belongs.
	group *group.Group
	// Validator allowing to check public key and member index against
	// group members.
	membershipValidator *group.MembershipValidator
	// Identifier of the particular operator inactivity notification session
	// this member is part of.
	sessionID string
	// Hash of inactivity claim preferred by the current participant.
	preferredInactivityClaimHash ClaimSignatureHash
	// Signature over preferredInactivityClaimHash calculated by the member.
	selfInactivityClaimSignature []byte
}

// newSigningMember creates a new signingMember in the initial state.
func newSigningMember(
	logger log.StandardLogger,
	memberIndex group.MemberIndex,
	groupSize int,
	dishonestThreshold int,
	membershipValidator *group.MembershipValidator,
	sessionID string,
) *signingMember {
	return &signingMember{
		logger:      logger,
		memberIndex: memberIndex,
		// TODO: Check is this is a correct way to create the group.
		group:               group.NewGroup(dishonestThreshold, groupSize),
		membershipValidator: membershipValidator,
		sessionID:           sessionID,
	}
}

// shouldAcceptMessage indicates whether the given member should accept
// a message from the given sender.
func (sm *signingMember) shouldAcceptMessage(
	senderID group.MemberIndex,
	senderPublicKey []byte,
) bool {
	isMessageFromSelf := senderID == sm.memberIndex
	isSenderValid := sm.membershipValidator.IsValidMembership(
		senderID,
		senderPublicKey,
	)
	isSenderAccepted := sm.group.IsOperating(senderID)

	return !isMessageFromSelf && isSenderValid && isSenderAccepted
}

// initializeSubmittingMember performs a transition of a member state to the
// next phase of the protocol.
func (sm *signingMember) initializeSubmittingMember() *submittingMember {
	return &submittingMember{
		signingMember: sm,
	}
}

func (sm *signingMember) signClaim(
	claim *Claim,
	claimSigner ClaimSigner,
) (*claimSignatureMessage, error) {
	signedClaim, err := claimSigner.SignClaim(claim)
	if err != nil {
		return nil, fmt.Errorf("failed to sign inactivity claim [%v]", err)
	}

	// Register self signature and claim hash.
	sm.selfInactivityClaimSignature = signedClaim.Signature
	sm.preferredInactivityClaimHash = signedClaim.ClaimHash

	return &claimSignatureMessage{
		senderID:  sm.memberIndex,
		claimHash: signedClaim.ClaimHash,
		signature: signedClaim.Signature,
		publicKey: signedClaim.PublicKey,
		sessionID: sm.sessionID,
	}, nil
}

// verifyInactivityClaimSignatures verifies signatures received in messages from
// other group members. It collects signatures supporting only the same
// inactivity claim hash as the one preferred by the current member. Each member
// is allowed to broadcast only one signature over a preferred inactivity claim
// hash. The function assumes that the input messages list does not contain a
// message from self and that the public key presented in each message is the
// correct one. This key needs to be compared against the one used by network
// client earlier, before this function is called.
func (sm *signingMember) verifyInactivityClaimSignatures(
	messages []*claimSignatureMessage,
	resultSigner ClaimSigner,
) map[group.MemberIndex][]byte {
	// TODO: Implement
	return nil
}

// submittingMember represents a member submitting an inactivity claim to the
// blockchain along with signatures received from other group members supporting
// the claim.
type submittingMember struct {
	*signingMember
}

// submitClaim submits the inactivity claim along with the supporting signatures
// to the provided claim submitter.
func (sm *submittingMember) submitClaim(
	ctx context.Context,
	claim *Claim,
	signatures map[group.MemberIndex][]byte,
	claimSubmitter ClaimSubmitter,
) error {
	if err := claimSubmitter.SubmitClaim(
		ctx,
		sm.memberIndex,
		claim,
		signatures,
	); err != nil {
		return fmt.Errorf("failed to submit inactivity [%v]", err)
	}

	return nil
}
