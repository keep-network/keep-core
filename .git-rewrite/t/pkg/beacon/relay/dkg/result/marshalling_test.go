package result

import (
	"reflect"
	"testing"

	"github.com/keep-network/keep-core/pkg/internal/pbutils"
	"github.com/keep-network/keep-core/pkg/operator"
)

func TestDKGResultHashSignatureMessageRoundtrip(t *testing.T) {
	_, pubKey, err := operator.GenerateKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	msg := &DKGResultHashSignatureMessage{
		senderIndex: 1410,
		resultHash:  [32]byte{30},
		signature:   []byte("signature"),
		publicKey:   pubKey,
	}

	unmarshaled := &DKGResultHashSignatureMessage{}

	err = pbutils.RoundTrip(msg, unmarshaled)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(msg, unmarshaled) {
		t.Fatalf("unexpected content of unmarshaled message")
	}
}
