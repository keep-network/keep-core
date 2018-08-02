package thresholdsignature

import (
	"os"
	"testing"

	"github.com/keep-network/keep-core/pkg/internal/blsutils"
	"github.com/keep-network/keep-core/pkg/internal/pbutils"
	"github.com/keep-network/keep-core/pkg/internal/testutils"

	"github.com/dfinity/go-dfinity-crypto/bls"
)

func TestMain(m *testing.M) {
	bls.Init(bls.CurveSNARK1)

	os.Exit(m.Run())
}

func TestSignatureShareMessageRoundTrip(t *testing.T) {
	id := blsutils.GenerateID()

	msg := &SignatureShareMessage{id, make([]byte, 0)}
	unmarshaled := &SignatureShareMessage{}
	err := pbutils.RoundTrip(msg, unmarshaled)
	if err != nil {
		t.Fatal(err)
	}

	blsutils.AssertIDsEqual(t, msg.ID, unmarshaled.ID)
	testutils.AssertBytesEqual(t, msg.ShareBytes, unmarshaled.ShareBytes)
}
