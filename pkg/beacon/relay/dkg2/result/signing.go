package result

import (
	"fmt"

	"github.com/keep-network/keep-core/pkg/operator"

	relayChain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/beacon/relay/member"
	"github.com/keep-network/keep-core/pkg/chain"
)

// SigningMember represents a member sharing preferred DKG result hash
// and signature over this hash with peer members.
type SigningMember struct {
	index member.Index

	chainHandle chain.Handle

	// Key used for signing the DKG result hash.
	privateKey *operator.PrivateKey

	// Hash of DKG result preferred by the current participant.
	preferredDKGResultHash relayChain.DKGResultHash
	// Received valid signatures supporting the same DKG result as current's
	// participant prefers. Contains also current's participant's signature.
	receivedValidResultSignatures map[member.Index]operator.Signature
}

// SignDKGResult calculates hash of DKG result and member's signature over this
// hash. It packs the hash and signature into a broadcast message.
//
// See Phase 13 of the protocol specification.
func (fm *SigningMember) SignDKGResult(dkgResult *relayChain.DKGResult) (
	*DKGResultHashSignatureMessage,
	error,
) {
	resultHash, err := fm.chainHandle.ThresholdRelay().CalculateDKGResultHash(dkgResult)
	if err != nil {
		return nil, fmt.Errorf("dkg result hash calculation failed [%v]", err)
	}
	fm.preferredDKGResultHash = resultHash

	signature, err := operator.Sign(resultHash[:], fm.privateKey)
	if err != nil {
		return nil, fmt.Errorf("dkg result hash signing failed [%v]", err)
	}

	// Register self signature.
	fm.receivedValidResultSignatures[fm.index] = signature

	return &DKGResultHashSignatureMessage{
		senderIndex: fm.index,
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
func (fm *SigningMember) VerifyDKGResultSignatures(
	messages []*DKGResultHashSignatureMessage,
) error {
	duplicatedMessagesFromSender := func(senderIndex member.Index) bool {
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

	for _, message := range messages {
		// Check if message from self.
		if message.senderIndex == fm.index {
			continue
		}

		// Check if sender sent multiple messages.
		if duplicatedMessagesFromSender(message.senderIndex) {
			fmt.Printf(
				"[member: %v] received multiple messages from sender [%d]",
				fm.index,
				message.senderIndex,
			)
			continue
		}

		// Sender's preferred DKG result hash doesn't match current member's
		// preferred DKG result hash.
		if message.resultHash != fm.preferredDKGResultHash {
			fmt.Printf(
				"[member: %v] signature from sender [%d] supports result different than preferred",
				fm.index,
				message.senderIndex,
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
				message.senderIndex,
				message,
			)
			continue
		}

		fm.receivedValidResultSignatures[message.senderIndex] = message.signature
	}

	return nil
}
