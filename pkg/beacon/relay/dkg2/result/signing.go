package result

import (
	"fmt"

	"github.com/keep-network/keep-core/pkg/beacon/relay/gjkr"
	"github.com/keep-network/keep-core/pkg/chain"
	"github.com/keep-network/keep-core/pkg/operator"

	relayChain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
)

// SigningMember represents a member sharing preferred DKG result hash
// and signature over this hash with peer members.
type SigningMember struct {
	index gjkr.MemberID

	// Key used for signing the DKG result hash.
	privateKey *operator.PrivateKey

	// Hash of DKG result preferred by the current participant.
	preferredDKGResultHash relayChain.DKGResultHash

	selfDKGResultSignature operator.Signature
}

type SignaturesVerifyingMember struct {
	*SigningMember
}

func NewSigningMember(memberindex gjkr.MemberID, operatorPrivateKey *operator.PrivateKey) *SigningMember {
	return &SigningMember{
		index:      memberindex,
		privateKey: operatorPrivateKey,
	}
}

// SignDKGResult calculates hash of DKG result and member's signature over this
// hash. It packs the hash and signature into a broadcast message.
//
// See Phase 13 of the protocol specification.
func (fm *SigningMember) SignDKGResult(dkgResult *relayChain.DKGResult, chainHandle chain.Handle) (
	*DKGResultHashSignatureMessage,
	error,
) {
	resultHash, err := chainHandle.ThresholdRelay().CalculateDKGResultHash(dkgResult)
	if err != nil {
		return nil, fmt.Errorf("dkg result hash calculation failed [%v]", err)
	}
	fm.preferredDKGResultHash = resultHash

	signature, err := operator.Sign(resultHash[:], fm.privateKey)
	if err != nil {
		return nil, fmt.Errorf("dkg result hash signing failed [%v]", err)
	}

	fm.selfDKGResultSignature = signature

	return &DKGResultHashSignatureMessage{
		senderindex: fm.index,
		resultHash:  resultHash,
		signature:   signature,
		publicKey:   &fm.privateKey.PublicKey,
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
// See Phase 13 of the protocol specification.
func (fm *SignaturesVerifyingMember) VerifyDKGResultSignatures(
	messages []*DKGResultHashSignatureMessage,
) (map[gjkr.MemberID]operator.Signature, error) {
	// Received valid signatures supporting the same DKG result as current's
	// participant prefers. Contains also current's participant's signature.
	receivedValidResultSignatures := make(map[gjkr.MemberID]operator.Signature)

	duplicatedMessagesFromSender := func(senderindex gjkr.MemberID) bool {
		messageFromSenderAlreadySeen := false
		for _, message := range messages {
			if message.senderindex == senderindex {
				if messageFromSenderAlreadySeen {
					return true
				}
				messageFromSenderAlreadySeen = true
			}
		}
		return false
	}

	for _, message := range messages {
		// Check if message from self.
		if message.senderindex == fm.index {
			continue
		}

		// Check if sender sent multiple messages.
		if duplicatedMessagesFromSender(message.senderindex) {
			fmt.Printf(
				"[member: %v] received multiple messages from sender [%d]",
				fm.index,
				message.senderindex,
			)
			continue
		}

		// Sender's preferred DKG result hash doesn't match current member's
		// preferred DKG result hash.
		if message.resultHash != fm.preferredDKGResultHash {
			fmt.Printf(
				"[member: %v] signature from sender [%d] supports result different than preferred",
				fm.index,
				message.senderindex,
			)
			continue
		}

		// Signature is invalid.
		err := operator.VerifySignature(
			message.publicKey,
			message.resultHash[:],
			message.signature,
		)
		if err != nil {
			fmt.Printf(
				"[member: %v] verification of signature from sender [%d] failed [%+v]",
				fm.index,
				message.senderindex,
				message,
			)
			continue
		}

		receivedValidResultSignatures[message.senderindex] = message.signature
	}

	// Register member's self signature.
	receivedValidResultSignatures[fm.index] = fm.selfDKGResultSignature

	return receivedValidResultSignatures, nil
}

func (sm *SigningMember) InitializeSignaturesVerification() *SignaturesVerifyingMember {
	return &SignaturesVerifyingMember{sm}
}
