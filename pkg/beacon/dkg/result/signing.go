package result

import (
	"fmt"
	"github.com/ipfs/go-log/v2"

	"github.com/keep-network/keep-core/pkg/chain"

	beaconchain "github.com/keep-network/keep-core/pkg/beacon/chain"
	"github.com/keep-network/keep-core/pkg/protocol/group"
)

// SigningMember represents a group member sharing their preferred DKG result hash
// and signature (over this hash) with other peer members.
type SigningMember struct {
	logger log.StandardLogger

	index group.MemberIndex

	// Group to which this member belongs.
	group *group.Group

	// Validator allowing to check public key and member index
	// against group members
	membershipValidator *group.MembershipValidator

	// Hash of DKG result preferred by the current participant.
	preferredDKGResultHash beaconchain.DKGResultHash
	// Signature over preferredDKGResultHash calculated by the member.
	selfDKGResultSignature []byte

	// Identifier of the particular DKG session this member is part of.
	sessionID string
}

// NewSigningMember creates a member to execute signing DKG result hash.
func NewSigningMember(
	logger log.StandardLogger,
	memberIndex group.MemberIndex,
	dkgGroup *group.Group,
	membershipValidator *group.MembershipValidator,
	sessionID string,
) *SigningMember {
	return &SigningMember{
		logger:              logger,
		index:               memberIndex,
		group:               dkgGroup,
		membershipValidator: membershipValidator,
		sessionID:           sessionID,
	}
}

// SignDKGResult calculates hash of DKG result and member's signature over this
// hash. It packs the hash and signature into a broadcast message.
//
// See Phase 13 of the protocol specification.
func (sm *SigningMember) SignDKGResult(
	dkgResult *beaconchain.DKGResult,
	beaconChain beaconchain.Interface,
) (
	*DKGResultHashSignatureMessage,
	error,
) {
	resultHash, err := beaconChain.CalculateDKGResultHash(dkgResult)
	if err != nil {
		return nil, fmt.Errorf("dkg result hash calculation failed [%v]", err)
	}
	sm.preferredDKGResultHash = resultHash

	signing := beaconChain.Signing()

	signature, err := signing.Sign(resultHash[:])
	if err != nil {
		return nil, fmt.Errorf("dkg result hash signing failed [%v]", err)
	}

	// Register self signature.
	sm.selfDKGResultSignature = signature

	return &DKGResultHashSignatureMessage{
		senderIndex: sm.index,
		resultHash:  resultHash,
		signature:   signature,
		publicKey:   signing.PublicKey(),
		sessionID:   sm.sessionID,
	}, nil
}

// VerifyDKGResultSignatures verifies signatures received in messages from other
// group members.
//
// It collects signatures supporting only the same DKG result hash as the one
// preferred by the current member.
//
// Each member is allowed to broadcast only one signature over a preferred DKG
// result hash.
//
// The function assumes that the public key presented in the message is the
// correct one. This key needs to be compared against the one used by network
// client earlier, before this function is called.
//
// See Phase 13 of the protocol specification.
func (sm *SigningMember) VerifyDKGResultSignatures(
	messages []*DKGResultHashSignatureMessage,
	signing chain.Signing,
) (map[group.MemberIndex][]byte, error) {
	duplicatedMessagesFromSender := func(senderIndex group.MemberIndex) bool {
		messageFromSenderAlreadySeen := false
		for _, message := range messages {
			if message.senderIndex == senderIndex {
				if messageFromSenderAlreadySeen {
					return true
				}
				messageFromSenderAlreadySeen = true
			}
		}
		return false
	}

	receivedValidResultSignatures := make(map[group.MemberIndex][]byte)

	for _, message := range messages {
		// Check if message from self.
		if message.senderIndex == sm.index {
			continue
		}

		// Check if sender sent multiple messages.
		if duplicatedMessagesFromSender(message.senderIndex) {
			sm.logger.Infof(
				"[member:%v] received multiple messages from sender: [%d]",
				sm.index,
				message.senderIndex,
			)
			continue
		}

		// Sender's preferred DKG result hash doesn't match current member's
		// preferred DKG result hash.
		if message.resultHash != sm.preferredDKGResultHash {
			sm.logger.Infof(
				"[member:%v] signature from sender [%d] supports result different than preferred",
				sm.index,
				message.senderIndex,
			)
			continue
		}

		// Check if the signature is valid.
		ok, err := signing.VerifyWithPublicKey(
			message.resultHash[:],
			message.signature,
			message.publicKey,
		)
		if err != nil {
			sm.logger.Infof(
				"[member:%v] verification of signature from sender [%d] failed: [%v]",
				sm.index,
				message.senderIndex,
				err,
			)
			continue
		}
		if !ok {
			sm.logger.Infof(
				"[member:%v] sender [%d] provided invalid signature",
				sm.index,
				message.senderIndex,
			)
			continue
		}

		receivedValidResultSignatures[message.senderIndex] = message.signature
	}

	// Register member's self signature.
	receivedValidResultSignatures[sm.index] = sm.selfDKGResultSignature

	return receivedValidResultSignatures, nil
}

// shouldAcceptMessage indicates whether the given member should accept
// a message from the given sender.
func (sm *SigningMember) shouldAcceptMessage(
	senderID group.MemberIndex,
	senderPublicKey []byte,
) bool {
	isMessageFromSelf := senderID == sm.index
	isSenderValid := sm.membershipValidator.IsValidMembership(
		senderID,
		senderPublicKey,
	)
	isSenderAccepted := sm.group.IsOperating(senderID)

	return !isMessageFromSelf && isSenderValid && isSenderAccepted
}
