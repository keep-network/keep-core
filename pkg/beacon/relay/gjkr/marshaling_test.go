package gjkr

import (
	"math/big"
	"reflect"
	"testing"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/keep-network/keep-core/pkg/beacon/relay/member"
	"github.com/keep-network/keep-core/pkg/internal/pbutils"
	"github.com/keep-network/keep-core/pkg/net/ephemeral"
)

func TestJoinMessageRoundtrip(t *testing.T) {
	msg := &JoinMessage{member.MemberIndex(1337)}
	unmarshaled := &JoinMessage{}

	err := pbutils.RoundTrip(msg, unmarshaled)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(msg, unmarshaled) {
		t.Fatalf("unexpected content of unmarshaled message")
	}
}

func TestEphemeralPublicKeyMessageRoundtrip(t *testing.T) {
	keyPair1, err := ephemeral.GenerateKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	keyPair2, err := ephemeral.GenerateKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	publicKeys := make(map[member.MemberIndex]*ephemeral.PublicKey)
	publicKeys[member.MemberIndex(2181)] = keyPair1.PublicKey
	publicKeys[member.MemberIndex(9119)] = keyPair2.PublicKey

	msg := &EphemeralPublicKeyMessage{
		senderID:            member.MemberIndex(3548),
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
		senderID: member.MemberIndex(1410),
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
	shares := make(map[member.MemberIndex]*peerShares)
	shares[member.MemberIndex(112)] = &peerShares{
		encryptedShareS: []byte{0x01, 0x02, 0x03, 0x04, 0x05},
		encryptedShareT: []byte{0x0F, 0x0E, 0x0D, 0x0C, 0x0B},
	}
	shares[member.MemberIndex(223)] = &peerShares{
		encryptedShareS: []byte{0x0A, 0x0E, 0x0F, 0x0F, 0x0F},
		encryptedShareT: []byte{0x01, 0x0F, 0x0E, 0x0E, 0x0D},
	}

	msg := &PeerSharesMessage{
		senderID: member.MemberIndex(997),
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
		senderID: member.MemberIndex(12121),
		accusedMembersKeys: map[member.MemberIndex]*ephemeral.PrivateKey{
			member.MemberIndex(1283): keyPair1.PrivateKey,
			member.MemberIndex(9712): keyPair2.PrivateKey,
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
		senderID: member.MemberIndex(987112),
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
		senderID: member.MemberIndex(129841),
		accusedMembersKeys: map[member.MemberIndex]*ephemeral.PrivateKey{
			member.MemberIndex(12341): keyPair1.PrivateKey,
			member.MemberIndex(51111): keyPair2.PrivateKey,
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

func TestDisqualifiedEphemeralKeysMessageRoundtrip(t *testing.T) {
	keyPair1, err := ephemeral.GenerateKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	keyPair2, err := ephemeral.GenerateKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	msg := &DisqualifiedEphemeralKeysMessage{
		senderID: member.MemberIndex(181811),
		privateKeys: map[member.MemberIndex]*ephemeral.PrivateKey{
			member.MemberIndex(1821881): keyPair1.PrivateKey,
			member.MemberIndex(8181818): keyPair2.PrivateKey,
		},
	}
	unmarshaled := &DisqualifiedEphemeralKeysMessage{}

	err = pbutils.RoundTrip(msg, unmarshaled)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(msg, unmarshaled) {
		t.Fatalf("unexpected content of unmarshaled message")
	}
}
