package dkg

import (
	"reflect"
	"testing"

	fuzz "github.com/google/gofuzz"

	"github.com/keep-network/keep-core/pkg/crypto/ephemeral"
	"github.com/keep-network/keep-core/pkg/internal/pbutils"
	"github.com/keep-network/keep-core/pkg/internal/tecdsatest"
	"github.com/keep-network/keep-core/pkg/protocol/group"
)

func TestEphemeralPublicKeyMessage_MarshalingRoundtrip(t *testing.T) {
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

	msg := &ephemeralPublicKeyMessage{
		senderID:            group.MemberIndex(38),
		ephemeralPublicKeys: publicKeys,
		sessionID:           "session-1",
	}
	unmarshaled := &ephemeralPublicKeyMessage{}

	err = pbutils.RoundTrip(msg, unmarshaled)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(msg, unmarshaled) {
		t.Fatalf("unexpected content of unmarshaled message")
	}
}

func TestFuzzEphemeralPublicKeyMessage_MarshalingRoundtrip(t *testing.T) {
	for i := 0; i < 10; i++ {
		var (
			senderID            group.MemberIndex
			ephemeralPublicKeys map[group.MemberIndex]*ephemeral.PublicKey
			sessionID           string
		)

		f := fuzz.New().NilChance(0.1).
			NumElements(0, 512).
			Funcs(pbutils.FuzzFuncs()...)

		f.Fuzz(&senderID)
		f.Fuzz(&ephemeralPublicKeys)
		f.Fuzz(&sessionID)

		message := &ephemeralPublicKeyMessage{
			senderID:            senderID,
			ephemeralPublicKeys: ephemeralPublicKeys,
			sessionID:           sessionID,
		}

		_ = pbutils.RoundTrip(message, &ephemeralPublicKeyMessage{})
	}
}

func TestFuzzEphemeralPublicKeyMessage_Unmarshaler(t *testing.T) {
	pbutils.FuzzUnmarshaler(&ephemeralPublicKeyMessage{})
}

func TestTssRoundOneMessage_MarshalingRoundtrip(t *testing.T) {
	msg := &tssRoundOneMessage{
		senderID:         group.MemberIndex(50),
		broadcastPayload: []byte{1, 2, 3, 4, 5},
		sessionID:        "session-1",
	}
	unmarshaled := &tssRoundOneMessage{}

	err := pbutils.RoundTrip(msg, unmarshaled)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(msg, unmarshaled) {
		t.Fatalf("unexpected content of unmarshaled message")
	}
}

func TestFuzzTssRoundOneMessage_MarshalingRoundtrip(t *testing.T) {
	for i := 0; i < 10; i++ {
		var (
			senderID  group.MemberIndex
			payload   []byte
			sessionID string
		)

		f := fuzz.New().NilChance(0.1).
			NumElements(0, 512).
			Funcs(pbutils.FuzzFuncs()...)

		f.Fuzz(&senderID)
		f.Fuzz(&payload)
		f.Fuzz(&sessionID)

		message := &tssRoundOneMessage{
			senderID:         senderID,
			broadcastPayload: payload,
			sessionID:        sessionID,
		}

		_ = pbutils.RoundTrip(message, &tssRoundOneMessage{})
	}
}

func TestFuzzTssRoundOneMessage_Unmarshaler(t *testing.T) {
	pbutils.FuzzUnmarshaler(&tssRoundOneMessage{})
}

func TestTssRoundTwoMessage_MarshalingRoundtrip(t *testing.T) {
	msg := &tssRoundTwoMessage{
		senderID:         group.MemberIndex(50),
		broadcastPayload: []byte{1, 2, 3, 4, 5},
		peersPayload: map[group.MemberIndex][]byte{
			1: {6, 7, 8, 9, 10},
			2: {11, 12, 13, 14, 15},
		},
		sessionID: "session-1",
	}
	unmarshaled := &tssRoundTwoMessage{}

	err := pbutils.RoundTrip(msg, unmarshaled)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(msg, unmarshaled) {
		t.Fatalf("unexpected content of unmarshaled message")
	}
}

func TestFuzzTssRoundTwoMessage_MarshalingRoundtrip(t *testing.T) {
	for i := 0; i < 10; i++ {
		var (
			senderID         group.MemberIndex
			broadcastPayload []byte
			peersPayload     map[group.MemberIndex][]byte
			sessionID        string
		)

		f := fuzz.New().NilChance(0.1).
			NumElements(0, 512).
			Funcs(pbutils.FuzzFuncs()...)

		f.Fuzz(&senderID)
		f.Fuzz(&broadcastPayload)
		f.Fuzz(&peersPayload)
		f.Fuzz(&sessionID)

		message := &tssRoundTwoMessage{
			senderID:         senderID,
			broadcastPayload: broadcastPayload,
			peersPayload:     peersPayload,
			sessionID:        sessionID,
		}

		_ = pbutils.RoundTrip(message, &tssRoundTwoMessage{})
	}
}

func TestFuzzTssRoundTwoMessage_Unmarshaler(t *testing.T) {
	pbutils.FuzzUnmarshaler(&tssRoundTwoMessage{})
}

