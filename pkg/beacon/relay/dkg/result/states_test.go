package result

import (
	"crypto/ecdsa"
	"reflect"
	"testing"

	relayChain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/operator"
)

func TestAcceptValidSignatureHashMessage(t *testing.T) {
	groupSize := 2

	dkgResult := &relayChain.DKGResult{
		GroupPublicKey: []byte("He’s the hero Gotham deserves."),
	}

	members, chainHandles, publicKeys, err := initializeSigningMembers(groupSize)
	if err != nil {
		t.Fatal(err)
	}

	member, _ := members[0], chainHandles[0]
	member2, chain2, publicKey2 := members[1], chainHandles[1], publicKeys[1]

	message2, err := member2.SignDKGResult(
		dkgResult,
		chain2.ThresholdRelay(),
		chain2.Signing(),
	)

	state := &resultSigningState{
		member:            member,
		signatureMessages: make([]*DKGResultHashSignatureMessage, 0),
	}

	state.Receive(&mockSignatureMessage{
		message2,
		*publicKey2,
		chain2.Signing().PublicKey(),
	})

	if len(state.signatureMessages) != 1 {
		t.Fatalf("Expected one signature hash message accepted")
	}
	if !reflect.DeepEqual(state.signatureMessages[0], message2) {
		t.Fatalf(
			"Unexpected accepted message\nExpected: %v\nActual:   %v\n",
			message2,
			state.signatureMessages[0],
		)
	}
}

func TestDoNotAcceptMessageWithSwappedKey(t *testing.T) {
	groupSize := 2

	dkgResult := &relayChain.DKGResult{
		GroupPublicKey: []byte("But not the one it needs right now."),
	}

	members, chainHandles, _, err := initializeSigningMembers(groupSize)
	if err != nil {
		t.Fatal(err)
	}

	member, _ := members[0], chainHandles[0]
	member2, chain2 := members[1], chainHandles[1]

	message2, err := member2.SignDKGResult(
		dkgResult,
		chain2.ThresholdRelay(),
		chain2.Signing(),
	)

	state := &resultSigningState{
		member:            member,
		signatureMessages: make([]*DKGResultHashSignatureMessage, 0),
	}

	_, publicKey, err := operator.GenerateKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	state.Receive(&mockSignatureMessage{
		message2,
		*publicKey,
		[]byte("operator uses another key"),
	})

	if len(state.signatureMessages) != 0 {
		t.Fatalf("Expected no signature hash message accepted")
	}
}

type mockSignatureMessage struct {
	payload              *DKGResultHashSignatureMessage
	senderPublicKey      ecdsa.PublicKey
	senderPublicKeyBytes []byte
}

func (msm *mockSignatureMessage) TransportSenderID() net.TransportIdentifier {
	panic("not implemented")
}
func (msm *mockSignatureMessage) Payload() interface{} {
	return msm.payload
}
func (msm *mockSignatureMessage) Type() string {
	panic("not implemented")
}
func (msm *mockSignatureMessage) SenderPublicKey() ecdsa.PublicKey {
	return msm.senderPublicKey
}
func (msm *mockSignatureMessage) SenderPublicKeyBytes() []byte {
	return msm.senderPublicKeyBytes
}
