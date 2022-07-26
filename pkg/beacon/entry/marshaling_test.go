package entry

import (
	"testing"

	fuzz "github.com/google/gofuzz"
	"github.com/keep-network/keep-core/pkg/protocol/group"

	"github.com/keep-network/keep-core/pkg/internal/pbutils"
	"github.com/keep-network/keep-core/pkg/internal/testutils"
)

func TestSignatureShareMessageRoundTrip(t *testing.T) {
	msg := &SignatureShareMessage{123, make([]byte, 0)}
	unmarshaled := &SignatureShareMessage{}

	err := pbutils.RoundTrip(msg, unmarshaled)
	if err != nil {
		t.Fatal(err)
	}

	if msg.senderID != unmarshaled.senderID {
		t.Errorf(
			"unexpected sender ID\nexpected: [%v]\nactual:   [%v]",
			msg.senderID,
			unmarshaled.senderID,
		)
	}

	testutils.AssertBytesEqual(t, msg.shareBytes, unmarshaled.shareBytes)
}

func TestFuzzSignatureShareMessageRoundtrip(t *testing.T) {
	for i := 0; i < 10; i++ {
		var (
			senderID   group.MemberIndex
			shareBytes []byte
		)

		f := fuzz.New().NilChance(0.1).NumElements(0, 512)

		f.Fuzz(&senderID)
		f.Fuzz(&shareBytes)

		message := &SignatureShareMessage{
			senderID:   senderID,
			shareBytes: shareBytes,
		}

		_ = pbutils.RoundTrip(message, &SignatureShareMessage{})
	}
}

func TestFuzzSignatureShareMessageUnmarshaler(t *testing.T) {
	pbutils.FuzzUnmarshaler(&SignatureShareMessage{})
}