func TestTssRoundThreeMessage_MarshalingRoundtrip(t *testing.T) {
	msg := &tssRoundThreeMessage{
		senderID:         group.MemberIndex(50),
		broadcastPayload: []byte{1, 2, 3, 4, 5},
		sessionID:        "session-1",
	}
	unmarshaled := &tssRoundThreeMessage{}

	err := pbutils.RoundTrip(msg, unmarshaled)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(msg, unmarshaled) {
		t.Fatalf("unexpected content of unmarshaled message")
	}
}

func TestFuzzTssRoundThreeMessage_MarshalingRoundtrip(t *testing.T) {
	for i := 0; i < 10; i++ {
		var (
			senderID  group.MemberIndex
			payload   []byte
			sessionID string
		)

		f := fuzz.New().NilChance(0.1).
			NumElements(0, 512).
			Funcs(pbutils.FuzzFuncs()...)

		f.Fuzz(&senderID)
		f.Fuzz(&payload)
		f.Fuzz(&sessionID)

		message := &tssRoundThreeMessage{
			senderID:         senderID,
			broadcastPayload: payload,
			sessionID:        sessionID,
		}

		_ = pbutils.RoundTrip(message, &tssRoundThreeMessage{})
	}
}

func TestFuzzTssRoundThreeMessage_Unmarshaler(t *testing.T) {
	pbutils.FuzzUnmarshaler(&tssRoundThreeMessage{})
}

func TestTssFinalizationMessage_MarshalingRoundtrip(t *testing.T) {
	msg := &tssFinalizationMessage{
		senderID:  group.MemberIndex(50),
		sessionID: "session-1",
	}
	unmarshaled := &tssFinalizationMessage{}

	err := pbutils.RoundTrip(msg, unmarshaled)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(msg, unmarshaled) {
		t.Fatalf("unexpected content of unmarshaled message")
	}
}

func TestFuzzTssFinalizationMessage_MarshalingRoundtrip(t *testing.T) {
	for i := 0; i < 10; i++ {
		var (
			senderID  group.MemberIndex
			sessionID string
		)

		f := fuzz.New().NilChance(0.1).
			NumElements(0, 512).
			Funcs(pbutils.FuzzFuncs()...)

		f.Fuzz(&senderID)
		f.Fuzz(&sessionID)

		message := &tssFinalizationMessage{
			senderID:  senderID,
			sessionID: sessionID,
		}

		_ = pbutils.RoundTrip(message, &tssFinalizationMessage{})
	}
}

func TestFuzzTssFinalizationMessage_Unmarshaler(t *testing.T) {
	pbutils.FuzzUnmarshaler(&tssFinalizationMessage{})
}

func TestResultSignatureMessage_MarshalingRoundtrip(t *testing.T) {
	msg := &resultSignatureMessage{
		senderID:   123,
		resultHash: [32]byte{0: 11, 10: 22, 31: 33},
		signature:  []byte("signature"),
		publicKey:  []byte("pubkey"),
		sessionID:  "session-1",
	}
	unmarshaled := &resultSignatureMessage{}

	err := pbutils.RoundTrip(msg, unmarshaled)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(msg, unmarshaled) {
		t.Fatalf("unexpected content of unmarshaled message")
	}
}

func TestFuzzResultSignatureMessage_MarshalingRoundtrip(t *testing.T) {
	for i := 0; i < 10; i++ {
		var (
			senderID   group.MemberIndex
			resultHash ResultSignatureHash
			signature  []byte
			publicKey  []byte
			sessionID  string
		)

		f := fuzz.New().NilChance(0.1).NumElements(0, 512)

		f.Fuzz(&senderID)
		f.Fuzz(&resultHash)
		f.Fuzz(&signature)
		f.Fuzz(&publicKey)
		f.Fuzz(&sessionID)

		message := &resultSignatureMessage{
			senderID:   senderID,
			resultHash: resultHash,
			signature:  signature,
			publicKey:  publicKey,
			sessionID:  sessionID,
		}

		_ = pbutils.RoundTrip(message, &resultSignatureMessage{})
	}
}

func TestFuzzResultSignatureMessage_Unmarshaler(t *testing.T) {
	pbutils.FuzzUnmarshaler(&resultSignatureMessage{})
}

func TestPreParamsMarshalling(t *testing.T) {
	testData, err := tecdsatest.LoadPrivateKeyShareTestFixtures(1)
	if err != nil {
		t.Fatalf("failed to load test data: [%v]", err)
	}

	localPreParams := testData[0].LocalPreParams

	preParams := newPreParams(&localPreParams)

	unmarshaled := &PreParams{}

	if err := pbutils.RoundTrip(preParams, unmarshaled); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(preParams, unmarshaled) {
		t.Errorf(
			"unexpected content of unmarshaled pre-params\nexpected: %+v\nactual:   %+v\n",
			preParams,
			unmarshaled,
		)
	}

	// Check if PreParams Data pass the tss-lib validation.
	if !unmarshaled.data.ValidateWithProof() {
		t.Errorf("unmarshaled pre params data are invalid")
	}
}
