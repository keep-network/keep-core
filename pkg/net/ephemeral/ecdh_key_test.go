package ephemeral

import (
	"crypto/rand"
	"reflect"
	"testing"
)

func TestMarshalUnmarshalPublicKey(t *testing.T) {
	_, pubKey, err := GenerateEphemeralKeypair(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	marshalled, err := pubKey.Marshal()
	if err != nil {
		t.Fatal(err)
	}

	unmarshalled := &PublicKey{}
	err = unmarshalled.Unmarshal(marshalled)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(unmarshalled, pubKey) {
		t.Fatal("unmarshalled public key does not match the original one")
	}
}

func TestMarshalUnmarshalPrivateKey(t *testing.T) {
	privKey, _, err := GenerateEphemeralKeypair(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}

	marshalled, err := privKey.Marshal()
	if err != nil {
		t.Fatal(err)
	}

	unmarshalled := &PrivateKey{}
	err = unmarshalled.Unmarshal(marshalled)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(unmarshalled, privKey) {
		t.Fatal("unmarshalled private key does not match the original one")
	}
}
