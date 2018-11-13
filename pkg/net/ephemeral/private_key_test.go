package ephemeral

import (
	"reflect"
	"testing"
)

func TestMarshalUnmarshalPublicKey(t *testing.T) {
	_, pubKey, err := GenerateEphemeralKeypair()
	if err != nil {
		t.Fatal(err)
	}

	marshalled := pubKey.Marshal()

	unmarshalled, err := UnmarshalPublicKey(marshalled)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(unmarshalled, pubKey) {
		t.Fatal("unmarshalled public key does not match the original one")
	}
}

func TestMarshalUnmarshalPrivateKey(t *testing.T) {
	privKey, _, err := GenerateEphemeralKeypair()
	if err != nil {
		t.Fatal(err)
	}

	marshalled := privKey.Marshal()

	unmarshalled := UnmarshalPrivateKey(marshalled)

	if !reflect.DeepEqual(unmarshalled, privKey) {
		t.Fatal("unmarshalled private key does not match the original one")
	}
}
