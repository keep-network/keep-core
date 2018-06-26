package dkg

import (
	"os"
	"testing"

	"github.com/dfinity/go-dfinity-crypto/bls"
	"github.com/keep-network/keep-core/pkg/internal/blsutils"
	"github.com/keep-network/keep-core/pkg/internal/pbutils"
)

func TestMain(m *testing.M) {
	bls.Init(bls.CurveSNARK1)

	os.Exit(m.Run())
}

func TestJoinMessageRoundTrip(t *testing.T) {
	id := blsutils.GenerateID()

	msg := &JoinMessage{id}
	unmarshaled := &JoinMessage{}
	err := pbutils.RoundTrip(msg, unmarshaled)
	if err != nil {
		t.Fatal(err)
	}

	blsutils.AssertIDsEqual(t, msg.id, unmarshaled.id)
}

func TestMemberCommitmentsMessageRoundTrip(t *testing.T) {
	id := blsutils.GenerateID()
	commitments := make([]bls.PublicKey, 0)
	for i := 0; i < 10; i++ {
		pk := blsutils.GeneratePublicKey()
		commitments = append(commitments, *pk)
	}

	msg := &MemberCommitmentsMessage{id, commitments}
	unmarshaled := &MemberCommitmentsMessage{}
	err := pbutils.RoundTrip(msg, unmarshaled)
	if err != nil {
		t.Fatal(err)
	}

	blsutils.AssertIDsEqual(t, msg.id, unmarshaled.id)
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
	id := blsutils.GenerateID()
	receiverID := blsutils.GenerateID()
	share := blsutils.GenerateSecretKey()

	msg := &MemberShareMessage{id, receiverID, share}
	unmarshaled := &MemberShareMessage{}
	err := pbutils.RoundTrip(msg, unmarshaled)
	if err != nil {
		t.Fatal(err)
	}

	blsutils.AssertIDsEqual(t, msg.id, unmarshaled.id)
	blsutils.AssertIDsEqual(t, msg.receiverID, unmarshaled.receiverID)
	assertSecretKeyRoundTrip(t, msg.Share, unmarshaled.Share)
}

func TestAccusationsMessageRoundTrip(t *testing.T) {
	id := blsutils.GenerateID()
	accusedIDs := make([]*bls.ID, 0)
	for i := 0; i < 10; i++ {
		accusedIDs = append(accusedIDs, blsutils.GenerateID())
	}

	msg := &AccusationsMessage{id, accusedIDs}
	unmarshaled := &AccusationsMessage{}
	err := pbutils.RoundTrip(msg, unmarshaled)
	if err != nil {
		t.Fatal(err)
	}

	blsutils.AssertIDsEqual(t, msg.id, unmarshaled.id)
	assertEqual(
		t,
		len(msg.accusedIDs),
		len(unmarshaled.accusedIDs),
		"Expected accused IDs length to be equal pre- and post-round-trip")
	for i, id := range msg.accusedIDs {
		blsutils.AssertIDsEqual(t, id, unmarshaled.accusedIDs[i])
	}
}

func TestJustificationsMessageRoundTrip(t *testing.T) {
	justifications := make(map[bls.ID]*bls.SecretKey)
	for i := 0; i < 10; i++ {
		justificationID := blsutils.GenerateID()
		sk := blsutils.GenerateSecretKey()
		justifications[*justificationID] = sk
	}

	msg := &JustificationsMessage{blsutils.GenerateID(), justifications}
	unmarshaled := &JustificationsMessage{}
	err := pbutils.RoundTrip(msg, unmarshaled)
	if err != nil {
		t.Fatal(err)
	}

	blsutils.AssertIDsEqual(t, msg.id, unmarshaled.id)
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
