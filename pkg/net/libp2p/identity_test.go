package libp2p

import (
	"testing"

	testutils "github.com/libp2p/go-testutil"
)

func TestAddIdentityToStore(t *testing.T) {
	pi := generateDeterministicIdentity(t)

	ps, err := addIdentityToStore(pi)
	if err != nil {
		t.Fatal(err)
	}

	var match bool
	for _, p := range ps.Peers() {
		if p == pi.id {
			match = true
		}
	}
	if !match {
		t.Fatalf("Failed to add Identity with ID %+v to the PeerStore", pi.id)
	}
}

func TestPublicKeyFunctions(t *testing.T) {
	pi := generateDeterministicIdentity(t)

	ps, err := addIdentityToStore(pi)
	if err != nil {
		t.Fatalf("Failed to add Identity to store with err: %s", err)
	}

	msg := []byte("so random you can't fake it.")

	privKey := ps.PrivKey(pi.id)
	sig, err := privKey.Sign(msg)
	if err != nil {
		t.Fatalf("Failed to sign msg with err: %s", err)
	}

	pubKey := pubKeyFromIdentifier(pi)

	ok, err := pubKey.Verify(msg, sig)
	if err != nil {
		t.Fatalf("Failed to verify msg with err: %s", err)
	}
	if !ok {
		t.Fatal("Failed to verify signature")
	}

	msg[0] = ^msg[0]
	ok, err = pubKey.Verify(msg, sig)
	if err == nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("Signature should not have matched with mutated data")
	}
}

func generateDeterministicIdentity(t *testing.T) *peerIdentifier {
	p := testutils.RandPeerNetParamsOrFatal(t)
	return &peerIdentifier{id: p.ID, sk: p.PrivKey}
}
