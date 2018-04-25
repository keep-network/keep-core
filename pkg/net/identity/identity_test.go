package identity

import (
	"testing"

	"github.com/keep-network/keep-core/pkg/net"
	peer "github.com/libp2p/go-libp2p-peer"
)

func TestAddIdentityToStore(t *testing.T) {
	identity, err := LoadOrGenerateIdentity(0, "")
	if err != nil {
		t.Fatalf("Failed to generate valid PeerIdentity with err: %s", err)
	}
	pi := identity.(*peerIdentity)
	ps, err := AddIdentityToStore(pi)
	if err != nil {
		t.Fatalf("Failed to add Identity to store with err: %s", err)
	}

	var match bool
	for _, p := range ps.Peers() {
		if p == toPeerID(pi.ID()) {
			match = true
		}
	}
	if !match {
		t.Fatalf("Failed to add Identity with ID %+v to the PeerStore", toPeerID(pi.ID()))
	}
}

func TestPublicKeyFunctions(t *testing.T) {
	identity, err := LoadOrGenerateIdentity(0, "")
	if err != nil {
		t.Fatalf("Failed to generate valid PeerIdentity with err: %s", err)
	}
	pi := identity.(*peerIdentity)
	ps, err := AddIdentityToStore(pi)
	if err != nil {
		t.Fatalf("Failed to add Identity to store with err: %s", err)
	}

	peerID := toPeerID(pi.ID())
	msg := []byte("so random you can't fake it.")

	privKey := ps.PrivKey(peerID)
	sig, err := privKey.Sign(msg)
	if err != nil {
		t.Fatalf("Failed to sign msg with err: %s", err)
	}

	pubKey, err := pi.PubKeyFromID(pi.ID())
	if err != nil {
		t.Fatalf("Failed to generate public key from id with err %s", err)
	}

	ok, err := pubKey.Verify(msg, sig)
	if err != nil {
		t.Fatalf("Failed to verify msg with err: %s", err)
	}
	if !ok {
		t.Fatal("Failed to verify signature")
	}

	msg[0] = ^msg[0]
	ok, err = pubKey.Verify(msg, sig)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("Signature should not have matched with mutated data")
	}
}

func toPeerID(pi net.ID) peer.ID {
	return peer.ID(pi)
}
