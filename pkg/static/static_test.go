package static

import (
	crand "crypto/rand"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func TestSignAndVerifyRoundTrip(t *testing.T) {
	privateKey, publicKey, err := GenerateStaticKeyPair(crand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	hash := hexutil.MustDecode("0xab530a13e45914982b79f9b7e3fba994cfd1f3fb22f71cea1afbf02b460c6d1d")

	signature, err := privateKey.Sign(hash)
	if err != nil {
		t.Fatal(err)
	}

	if !VerifySignature(publicKey, hash, signature) {
		t.Error("invalid signature")
	}
}
