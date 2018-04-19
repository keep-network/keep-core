package identity

import "testing"

func TestSignAndVerifytMessage(t *testing.T) {
	pi, err := LoadOrGenerateIdentity(0, "")
	if err != nil {
		t.Fatalf("Failed to generate valid PeerIdentity with err: %s", err)
	}
	ps, err := pi.AddIdentityToStore()
	if err != nil {
		t.Fatalf("Failed to add Identity to store with err: %s", err)
	}

	msg := []byte("so random you can't fake it.")
	privKey := ps.PrivKey(pi.ID())

	sig, err := privKey.Sign(msg)
	if err != nil {
		t.Fatalf("Failed to sign msg with err: %s", err)
	}

	ok, err := pi.PubKey().Verify(msg, sig)
	if err != nil {
		t.Fatalf("Failed to verify msg with err: %s", err)
	}
	if !ok {
		t.Fatal("Failed to verify signature")
	}

	msg[0] = ^msg[0]
	ok, err = pi.PubKey().Verify(msg, sig)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("Signature should not have matched with mutated data")
	}

}
