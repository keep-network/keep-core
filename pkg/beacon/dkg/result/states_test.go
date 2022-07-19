package result

import (
	beaconchain "github.com/keep-network/keep-core/pkg/beacon/chain"
	"reflect"
	"testing"

	"github.com/keep-network/keep-core/pkg/net"
)

func TestAcceptValidSignatureHashMessage(t *testing.T) {
	groupSize := 2

	dkgResult := &beaconchain.DKGResult{
		GroupPublicKey: []byte("He’s the hero Gotham deserves."),
	}

	members, beaconChains, err := initializeSigningMembers(groupSize)
	if err != nil {
		t.Fatal(err)
	}

	member, _ := members[0], beaconChains[0]
	member2, beaconChain2 := members[1], beaconChains[1]

	message2, err := member2.SignDKGResult(
		dkgResult,
		beaconChain2,
	)

	state := &resultSigningState{
		member:            member,
		signatureMessages: make([]*DKGResultHashSignatureMessage, 0),
	}

	state.Receive(&mockSignatureMessage{
		message2,
		beaconChain2.Signing().PublicKey(),
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

	dkgResult := &beaconchain.DKGResult{
		GroupPublicKey: []byte("But not the one it needs right now."),
	}

	members, beaconChains, err := initializeSigningMembers(groupSize)
	if err != nil {
		t.Fatal(err)
	}

	member, _ := members[0], beaconChains[0]
	member2, beaconChain2 := members[1], beaconChains[1]

	message2, err := member2.SignDKGResult(
		dkgResult,
		beaconChain2,
	)

	state := &resultSigningState{
		member:            member,
		signatureMessages: make([]*DKGResultHashSignatureMessage, 0),
	}

	state.Receive(&mockSignatureMessage{
		message2,
		[]byte("operator uses another key"),
	})

	if len(state.signatureMessages) != 0 {
		t.Fatalf("Expected no signature hash message accepted")
	}
}

type mockSignatureMessage struct {
	payload         *DKGResultHashSignatureMessage
	senderPublicKey []byte
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
func (msm *mockSignatureMessage) SenderPublicKey() []byte {
	return msm.senderPublicKey
}
func (msm *mockSignatureMessage) Seqno() uint64 {
	panic("not implemented")
}
