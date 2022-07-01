package gjkr

import (
	"math/big"
	"reflect"
	"testing"

	fuzz "github.com/google/gofuzz"

	bn256 "github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"
	"github.com/keep-network/keep-core/pkg/beacon/relay/group"
	"github.com/keep-network/keep-core/pkg/internal/pbutils"
	"github.com/keep-network/keep-core/pkg/crypto/ephemeral"
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

func TestFuzzEphemeralPublicKeyMessageRoundtrip(t *testing.T) {
	for i := 0; i < 10; i++ {
		var (
			senderID            group.MemberIndex
			ephemeralPublicKeys map[group.MemberIndex]*ephemeral.PublicKey
		)

		f := fuzz.New().NilChance(0.1).
			NumElements(0, 512).
			Funcs(pbutils.FuzzFuncs()...)

		f.Fuzz(&senderID)
		f.Fuzz(&ephemeralPublicKeys)

		message := &EphemeralPublicKeyMessage{
			senderID:            senderID,
			ephemeralPublicKeys: ephemeralPublicKeys,
		}

		_ = pbutils.RoundTrip(message, &EphemeralPublicKeyMessage{})
	}
}

func TestFuzzEphemeralPublicKeyMessageUnmarshaler(t *testing.T) {
	pbutils.FuzzUnmarshaler(&EphemeralPublicKeyMessage{})
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

func TestFuzzMemberCommitmentsMessageRoundtrip(t *testing.T) {
	for i := 0; i < 10; i++ {
		var (
			senderID    group.MemberIndex
			commitments []*bn256.G1
		)

		f := fuzz.New().NilChance(0.1).
			NumElements(0, 512).
			Funcs(pbutils.FuzzFuncs()...)

		f.Fuzz(&senderID)
		f.Fuzz(&commitments)

		message := &MemberCommitmentsMessage{
			senderID:    senderID,
			commitments: commitments,
		}

		_ = pbutils.RoundTrip(message, &MemberCommitmentsMessage{})
	}
}

func TestFuzzMemberCommitmentsMessageUnmarshaler(t *testing.T) {
	pbutils.FuzzUnmarshaler(&MemberCommitmentsMessage{})
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

func TestFuzzPeerSharesMessageRoundtrip(t *testing.T) {
	for i := 0; i < 10; i++ {
		var (
			senderID group.MemberIndex
			shares   map[group.MemberIndex]*peerShares
		)

		fuzzPeerShares := func(shares *peerShares, c fuzz.Continue) {
			var encryptedShareS, encryptedShareT []byte

			c.Fuzz(&encryptedShareS)
			c.Fuzz(&encryptedShareT)

			shares.encryptedShareS = encryptedShareS
			shares.encryptedShareT = encryptedShareT
		}

		fuzzFuncs := []interface{}{fuzzPeerShares}
		fuzzFuncs = append(fuzzFuncs, pbutils.FuzzFuncs()...)

		f := fuzz.New().NilChance(0.1).
			NumElements(0, 512).
			Funcs(fuzzFuncs...)

		f.Fuzz(&senderID)
		f.Fuzz(&shares)

		message := &PeerSharesMessage{
			senderID: senderID,
			shares:   shares,
		}

		_ = pbutils.RoundTrip(message, &PeerSharesMessage{})
	}
}

func TestFuzzPeerSharesMessageUnmarshaler(t *testing.T) {
	pbutils.FuzzUnmarshaler(&PeerSharesMessage{})
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

func TestFuzzSecretSharesAccusationsMessageRoundtrip(t *testing.T) {
	for i := 0; i < 10; i++ {
		var (
			senderID           group.MemberIndex
			accusedMembersKeys map[group.MemberIndex]*ephemeral.PrivateKey
		)

		f := fuzz.New().NilChance(0.1).
			NumElements(0, 512).
			Funcs(pbutils.FuzzFuncs()...)

		f.Fuzz(&senderID)
		f.Fuzz(&accusedMembersKeys)

		message := &SecretSharesAccusationsMessage{
			senderID:           senderID,
			accusedMembersKeys: accusedMembersKeys,
		}

		_ = pbutils.RoundTrip(message, &SecretSharesAccusationsMessage{})
	}
}

func TestFuzzSecretSharesAccusationsMessageUnmarshaler(t *testing.T) {
	pbutils.FuzzUnmarshaler(&SecretSharesAccusationsMessage{})
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

func TestFuzzMemberPublicKeySharePointsMessageRoundtrip(t *testing.T) {
	for i := 0; i < 10; i++ {
		var (
			senderID             group.MemberIndex
			publicKeySharePoints []*bn256.G2
		)

		f := fuzz.New().NilChance(0.1).
			NumElements(0, 512).
			Funcs(pbutils.FuzzFuncs()...)

		f.Fuzz(&senderID)
		f.Fuzz(&publicKeySharePoints)

		message := &MemberPublicKeySharePointsMessage{
			senderID:             senderID,
			publicKeySharePoints: publicKeySharePoints,
		}

		_ = pbutils.RoundTrip(message, &MemberPublicKeySharePointsMessage{})
	}
}

func TestFuzzMemberPublicKeySharePointsMessageUnmarshaler(t *testing.T) {
	pbutils.FuzzUnmarshaler(&MemberPublicKeySharePointsMessage{})
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

func TestFuzzPointsAccusationsMessageRoundtrip(t *testing.T) {
	for i := 0; i < 10; i++ {
		var (
			senderID           group.MemberIndex
			accusedMembersKeys map[group.MemberIndex]*ephemeral.PrivateKey
		)

		f := fuzz.New().NilChance(0.1).
			NumElements(0, 512).
			Funcs(pbutils.FuzzFuncs()...)

		f.Fuzz(&senderID)
		f.Fuzz(&accusedMembersKeys)

		message := &PointsAccusationsMessage{
			senderID:           senderID,
			accusedMembersKeys: accusedMembersKeys,
		}

		_ = pbutils.RoundTrip(message, &PointsAccusationsMessage{})
	}
}

func TestFuzzPointsAccusationsMessageUnmarshaler(t *testing.T) {
	pbutils.FuzzUnmarshaler(&PointsAccusationsMessage{})
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

func TestFuzzMisbehavedEphemeralKeysMessageRoundtrip(t *testing.T) {
	for i := 0; i < 10; i++ {
		var (
			senderID    group.MemberIndex
			privateKeys map[group.MemberIndex]*ephemeral.PrivateKey
		)

		f := fuzz.New().NilChance(0.1).
			NumElements(0, 512).
			Funcs(pbutils.FuzzFuncs()...)

		f.Fuzz(&senderID)
		f.Fuzz(&privateKeys)

		message := &MisbehavedEphemeralKeysMessage{
			senderID:    senderID,
			privateKeys: privateKeys,
		}

		_ = pbutils.RoundTrip(message, &MisbehavedEphemeralKeysMessage{})
	}
}

func TestFuzzMisbehavedEphemeralKeysMessageUnmarshaler(t *testing.T) {
	pbutils.FuzzUnmarshaler(&MisbehavedEphemeralKeysMessage{})
}
