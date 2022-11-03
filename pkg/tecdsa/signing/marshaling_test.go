package signing

import (
	"reflect"
	"testing"

	fuzz "github.com/google/gofuzz"
	"github.com/keep-network/keep-core/pkg/crypto/ephemeral"
	"github.com/keep-network/keep-core/pkg/internal/pbutils"
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

func TestTssRoundOneCompositeMessage_MarshalingRoundtrip(t *testing.T) {
	msg := &tssRoundOneCompositeMessage{
		senderID:  group.MemberIndex(50),
		sessionID: "session-1",
		tssRoundOneMessages: map[string]*tssRoundOneMessage{
			"one": &tssRoundOneMessage{
				broadcastPayload: []byte{1, 2, 3, 4, 5},
				peersPayload: map[group.MemberIndex][]byte{
					1: {6, 7, 8, 9, 10},
					2: {11, 12, 13, 14, 15},
				},
			},
			"two": &tssRoundOneMessage{
				broadcastPayload: []byte{6, 7, 8, 9},
				peersPayload: map[group.MemberIndex][]byte{
					1: {16, 17, 18, 19, 20},
					2: {21, 22, 23, 24, 25},
				},
			},
		},
	}
	unmarshaled := &tssRoundOneCompositeMessage{}

	err := pbutils.RoundTrip(msg, unmarshaled)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(msg, unmarshaled) {
		t.Fatalf("unexpected content of unmarshaled message")
	}
}

func TestFuzzTssRoundOneCompositeMessage_MarshalingRoundtrip(t *testing.T) {
	for i := 0; i < 10; i++ {
		var (
			senderID            group.MemberIndex
			sessionID           string
			tssRoundOneMessages map[string]*tssRoundOneMessage
		)

		f := fuzz.New().NilChance(0.1).
			NumElements(0, 512).
			Funcs(pbutils.FuzzFuncs()...)

		f.Fuzz(&senderID)
		f.Fuzz(&sessionID)
		f.Fuzz(&tssRoundOneMessages)

		message := &tssRoundOneCompositeMessage{
			senderID:            senderID,
			sessionID:           sessionID,
			tssRoundOneMessages: tssRoundOneMessages,
		}

		_ = pbutils.RoundTrip(message, &tssRoundOneCompositeMessage{})
	}
}

func TestFuzzTssRoundOneCompositeMessage_Unmarshaler(t *testing.T) {
	pbutils.FuzzUnmarshaler(&tssRoundOneCompositeMessage{})
}

func TestTssRoundTwoCompositeMessage_MarshalingRoundtrip(t *testing.T) {
	msg := &tssRoundTwoCompositeMessage{
		senderID:  group.MemberIndex(50),
		sessionID: "session-1",
		tssRoundTwoMessages: map[string]*tssRoundTwoMessage{
			"one": &tssRoundTwoMessage{
				peersPayload: map[group.MemberIndex][]byte{
					1: {6, 7, 8, 9, 10},
					2: {11, 12, 13, 14, 15},
				},
			},
			"two": &tssRoundTwoMessage{
				peersPayload: map[group.MemberIndex][]byte{
					1: {16, 17, 18, 19, 20},
					2: {21, 22, 23, 24, 25},
				},
			},
		},
	}
	unmarshaled := &tssRoundTwoCompositeMessage{}

	err := pbutils.RoundTrip(msg, unmarshaled)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(msg, unmarshaled) {
		t.Fatalf("unexpected content of unmarshaled message")
	}
}

func TestFuzzTssRoundTwoCompositeMessage_MarshalingRoundtrip(t *testing.T) {
	for i := 0; i < 10; i++ {
		var (
			senderID            group.MemberIndex
			sessionID           string
			tssRoundTwoMessages map[string]*tssRoundTwoMessage
		)

		f := fuzz.New().NilChance(0.1).
			NumElements(0, 512).
			Funcs(pbutils.FuzzFuncs()...)

		f.Fuzz(&senderID)
		f.Fuzz(&sessionID)
		f.Fuzz(&tssRoundTwoMessages)

		message := &tssRoundTwoCompositeMessage{
			senderID:            senderID,
			sessionID:           sessionID,
			tssRoundTwoMessages: tssRoundTwoMessages,
		}

		_ = pbutils.RoundTrip(message, &tssRoundTwoCompositeMessage{})
	}
}

func TestFuzzTssRoundTwoCompositeMessage_Unmarshaler(t *testing.T) {
	pbutils.FuzzUnmarshaler(&tssRoundTwoCompositeMessage{})
}

func TestTssRoundThreeCompositeMessage_MarshalingRoundtrip(t *testing.T) {
	msg := &tssRoundThreeCompositeMessage{
		senderID:  group.MemberIndex(50),
		sessionID: "session-1",
		tssRoundThreeMessages: map[string]*tssRoundThreeMessage{
			"one": &tssRoundThreeMessage{
				broadcastPayload: []byte{1, 2, 3, 4, 5},
			},
			"two": &tssRoundThreeMessage{
				broadcastPayload: []byte{6, 7, 8, 9, 10},
			},
		},
	}
	unmarshaled := &tssRoundThreeCompositeMessage{}

	err := pbutils.RoundTrip(msg, unmarshaled)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(msg, unmarshaled) {
		t.Fatalf("unexpected content of unmarshaled message")
	}
}

func TestFuzzTssRoundThreeCompositeMessage_MarshalingRoundtrip(t *testing.T) {
	for i := 0; i < 10; i++ {
		var (
			senderID              group.MemberIndex
			sessionID             string
			tssRoundThreeMessages map[string]*tssRoundThreeMessage
		)

		f := fuzz.New().NilChance(0.1).
			NumElements(0, 512).
			Funcs(pbutils.FuzzFuncs()...)

		f.Fuzz(&senderID)
		f.Fuzz(&sessionID)
		f.Fuzz(&tssRoundThreeMessages)

		message := &tssRoundThreeCompositeMessage{
			senderID:              senderID,
			sessionID:             sessionID,
			tssRoundThreeMessages: tssRoundThreeMessages,
		}

		_ = pbutils.RoundTrip(message, &tssRoundThreeCompositeMessage{})
	}
}

func TestFuzzTssRoundThreeCompositeMessage_Unmarshaler(t *testing.T) {
	pbutils.FuzzUnmarshaler(&tssRoundThreeCompositeMessage{})
}

func TestTssRoundFourCompositeMessage_MarshalingRoundtrip(t *testing.T) {
	msg := &tssRoundFourCompositeMessage{
		senderID:  group.MemberIndex(50),
		sessionID: "session-1",
		tssRoundFourMessages: map[string]*tssRoundFourMessage{
			"one": &tssRoundFourMessage{
				broadcastPayload: []byte{1, 2, 3, 4, 5},
			},
			"two": &tssRoundFourMessage{
				broadcastPayload: []byte{6, 7, 8, 9, 10},
			},
		},
	}
	unmarshaled := &tssRoundFourCompositeMessage{}

	err := pbutils.RoundTrip(msg, unmarshaled)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(msg, unmarshaled) {
		t.Fatalf("unexpected content of unmarshaled message")
	}
}

func TestFuzzTssRoundFourCompositeMessage_MarshalingRoundtrip(t *testing.T) {
	for i := 0; i < 10; i++ {
		var (
			senderID             group.MemberIndex
			sessionID            string
			tssRoundFourMessages map[string]*tssRoundFourMessage
		)

		f := fuzz.New().NilChance(0.1).
			NumElements(0, 512).
			Funcs(pbutils.FuzzFuncs()...)

		f.Fuzz(&senderID)
		f.Fuzz(&sessionID)
		f.Fuzz(&tssRoundFourMessages)

		message := &tssRoundFourCompositeMessage{
			senderID:             senderID,
			sessionID:            sessionID,
			tssRoundFourMessages: tssRoundFourMessages,
		}

		_ = pbutils.RoundTrip(message, &tssRoundFourCompositeMessage{})
	}
}

func TestFuzzTssRoundFourCompositeMessage_Unmarshaler(t *testing.T) {
	pbutils.FuzzUnmarshaler(&tssRoundFourCompositeMessage{})
}

func TestTssRoundFiveCompositeMessage_MarshalingRoundtrip(t *testing.T) {
	msg := &tssRoundFiveCompositeMessage{
		senderID:  group.MemberIndex(50),
		sessionID: "session-1",
		tssRoundFiveMessages: map[string]*tssRoundFiveMessage{
			"one": &tssRoundFiveMessage{
				broadcastPayload: []byte{1, 2, 3, 4, 5},
			},
			"two": &tssRoundFiveMessage{
				broadcastPayload: []byte{6, 7, 8, 9, 10},
			},
		},
	}
	unmarshaled := &tssRoundFiveCompositeMessage{}

	err := pbutils.RoundTrip(msg, unmarshaled)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(msg, unmarshaled) {
		t.Fatalf("unexpected content of unmarshaled message")
	}
}

func TestFuzzTssRoundFiveCompositeMessage_MarshalingRoundtrip(t *testing.T) {
	for i := 0; i < 10; i++ {
		var (
			senderID             group.MemberIndex
			sessionID            string
			tssRoundFiveMessages map[string]*tssRoundFiveMessage
		)

		f := fuzz.New().NilChance(0.1).
			NumElements(0, 512).
			Funcs(pbutils.FuzzFuncs()...)

		f.Fuzz(&senderID)
		f.Fuzz(&sessionID)
		f.Fuzz(&tssRoundFiveMessages)

		message := &tssRoundFiveCompositeMessage{
			senderID:             senderID,
			sessionID:            sessionID,
			tssRoundFiveMessages: tssRoundFiveMessages,
		}

		_ = pbutils.RoundTrip(message, &tssRoundFiveCompositeMessage{})
	}
}

func TestFuzzTssRoundFiveCompositeMessage_Unmarshaler(t *testing.T) {
	pbutils.FuzzUnmarshaler(&tssRoundFiveCompositeMessage{})
}

func TestTssRoundSixCompositeMessage_MarshalingRoundtrip(t *testing.T) {
	msg := &tssRoundSixCompositeMessage{
		senderID:  group.MemberIndex(50),
		sessionID: "session-1",
		tssRoundSixMessages: map[string]*tssRoundSixMessage{
			"one": &tssRoundSixMessage{
				broadcastPayload: []byte{1, 2, 3, 4, 5},
			},
			"two": &tssRoundSixMessage{
				broadcastPayload: []byte{6, 7, 8, 9, 10},
			},
		},
	}
	unmarshaled := &tssRoundSixCompositeMessage{}

	err := pbutils.RoundTrip(msg, unmarshaled)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(msg, unmarshaled) {
		t.Fatalf("unexpected content of unmarshaled message")
	}
}

func TestFuzzTssRoundSixCompositeMessage_MarshalingRoundtrip(t *testing.T) {
	for i := 0; i < 10; i++ {
		var (
			senderID            group.MemberIndex
			sessionID           string
			tssRoundSixMessages map[string]*tssRoundSixMessage
		)

		f := fuzz.New().NilChance(0.1).
			NumElements(0, 512).
			Funcs(pbutils.FuzzFuncs()...)

		f.Fuzz(&senderID)
		f.Fuzz(&sessionID)
		f.Fuzz(&tssRoundSixMessages)

		message := &tssRoundSixCompositeMessage{
			senderID:            senderID,
			sessionID:           sessionID,
			tssRoundSixMessages: tssRoundSixMessages,
		}

		_ = pbutils.RoundTrip(message, &tssRoundSixCompositeMessage{})
	}
}

func TestFuzzTssRoundSixCompositeMessage_Unmarshaler(t *testing.T) {
	pbutils.FuzzUnmarshaler(&tssRoundSixCompositeMessage{})
}

func TestTssRoundSevenCompositeMessage_MarshalingRoundtrip(t *testing.T) {
	msg := &tssRoundSevenCompositeMessage{
		senderID:  group.MemberIndex(50),
		sessionID: "session-1",
		tssRoundSevenMessages: map[string]*tssRoundSevenMessage{
			"one": &tssRoundSevenMessage{
				broadcastPayload: []byte{1, 2, 3, 4, 5},
			},
			"two": &tssRoundSevenMessage{
				broadcastPayload: []byte{6, 7, 8, 9, 10},
			},
		},
	}
	unmarshaled := &tssRoundSevenCompositeMessage{}

	err := pbutils.RoundTrip(msg, unmarshaled)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(msg, unmarshaled) {
		t.Fatalf("unexpected content of unmarshaled message")
	}
}

func TestFuzzTssRoundSevenCompositeMessage_MarshalingRoundtrip(t *testing.T) {
	for i := 0; i < 10; i++ {
		var (
			senderID              group.MemberIndex
			sessionID             string
			tssRoundSevenMessages map[string]*tssRoundSevenMessage
		)

		f := fuzz.New().NilChance(0.1).
			NumElements(0, 512).
			Funcs(pbutils.FuzzFuncs()...)

		f.Fuzz(&senderID)
		f.Fuzz(&sessionID)
		f.Fuzz(&tssRoundSevenMessages)

		message := &tssRoundSevenCompositeMessage{
			senderID:  senderID,
			sessionID: sessionID,
		}

		_ = pbutils.RoundTrip(message, &tssRoundSevenCompositeMessage{})
	}
}

func TestFuzzTssRoundSevenCompositeMessage_Unmarshaler(t *testing.T) {
	pbutils.FuzzUnmarshaler(&tssRoundSevenCompositeMessage{})
}

func TestTssRoundEightCompositeMessage_MarshalingRoundtrip(t *testing.T) {
	msg := &tssRoundEightCompositeMessage{
		senderID:  group.MemberIndex(50),
		sessionID: "session-1",
		tssRoundEightMessages: map[string]*tssRoundEightMessage{
			"one": &tssRoundEightMessage{
				broadcastPayload: []byte{1, 2, 3, 4, 5},
			},
			"two": &tssRoundEightMessage{
				broadcastPayload: []byte{6, 7, 8, 9, 10},
			},
		},
	}
	unmarshaled := &tssRoundEightCompositeMessage{}

	err := pbutils.RoundTrip(msg, unmarshaled)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(msg, unmarshaled) {
		t.Fatalf("unexpected content of unmarshaled message")
	}
}

func TestFuzzTssRoundEightCompositeMessage_MarshalingRoundtrip(t *testing.T) {
	for i := 0; i < 10; i++ {
		var (
			senderID              group.MemberIndex
			sessionID             string
			tssRoundEightMessages map[string]*tssRoundEightMessage
		)

		f := fuzz.New().NilChance(0.1).
			NumElements(0, 512).
			Funcs(pbutils.FuzzFuncs()...)

		f.Fuzz(&senderID)
		f.Fuzz(&sessionID)
		f.Fuzz(&tssRoundEightMessages)

		message := &tssRoundEightCompositeMessage{
			senderID:              senderID,
			sessionID:             sessionID,
			tssRoundEightMessages: tssRoundEightMessages,
		}

		_ = pbutils.RoundTrip(message, &tssRoundEightCompositeMessage{})
	}
}

func TestFuzzTssRoundEightCompositeMessage_Unmarshaler(t *testing.T) {
	pbutils.FuzzUnmarshaler(&tssRoundEightCompositeMessage{})
}

func TestTssRoundNineCompositeMessage_MarshalingRoundtrip(t *testing.T) {
	msg := &tssRoundNineCompositeMessage{
		senderID:  group.MemberIndex(50),
		sessionID: "session-1",
		tssRoundNineMessages: map[string]*tssRoundNineMessage{
			"one": &tssRoundNineMessage{
				broadcastPayload: []byte{1, 2, 3, 4, 5},
			},
			"two": &tssRoundNineMessage{
				broadcastPayload: []byte{6, 7, 8, 9, 10},
			},
		},
	}
	unmarshaled := &tssRoundNineCompositeMessage{}

	err := pbutils.RoundTrip(msg, unmarshaled)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(msg, unmarshaled) {
		t.Fatalf("unexpected content of unmarshaled message")
	}
}

func TestFuzzTssRoundNineCompositeMessage_MarshalingRoundtrip(t *testing.T) {
	for i := 0; i < 10; i++ {
		var (
			senderID             group.MemberIndex
			sessionID            string
			tssRoundNineMessages map[string]*tssRoundNineMessage
		)

		f := fuzz.New().NilChance(0.1).
			NumElements(0, 512).
			Funcs(pbutils.FuzzFuncs()...)

		f.Fuzz(&senderID)
		f.Fuzz(&sessionID)
		f.Fuzz(&tssRoundNineMessages)

		message := &tssRoundNineCompositeMessage{
			senderID:             senderID,
			sessionID:            sessionID,
			tssRoundNineMessages: tssRoundNineMessages,
		}

		_ = pbutils.RoundTrip(message, &tssRoundNineCompositeMessage{})
	}
}

func TestFuzzTssRoundNineCompositeMessage_Unmarshaler(t *testing.T) {
	pbutils.FuzzUnmarshaler(&tssRoundNineCompositeMessage{})
}
