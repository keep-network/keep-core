package result

import (
	"crypto/ecdsa"
	"fmt"
	"os"

	"github.com/keep-network/keep-core/pkg/beacon/relay/gjkr"

	relayChain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/chain"
)

// SigningMember represents a member sharing preferred DKG result hash
// and signature over this hash with peer members.
type SigningMember struct {
	index gjkr.MemberID

	chainHandle chain.Handle

	// Key used for signing the DKG result hash.
	privateKey *ecdsa.PrivateKey // TODO: Change to static.PrivateKey

	// Hash of DKG result preferred by the current participant.
	preferredDKGResultHash relayChain.DKGResultHash
	// Received valid signatures supporting the same DKG result as current's
	// participants prefer. Contains also current's participant's signature.
	receivedValidResultSignatures map[gjkr.MemberID]Signature // TODO: Change to static.Signature
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

	signature, err := fm.sign(resultHash) // TODO: Update to fm.privateKey.Sign(resultHash)
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
// result hash. This function tracks members who delivered multiple signatures.
// It returns the map of members' indices along with signatures of members who
// delivered multiple signatures.
//
// See Phase 13 of the protocol specification.
func (fm *SigningMember) VerifyDKGResultSignatures(
	messages []*DKGResultHashSignatureMessage,
) (map[gjkr.MemberID][]Signature, error) {
	// alreadyReceivedSignature tracks if the other member already send a signature.
	alreadyReceivedSignature := make([]gjkr.MemberID, 0)
	// accusations collects indices and signatures of members who delivered
	// multiple signatures.
	accusations := make(map[gjkr.MemberID][]Signature)

messagesCheck:
	for _, message := range messages {
		// Check if message from self.
		if message.senderIndex == fm.index {
			continue
		}

		// Check if sender sent multiple signatures.
		for _, alreadySignedIndex := range alreadyReceivedSignature {
			if message.senderIndex == alreadySignedIndex {
				fmt.Println("message from member who already send a message")

				if signature, ok := fm.receivedValidResultSignatures[message.senderIndex]; ok {
					accusations[message.senderIndex] = append(
						accusations[message.senderIndex],
						signature,
					)

					delete(fm.receivedValidResultSignatures, message.senderIndex)
				}

				accusations[message.senderIndex] = append(
					accusations[message.senderIndex],
					message.signature,
				)

				continue messagesCheck
			}
		}
		alreadyReceivedSignature = append(alreadyReceivedSignature, message.senderIndex)

		// Sender's preferred DKG result hash doesn't match current member's
		// preferred DKG result hash.
		if message.resultHash != fm.preferredDKGResultHash {
			fmt.Println("signature for result different than preferred")
			continue
		}

		// Signature is invalid.
		if !fm.verifySignature( // TODO: Change to static.VerifySignature
			message.senderIndex, // TODO: Chagne to fm.otherMembersPublicKeys[message.senderIndex]
			message.resultHash,
			message.signature,
			message.publicKey,
		) {
			fmt.Fprintf(os.Stderr, "invalid signature in message: [%+v]", message)
			// TODO: Should we accuse the member who send invalid signature?
			continue
		}

		fm.receivedValidResultSignatures[message.senderIndex] = message.signature
	}

	return accusations, nil
}
