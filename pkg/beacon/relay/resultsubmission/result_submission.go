package resultsubmission

import (
	"fmt"
	"os"

	relayChain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
)

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

	signature := fm.sign(resultHash)
	// Register self signature.
	fm.validResultSignatures[fm.index] = signature

	return &DKGResultHashSignatureMessage{
		senderIndex: fm.index,
		resultHash:  resultHash,
		signature:   signature,
	}, nil
}

// VerifyDKGResultSignatures // TODO: Write docs
func (fm *ResultSigningMember) VerifyDKGResultSignatures(messages []*DKGResultHashSignatureMessage) error {
	for _, message := range messages {
		if !fm.verifySignature(message.senderIndex, message.resultHash, message.signature) {
			fmt.Fprintf(os.Stderr, "invalid signature in message: [%+v]", message)
			continue
		}

		if message.resultHash != fm.dkgResultHash {
			fmt.Printf("signature for result different than preferred")
			continue
		}

		fm.validResultSignatures[message.senderIndex] = message.signature
	}

	return nil
}

func (fm *ResultSigningMember) sign(resultHash relayChain.DKGResultHash) []byte {
	// TODO: Implement
	return append([]byte("Signed:"), resultHash[:]...)
}

func (fm *ResultSigningMember) verifySignature(
	participantIndex ParticipantIndex,
	hash relayChain.DKGResultHash,
	signature []byte,
) bool {
	// TODO: Implement
	// ecdsa.Verify(fm.publicKeys[participantIndex], hash, r, s)
	return true
}
