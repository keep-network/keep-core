package dkg

import (
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/dfinity/go-dfinity-crypto/bls"
	"github.com/gogo/protobuf/proto"
)

func TestMain(m *testing.M) {
	bls.Init(bls.CurveSNARK1)

	os.Exit(m.Run())
}

func TestJoinMessageRoundTrip(t *testing.T) {
	id := generateBlsID()

	msg := &JoinMessage{id}
	unmarshaled := &JoinMessage{}
	err := roundTrip(msg, unmarshaled)
	if err != nil {
		t.Fatal(err)
	}

	assertIDRoundTrip(t, msg.id, unmarshaled.id)
}

func TestMemberCommitmentsMessageRoundTrip(t *testing.T) {
	id := generateBlsID()
	commitments := make([]bls.PublicKey, 0)
	for i := 0; i < 10; i++ {
		pk := generateBlsPublicKey()
		commitments = append(commitments, *pk)
	}

	msg := &MemberCommitmentsMessage{id, commitments}
	unmarshaled := &MemberCommitmentsMessage{}
	err := roundTrip(msg, unmarshaled)
	if err != nil {
		t.Fatal(err)
	}

	assertIDRoundTrip(t, msg.id, unmarshaled.id)
	assertEqual(
		t,
		len(msg.Commitments),
		len(unmarshaled.Commitments),
		"Expected commitment length to be equal pre- and post-round-trip")

	for i, commitment := range msg.Commitments {
		assertPublicKeyRoundTrip(t, &commitment, &unmarshaled.Commitments[i])
	}
}

func TestMemberShareMessageRoundTrip(t *testing.T) {
	id := generateBlsID()
	receiverID := generateBlsID()
	share := generateBlsSecretKey()

	msg := &MemberShareMessage{id, receiverID, share}
	unmarshaled := &MemberShareMessage{}
	err := roundTrip(msg, unmarshaled)
	if err != nil {
		t.Fatal(err)
	}

	assertIDRoundTrip(t, msg.id, unmarshaled.id)
	assertIDRoundTrip(t, msg.receiverID, unmarshaled.receiverID)
	assertSecretKeyRoundTrip(t, msg.Share, unmarshaled.Share)
}

func TestAccusationsMessageRoundTrip(t *testing.T) {
	id := generateBlsID()
	accusedIDs := make([]*bls.ID, 0)
	for i := 0; i < 10; i++ {
		accusedIDs = append(accusedIDs, generateBlsID())
	}

	msg := &AccusationsMessage{id, accusedIDs}
	unmarshaled := &AccusationsMessage{}
	err := roundTrip(msg, unmarshaled)
	if err != nil {
		t.Fatal(err)
	}

	assertIDRoundTrip(t, msg.id, unmarshaled.id)
	assertEqual(
		t,
		len(msg.accusedIDs),
		len(unmarshaled.accusedIDs),
		"Expected accused IDs length to be equal pre- and post-round-trip")
	for i, id := range msg.accusedIDs {
		assertIDRoundTrip(t, id, unmarshaled.accusedIDs[i])
	}
}

func TestJustificationsMessageRoundTrip(t *testing.T) {
	justifications := make(map[bls.ID]*bls.SecretKey)
	for i := 0; i < 10; i++ {
		justificationID := generateBlsID()
		sk := generateBlsSecretKey()
		justifications[*justificationID] = sk
	}

	msg := &JustificationsMessage{generateBlsID(), justifications}
	unmarshaled := &JustificationsMessage{}
	err := roundTrip(msg, unmarshaled)
	if err != nil {
		t.Fatal(err)
	}

	assertIDRoundTrip(t, msg.id, unmarshaled.id)
	assertEqual(
		t,
		len(msg.justifications),
		len(unmarshaled.justifications),
		"Expected justifications length to be equal pre- and post-round-trip")
	for id, sk := range msg.justifications {
		if unmarshaledSk, ok := unmarshaled.justifications[id]; !ok {
			t.Errorf(
				"Expected starting id [%v] to exist in round-trip justifications",
				id.GetHexString())
		} else {
			assertSecretKeyRoundTrip(t, sk, unmarshaledSk)
		}
	}
}

func assertIDRoundTrip(t *testing.T, id *bls.ID, roundTripID *bls.ID) {
	if !id.IsEqual(roundTripID) {
		t.Errorf(
			"ID failed to round-trip: [%s] != [%s]",
			id.GetHexString(),
			roundTripID.GetHexString())
	}
}

func assertEqual(t *testing.T, n int, n2 int, msg string) {
	if n != n2 {
		t.Errorf("%v: [%v] != [%v]", msg, n, n2)
	}
}

func assertPublicKeyRoundTrip(t *testing.T, pk1 *bls.PublicKey, pk2 *bls.PublicKey) {
	if !pk1.IsEqual(pk2) {
		t.Errorf(
			"bls.PublicKey failed to round-trip: [%s] != [%s]",
			pk1.GetHexString(),
			pk2.GetHexString())
	}
}

func assertSecretKeyRoundTrip(t *testing.T, sk1 *bls.SecretKey, sk2 *bls.SecretKey) {
	if !sk1.IsEqual(sk2) {
		t.Errorf(
			"bls.SecretKey failed to round-trip: [%s] != [%s]",
			sk1.GetHexString(),
			sk2.GetHexString())
	}
}

func roundTrip(
	marshaler proto.Marshaler,
	unmarshaler proto.Unmarshaler) error {
	bytes, err := marshaler.Marshal()
	if err != nil {
		return err
	}

	err = unmarshaler.Unmarshal(bytes)
	if err != nil {
		return err
	}

	return nil
}

func generateBlsID() *bls.ID {
	id := bls.ID{}
	idValue := fmt.Sprintf("%v", rand.Int31())
	err := id.SetDecString(idValue)
	if err != nil {
		panic(fmt.Sprintf(
			"Failed to generate id from random number %v: [%v]",
			idValue,
			err))
	}

	return &id
}

func generateBlsPublicKey() *bls.PublicKey {
	sk := bls.SecretKey{}
	sk.SetByCSPRNG()
	pk := sk.GetPublicKey()

	return pk
}

func generateBlsSecretKey() *bls.SecretKey {
	sk := bls.SecretKey{}
	sk.SetByCSPRNG()

	return &sk
}
