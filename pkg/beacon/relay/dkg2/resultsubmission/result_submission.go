package resultsubmission

import (
	"fmt"
	"os"

	relayChain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
)

// Signature ...
type Signature []byte

// SignDKGResult  // TODO: Write docs
func (fm *ResultSigningMember) SignDKGResult(dkgResult *relayChain.DKGResult) (
	*DKGResultHashSignatureMessage,
	error,
) {
	resultHash, err := fm.chainHandle.ThresholdRelay().CalculateDKGResultHash(dkgResult)
	if err != nil {
		return nil, fmt.Errorf("dkg result hash calculation failed [%v]", err)
	}
	fm.dkgResultHash = resultHash

	signature, err := fm.sign(resultHash)
	if err != nil {
		return nil, fmt.Errorf("dkg result hash signing failed [%v]", err)
	}

	// Register self signature.
	fm.validResultSignatures[fm.index] = signature

	return &DKGResultHashSignatureMessage{
		senderIndex: fm.index,
		resultHash:  resultHash,
		signature:   signature,
	}, nil
}

// VerifyDKGResultSignatures // TODO: Write docs
func (fm *ResultSigningMember) VerifyDKGResultSignatures(
	messages []*DKGResultHashSignatureMessage,
) (map[ParticipantIndex][]Signature, error) {
	alreadyReceivedSignature := make([]ParticipantIndex, 0) // track if member already send signature
	accusations := make(map[ParticipantIndex][]Signature)

messagesCheck:
	for _, message := range messages {
		// Check if sender sent multiple signatures.
		for _, alreadySignedIndex := range alreadyReceivedSignature {
			if message.senderIndex == alreadySignedIndex {
				fmt.Println("message from member who already send a message")

				if signature, ok := fm.validResultSignatures[message.senderIndex]; ok {
					accusations[message.senderIndex] = append(
						accusations[message.senderIndex],
						signature,
					)

					delete(fm.validResultSignatures, message.senderIndex)
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
		if message.resultHash != fm.dkgResultHash {
			fmt.Println("signature for result different than preferred")
			continue
		}

		// Signature is invalid.
		if !fm.verifySignature(message.senderIndex, message.resultHash, message.signature) {
			fmt.Fprintf(os.Stderr, "invalid signature in message: [%+v]", message)
			// TODO: Should we accuse the member who send invalid signature?
			continue
		}

		fm.validResultSignatures[message.senderIndex] = message.signature
	}

	return accusations, nil
}

func (fm *ResultSigningMember) sign(resultHash relayChain.DKGResultHash) []byte {
	// TODO: Implement
	return append([]byte("Signed:"), resultHash[:]...)
}

func (fm *ResultSigningMember) verifySignature(
	participantIndex ParticipantIndex,
	hash relayChain.DKGResultHash,
	signature Signature,
) bool {
	// TODO: Implement
	// ecdsa.Verify(fm.publicKeys[participantIndex], hash, r, s)
	return true
}
