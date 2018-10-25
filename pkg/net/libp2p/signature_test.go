package libp2p

import (
	"strings"
	"testing"
)

func TestVerifyMessageSignature(t *testing.T) {
	identity, err := newTestIdentity()

	ch := &channel{
		clientIdentity: identity,
	}

	msg := []byte("It's not much of a tail, but I'm sort of attached to it.")

	signature, err := ch.sign(msg)
	if err != nil {
		t.Fatal(err)
	}

	err = ch.verify(identity.id, msg, signature)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDetectInvalidMessageSignature(t *testing.T) {
	identity, err := newTestIdentity()

	ch := &channel{
		clientIdentity: identity,
	}

	msg := []byte("It's not much of a tail, but I'm sort of attached to it.")

	signature, err := ch.sign(msg)
	if err != nil {
		t.Fatal(err)
	}

	anotherIdentity, err := newTestIdentity()
	if err != nil {
		t.Fatal(err)
	}

	err = ch.verify(anotherIdentity.id, msg, signature)
	if err == nil {
		t.Fatal("signature validation should fail")
	}

	if !strings.HasPrefix(err.Error(), "invalid signature") {
		t.Fatalf("error other than expected: %v", err)
	}
}
