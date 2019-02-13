package thresholdsignature

import (
	"testing"

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

	testutils.AssertBytesEqual(t, msg.ShareBytes, unmarshaled.ShareBytes)
}
