package result

import (
	"math/big"
	"reflect"
	"testing"

	relayChain "github.com/keep-network/keep-core/pkg/beacon/relay/chain"
	"github.com/keep-network/keep-core/pkg/chain/local"
	"github.com/keep-network/keep-core/pkg/net"
	"github.com/keep-network/keep-core/pkg/operator"
)

func TestAcceptValidSignatureHashMessage(t *testing.T) {
	groupSize := 2
	threshold := 2
	minimumStake := big.NewInt(200)

	dkgResult := &relayChain.DKGResult{
		GroupPublicKey: []byte("He’s the hero Gotham deserves."),
	}

	chainHandle := local.Connect(groupSize, threshold, minimumStake)

	members, err := initializeSigningMembers(groupSize)
	if err != nil {
		t.Fatal(err)
	}

	member := members[0]
	member2 := members[1]

	message2, err := member2.SignDKGResult(dkgResult, chainHandle.ThresholdRelay())

	state := &resultSigningState{
		member:            member,
		signatureMessages: make([]*DKGResultHashSignatureMessage, 0),
	}

	state.Receive(&mockSignatureMessage{
		message2,
		operator.Marshal(&member2.privateKey.PublicKey),
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
	threshold := 2
	minimumStake := big.NewInt(200)

	dkgResult := &relayChain.DKGResult{
		GroupPublicKey: []byte("But not the one it needs right now."),
	}

	chainHandle := local.Connect(groupSize, threshold, minimumStake)

	members, err := initializeSigningMembers(groupSize)
	if err != nil {
		t.Fatal(err)
	}

	member := members[0]
	member2 := members[1]

	state := &resultSigningState{
		member:            member,
		signatureMessages: make([]*DKGResultHashSignatureMessage, 0),
	}

	message2, err := member2.SignDKGResult(dkgResult, chainHandle.ThresholdRelay())

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
