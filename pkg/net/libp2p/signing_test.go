package libp2p

import (
	"testing"

	"github.com/keep-network/keep-core/pkg/net/gen/pb"
)

func TestSignAndVerify(t *testing.T) {
	identity, err := generateIdentity()
	if err != nil {
		t.Fatal(err)
	}

	identityBytes, err := identity.Marshal()
	if err != nil {
		t.Fatal(err)
	}

	message := &pb.UnicastNetworkMessage{
		Sender:  identityBytes,
		Payload: []byte{5, 15, 25, 30, 35},
		Type:    []byte{1},
	}

	err = signMessage(message, identity.privKey)
	if err != nil {
		t.Fatal(err)
	}

	err = verifyMessageSignature(message, identity.pubKey)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAllFieldsButSignatureAreTaken(t *testing.T) {
	identity, err := generateIdentity()
	if err != nil {
		t.Fatal(err)
	}

	identityBytes, err := identity.Marshal()
	if err != nil {
		t.Fatal(err)
	}

	message := &pb.UnicastNetworkMessage{
		Sender:    identityBytes,
		Payload:   []byte{5, 15, 25, 30, 35},
		Type:      []byte{1},
		Signature: []byte{100, 200}, // should be ignored while signing
	}

	err = signMessage(message, identity.privKey)
	if err != nil {
		t.Fatal(err)
	}

	err = verifyMessageSignature(message, identity.pubKey)
	if err != nil {
		t.Fatal(err)
	}

	messageCopy := *message
	messageCopy.Sender = []byte{0}
	err = verifyMessageSignature(&messageCopy, identity.pubKey)
	if err == nil {
		t.Fatal("signature should fail because `Sender` field changed")
	}

	messageCopy = *message
	messageCopy.Payload = []byte{0}
	err = verifyMessageSignature(&messageCopy, identity.pubKey)
	if err == nil {
		t.Fatal("signature should fail because `Payload` field changed")
	}

	messageCopy = *message
	messageCopy.Type = []byte{0}
	err = verifyMessageSignature(&messageCopy, identity.pubKey)
	if err == nil {
		t.Fatal("signature should fail because `Type` field changed")
	}
}
