package gjkr

import (
	"math/big"
	"reflect"
	"testing"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
	"github.com/keep-network/keep-core/pkg/internal/pbutils"
	"github.com/keep-network/keep-core/pkg/net/ephemeral"
)

func TestEphemeralPublicKeyMessageRoundtrip(t *testing.T) {
	keyPair1, err := ephemeral.GenerateKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	keyPair2, err := ephemeral.GenerateKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	publicKeys := make(map[group.MemberIndex]*ephemeral.PublicKey)
	publicKeys[group.MemberIndex(211)] = keyPair1.PublicKey
	publicKeys[group.MemberIndex(19)] = keyPair2.PublicKey

	msg := &EphemeralPublicKeyMessage{
		senderID:            group.MemberIndex(38),
		ephemeralPublicKeys: publicKeys,
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
		senderID: group.MemberIndex(141),
		commitments: []*bn256.G1{
			new(bn256.G1).ScalarBaseMult(big.NewInt(966)),
			new(bn256.G1).ScalarBaseMult(big.NewInt(1385)),
			new(bn256.G1).ScalarBaseMult(big.NewInt(1569)),
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

func TestPeerSharesMessageRoundtrip(t *testing.T) {
	shares := make(map[group.MemberIndex]*peerShares)
	shares[group.MemberIndex(112)] = &peerShares{
		encryptedShareS: []byte{0x01, 0x02, 0x03, 0x04, 0x05},
		encryptedShareT: []byte{0x0F, 0x0E, 0x0D, 0x0C, 0x0B},
	}
	shares[group.MemberIndex(223)] = &peerShares{
		encryptedShareS: []byte{0x0A, 0x0E, 0x0F, 0x0F, 0x0F},
		encryptedShareT: []byte{0x01, 0x0F, 0x0E, 0x0E, 0x0D},
	}

	msg := &PeerSharesMessage{
		senderID: group.MemberIndex(97),
		shares:   shares,
	}

	unmarshaled := &PeerSharesMessage{}

	err := pbutils.RoundTrip(msg, unmarshaled)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(msg, unmarshaled) {
		t.Fatalf("unexpected content of unmarshaled message")
	}
}

func TestSecretSharesAccusationsMessageRoundtrip(t *testing.T) {
	keyPair1, err := ephemeral.GenerateKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	keyPair2, err := ephemeral.GenerateKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	msg := &SecretSharesAccusationsMessage{
		senderID: group.MemberIndex(121),
		accusedMembersKeys: map[group.MemberIndex]*ephemeral.PrivateKey{
			group.MemberIndex(12): keyPair1.PrivateKey,
			group.MemberIndex(92): keyPair2.PrivateKey,
		},
	}
	unmarshaled := &SecretSharesAccusationsMessage{}

	err = pbutils.RoundTrip(msg, unmarshaled)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(msg, unmarshaled) {
		t.Fatalf("unexpected content of unmarshaled message")
	}
}

func TestMemberPublicKeySharePointsMessageRoundtrip(t *testing.T) {
	msg := &MemberPublicKeySharePointsMessage{
		senderID: group.MemberIndex(98),
		publicKeySharePoints: []*bn256.G2{
			new(bn256.G2).ScalarBaseMult(big.NewInt(18211)),
			new(bn256.G2).ScalarBaseMult(big.NewInt(12311)),
			new(bn256.G2).ScalarBaseMult(big.NewInt(18828)),
			new(bn256.G2).ScalarBaseMult(big.NewInt(88711)),
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

func TestPointsAccusationsMessageRoundtrip(t *testing.T) {
	keyPair1, err := ephemeral.GenerateKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	keyPair2, err := ephemeral.GenerateKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	msg := &PointsAccusationsMessage{
		senderID: group.MemberIndex(141),
		accusedMembersKeys: map[group.MemberIndex]*ephemeral.PrivateKey{
			group.MemberIndex(41): keyPair1.PrivateKey,
			group.MemberIndex(11): keyPair2.PrivateKey,
		},
	}
	unmarshaled := &PointsAccusationsMessage{}

	err = pbutils.RoundTrip(msg, unmarshaled)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(msg, unmarshaled) {
		t.Fatalf("unexpected content of unmarshaled message")
	}
}

func TestMisbehavedEphemeralKeysMessageRoundtrip(t *testing.T) {
	keyPair1, err := ephemeral.GenerateKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	keyPair2, err := ephemeral.GenerateKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	msg := &MisbehavedEphemeralKeysMessage{
		senderID: group.MemberIndex(18),
		privateKeys: map[group.MemberIndex]*ephemeral.PrivateKey{
			group.MemberIndex(181): keyPair1.PrivateKey,
			group.MemberIndex(88):  keyPair2.PrivateKey,
		},
	}
	unmarshaled := &MisbehavedEphemeralKeysMessage{}

	err = pbutils.RoundTrip(msg, unmarshaled)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(msg, unmarshaled) {
		t.Fatalf("unexpected content of unmarshaled message")
	}
}
