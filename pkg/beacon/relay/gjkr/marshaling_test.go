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

func TestOtherMemberSharesMessageRoundtrip(t *testing.T) {
	msg := &OtherMemberSharesMessage{
		senderID:        MemberID(997),
		receiverID:      MemberID(112),
		encryptedShareS: []byte{0x01, 0x02, 0x03, 0x04, 0x05},
		encryptedShareT: []byte{0x0F, 0x0E, 0x0D, 0x0C, 0x0B},
	}
	unmarshaled := &OtherMemberSharesMessage{}

	err := pbutils.RoundTrip(msg, unmarshaled)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(msg, unmarshaled) {
		t.Fatalf("unexpected content of unmarshaled message")
	}
}

func TestSecretSharesAccusationsMessageRoundtrip(t *testing.T) {
	msg := &SecretSharesAccusationsMessage{
		senderID: MemberID(12121),
		accusedIDs: []MemberID{
			MemberID(1283),
			MemberID(9712),
			MemberID(8141),
		},
	}
	unmarshaled := &SecretSharesAccusationsMessage{}

	err := pbutils.RoundTrip(msg, unmarshaled)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(msg, unmarshaled) {
		t.Fatalf("unexpected content of unmarshaled message")
	}
}

func TestMemberPublicKeySharePointsMessageRoundtrip(t *testing.T) {
	msg := &MemberPublicKeySharePointsMessage{
		senderID: MemberID(987112),
		publicKeySharePoints: []*big.Int{
			big.NewInt(18211),
			big.NewInt(12311),
			big.NewInt(18828),
			big.NewInt(88711),
		},
	}
	unmarshaled := &MemberPublicKeySharePointsMessage{}

	err := pbutils.RoundTrip(msg, unmarshaled)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(msg, unmarshaled) {
		t.Fatalf("unexpected content of unmarshaled message")
	}
}

func TestPointsAccusationsMessage(t *testing.T) {
	msg := &PointsAccusationsMessage{
		senderID: MemberID(129841),
		accusedIDs: []MemberID{
			MemberID(12818),
			MemberID(91819),
			MemberID(61616),
		},
	}
	unmarshaled := &PointsAccusationsMessage{}

	err := pbutils.RoundTrip(msg, unmarshaled)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(msg, unmarshaled) {
		t.Fatalf("unexpected content of unmarshaled message")
	}
}
