package gjkr

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/keep-network/keep-core/pkg/internal/pbutils"
	"github.com/keep-network/keep-core/pkg/net/ephemeral"
)

func TestEphemeralPublicKeyMessageRoundtrip(t *testing.T) {
	keyPair, err := ephemeral.GenerateKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	msg := &EphemeralPublicKeyMessage{
		senderID:           MemberID(123456789),
		receiverID:         MemberID(987654321),
		ephemeralPublicKey: keyPair.PublicKey,
	}
	unmarshaled := &EphemeralPublicKeyMessage{}

	err = pbutils.RoundTrip(msg, unmarshaled)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(msg, unmarshaled) {
		t.Fatalf("unexpected content of unmarshaled message")
	}
}

func TestMemberCommitmentsMessageRoundtrip(t *testing.T) {
	msg := &MemberCommitmentsMessage{
		senderID: MemberID(1410),
		commitments: []*big.Int{
			big.NewInt(966),
			big.NewInt(1385),
			big.NewInt(1569),
		},
	}
	unmarshaled := &MemberCommitmentsMessage{}

	err := pbutils.RoundTrip(msg, unmarshaled)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(msg, unmarshaled) {
		t.Fatalf("unexpected content of unmarshaled message")
	}
}
