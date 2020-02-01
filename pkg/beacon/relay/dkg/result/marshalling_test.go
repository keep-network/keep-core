package result

import (
	"reflect"
	"testing"

	"github.com/keep-network/keep-core/pkg/internal/pbutils"
)

func TestDKGResultHashSignatureMessageRoundtrip(t *testing.T) {
	msg := &DKGResultHashSignatureMessage{
		senderIndex: 10,
		resultHash:  [32]byte{30},
		signature:   []byte("signature"),
		publicKey:   []byte("pubkey"),
	}

	unmarshaled := &DKGResultHashSignatureMessage{}

	err := pbutils.RoundTrip(msg, unmarshaled)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(msg, unmarshaled) {
		t.Fatalf("unexpected content of unmarshaled message")
	}
}
